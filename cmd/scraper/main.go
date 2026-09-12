package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
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
	MaxPages  = 3
)

type DiscoveredItem struct {
	URL        string
	SourcePage string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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

	slog.Info("Stage 2 complete",
		"catalogue_pages", pageCount,
		"discovered", len(discoveredItems),
		"unique_urls", len(urlMap),
	)

	// -------------------------------------------------------------------------
	// STAGE 3: Extract detail records
	// -------------------------------------------------------------------------
	var rawRecords []*scraper.RawBookRecord

	for _, item := range discoveredItems {
		req, err := http.NewRequest(http.MethodGet, item.URL, nil)
		if err != nil {
			slog.Error("Failed to create detail request", "url", item.URL, "error", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Detail request failed", "url", item.URL, "error", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			slog.Error("Unexpected detail status code", "url", item.URL, "status", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			slog.Error("Failed to read detail body", "url", item.URL, "error", err)
			continue
		}

		record, err := scraper.ParseBookDetailPage(body, item.URL, item.SourcePage)
		if err != nil {
			slog.Error("Failed to parse detail page", "url", item.URL, "error", err)
			continue
		}

		rawRecords = append(rawRecords, record)
	}

	slog.Info("Stage 3 complete", "detail_pages", len(rawRecords))

	// Print one sample raw record for verification
	if len(rawRecords) > 0 {
		sampleJSON, _ := json.MarshalIndent(rawRecords[0], "", "  ")
		fmt.Printf("\n--- Sample Raw Record (Stage 3 Checkpoint) ---\n%s\n---------------------------------------------\n", string(sampleJSON))
	}
}
