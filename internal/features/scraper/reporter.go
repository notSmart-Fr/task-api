package scraper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RunReport stores overall run statistics and metrics
type RunReport struct {
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	DurationMs      int64     `json:"duration_ms"`
	CataloguePages  int       `json:"catalogue_pages"`
	TotalDiscovered int       `json:"total_discovered"`
	PagesFetched    int       `json:"pages_fetched"`
	CacheHits       int       `json:"cache_hits"`
	ValidRecords    int       `json:"valid_records"`
	InvalidRecords  int       `json:"invalid_records"`
	FailedPages     int       `json:"failed_pages"`
}

type MetricsCollector struct {
	mu           sync.Mutex
	report       RunReport
	cacheHits    int
	pagesFetched int
}

func NewMetricsCollector(cataloguePages int) *MetricsCollector {
	return &MetricsCollector{
		report: RunReport{
			StartTime:      time.Now().UTC(),
			CataloguePages: cataloguePages,
		},
	}
}

func (m *MetricsCollector) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheHits++
}

func (m *MetricsCollector) RecordFetch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pagesFetched++
}

func (m *MetricsCollector) WriteReport(outputDir string, totalDiscovered, valid, invalid, failed int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.report.EndTime = time.Now().UTC()
	m.report.DurationMs = m.report.EndTime.Sub(m.report.StartTime).Milliseconds()
	m.report.TotalDiscovered = totalDiscovered
	m.report.CacheHits = m.cacheHits
	m.report.PagesFetched = m.pagesFetched
	m.report.ValidRecords = valid
	m.report.InvalidRecords = invalid
	m.report.FailedPages = failed

	_ = os.MkdirAll(outputDir, 0755)
	data, err := json.MarshalIndent(m.report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputDir, "run-report.json"), data, 0644)
}
