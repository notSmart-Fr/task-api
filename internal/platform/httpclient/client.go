package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PoliteTransport wraps standard RoundTripper to inject User-Agent and timeouts
type PoliteTransport struct {
	Base      http.RoundTripper
	UserAgent string
}

func (t *PoliteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

// CacheTransport handles checking cache/ before making outbound HTTP calls
type CacheTransport struct {
	Base     http.RoundTripper
	CacheDir string
}

func (c *CacheTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet {
		return c.Base.RoundTrip(req)
	}

	cacheKey := sanitizeFilename(req.URL.Path)
	if cacheKey == "" || cacheKey == "/" {
		cacheKey = "catalogue-page-1.html"
	}
	cachePath := filepath.Join(c.CacheDir, cacheKey)

	// Check if cached file exists
	if data, err := os.ReadFile(cachePath); err == nil {
		slog.Info("CACHE HIT", "url", req.URL.String(), "size", len(data), "file", cachePath)
		return &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(bytes.NewReader(data)),
			Header:        make(http.Header),
			ContentLength: int64(len(data)),
			Request:       req,
		}, nil
	}

	// Cache Miss -> Perform actual HTTP request
	slog.Info("FETCH", "url", req.URL.String())
	resp, err := c.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		resp.Body.Close()

		_ = os.MkdirAll(c.CacheDir, 0755)
		if err := os.WriteFile(cachePath, data, 0644); err != nil {
			slog.Error("Failed to write to cache", "path", cachePath, "error", err)
		} else {
			slog.Info("SAVED TO CACHE", "file", cachePath, "size", len(data))
		}

		resp.Body = io.NopCloser(bytes.NewReader(data))
	}

	return resp, nil
}

func sanitizeFilename(path string) string {
	clean := strings.Trim(path, "/")
	clean = strings.ReplaceAll(clean, "/", "-")
	if clean == "" {
		return "catalogue-page-1.html"
	}
	if !strings.HasSuffix(clean, ".html") {
		clean += ".html"
	}
	return clean
}

// NewClient constructs a polite, caching HTTP client
func NewClient(userAgent string, timeout time.Duration, cacheDir string) *http.Client {
	polite := &PoliteTransport{
		UserAgent: userAgent,
	}

	caching := &CacheTransport{
		Base:     polite,
		CacheDir: cacheDir,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: caching,
	}
}
