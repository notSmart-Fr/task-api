package main

import (
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

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	client := httpclient.NewClient(UserAgent, Timeout, MinDelay, CacheDir)

	currentURL := StartURL
	pageCount := 0

	urlMap := make(map[string]bool)
	var discoveredURLs []string

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
				discoveredURLs = append(discoveredURLs, bookURL)
			}
		}

		currentURL = nextURL
	}

	slog.Info("Stage 2 complete",
		"catalogue_pages", pageCount,
		"discovered", len(discoveredURLs),
		"unique_urls", len(urlMap),
	)
}
