package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"task-api/internal/platform/httpclient"
)

const (
	TargetURL = "https://books.toscrape.com/catalogue/page-1.html"
	UserAgent = "FlyRankInternship-A9/1.0 (+https://github.com/yourusername/task-api)"
	Timeout   = 5 * time.Second
	CacheDir  = "cache"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	client := httpclient.NewClient(UserAgent, Timeout, CacheDir)

	req, err := http.NewRequest(http.MethodGet, TargetURL, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Request failed", "error", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Unexpected status code", "status", resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read body", "error", err)
		os.Exit(1)
	}

	slog.Info("Stage 1 complete", "url", TargetURL, "status", resp.StatusCode, "bytes", len(body))
}
