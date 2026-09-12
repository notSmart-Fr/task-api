package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"task-api/internal/features/scraper"
	"task-api/internal/platform/httpclient"
)

const (
	StartURL  = "https://books.toscrape.com/catalogue/page-1.html"
	UserAgent = "FlyRankInternship-A9/1.0 (+https://github.com/yourusername/task-api)"
	Timeout   = 5 * time.Second
	MinDelay  = 500 * time.Millisecond
	CacheDir  = "cache"
	OutputDir = "output"
	MaxPages  = 3
)

type DiscoveredItem struct {
	URL        string
	SourcePage string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	startTime := time.Now().UTC()
	client := httpclient.NewClient(UserAgent, Timeout, MinDelay, CacheDir)

	currentURL := StartURL
	pageCount := 0

	urlMap := make(map[string]bool)
	var discoveredItems []DiscoveredItem

	// -------------------------------------------------------------------------
	// STAGE 2: Discover catalogue pages
	// -------------------------------------------------------------------------
	for currentURL != "" && pageCount < MaxPages {
		pageCount++

		req, err := http.NewRequest(http.MethodGet, currentURL, nil)
		if err != nil {
			slog.Error("Failed to create request", "url", currentURL, "error", err)
			break
		}

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Request failed", "url", currentURL, "error", err)
			break
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			slog.Error("Unexpected status code", "url", currentURL, "status", resp.StatusCode)
			break
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			slog.Error("Failed to read body", "url", currentURL, "error", err)
			break
		}

		bookURLs, nextURL, err := scraper.ParseCataloguePage(body, currentURL)
		if err != nil {
			slog.Error("Failed to parse catalogue page", "url", currentURL, "error", err)
			break
		}

		for _, bookURL := range bookURLs {
			if !urlMap[bookURL] {
				urlMap[bookURL] = true
				discoveredItems = append(discoveredItems, DiscoveredItem{
					URL:        bookURL,
					SourcePage: currentURL,
				})
			}
		}

		currentURL = nextURL
	}

	// -------------------------------------------------------------------------
	// STAGE 5 RESILIENCE TEST: Inject 1 deliberately broken/fake URL
	// -------------------------------------------------------------------------
	discoveredItems = append(discoveredItems, DiscoveredItem{
		URL:        "https://books.toscrape.com/catalogue/deliberate-broken-fake-book_999999/index.html",
		SourcePage: StartURL,
	})

	// -------------------------------------------------------------------------
	// STAGE 3 & 5: Fetch detail pages with fault tolerance & single retry
	// -------------------------------------------------------------------------
	var rawRecords []*scraper.RawBookRecord
	failedPages := 0

	for _, item := range discoveredItems {
		body, statusCode, err := fetchWithRetry(client, item.URL)
		if err != nil || statusCode != http.StatusOK {
			slog.Warn("Skipping broken page", "url", item.URL, "status", statusCode, "error", err)
			failedPages++
			continue
		}

		record, err := scraper.ParseBookDetailPage(body, item.URL, item.SourcePage)
		if err != nil {
			slog.Warn("Failed to parse detail page", "url", item.URL, "error", err)
			failedPages++
			continue
		}

		rawRecords = append(rawRecords, record)
	}

	// -------------------------------------------------------------------------
	// STAGE 4: Clean, Normalize, Validate & Store
	// -------------------------------------------------------------------------
	seenCanonical := make(map[string]bool)
	var cleanRecords []*scraper.CleanBookRecord
	var validationErrors []*scraper.ValidationError

	for _, raw := range rawRecords {
		clean, valErr := scraper.NormalizeAndValidate(raw)
		if valErr != nil {
			validationErrors = append(validationErrors, valErr)
			continue
		}

		if !seenCanonical[clean.ProductURL] {
			seenCanonical[clean.ProductURL] = true
			cleanRecords = append(cleanRecords, clean)
		}
	}

	_ = os.MkdirAll(OutputDir, 0755)

	booksJSON, _ := json.MarshalIndent(cleanRecords, "", "  ")
	_ = os.WriteFile(filepath.Join(OutputDir, "books.json"), booksJSON, 0644)

	errorsJSON, _ := json.MarshalIndent(validationErrors, "", "  ")
	_ = os.WriteFile(filepath.Join(OutputDir, "errors.json"), errorsJSON, 0644)

	// Write run-report.json
	duration := time.Since(startTime)
	report := scraper.RunReport{
		StartTime:       startTime,
		EndTime:         time.Now().UTC(),
		DurationMs:      duration.Milliseconds(),
		CataloguePages:  pageCount,
		TotalDiscovered: len(discoveredItems),
		ValidRecords:    len(cleanRecords),
		InvalidRecords:  len(validationErrors),
		FailedPages:     failedPages,
	}

	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	_ = os.WriteFile(filepath.Join(OutputDir, "run-report.json"), reportJSON, 0644)

	slog.Info("Stage 5 complete",
		"valid_records", len(cleanRecords),
		"failed_pages", failedPages,
		"report_written", "output/run-report.json",
	)
}

// fetchWithRetry executes HTTP GET; retries once on 5xx or network errors, skips 404/403 without retrying
func fetchWithRetry(client *http.Client, targetURL string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		return body, http.StatusOK, err
	}

	if resp != nil {
		statusCode := resp.StatusCode
		resp.Body.Close()
		// Do not retry 404 Not Found or 403 Forbidden
		if statusCode == http.StatusNotFound || statusCode == http.StatusForbidden {
			return nil, statusCode, nil
		}
	}

	// Retry once for server errors (5xx) or network timeouts
	time.Sleep(1 * time.Second)
	reqRetry, _ := http.NewRequest(http.MethodGet, targetURL, nil)
	respRetry, errRetry := client.Do(reqRetry)
	if errRetry != nil {
		return nil, 0, errRetry
	}
	defer respRetry.Body.Close()

	if respRetry.StatusCode != http.StatusOK {
		return nil, respRetry.StatusCode, nil
	}

	body, err := io.ReadAll(respRetry.Body)
	return body, http.StatusOK, err
}
