package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CleanBookRecord represents a validated, normalized output record ready for storage
type CleanBookRecord struct {
	Title            string    `json:"title"`
	ProductURL       string    `json:"product_url"`
	PriceText        string    `json:"price_text"`
	PriceGBP         float64   `json:"price_gbp"`
	AvailabilityText string    `json:"availability_text"`
	RatingText       string    `json:"rating_text"`
	Description      *string   `json:"description"`
	SourcePage       string    `json:"source_page"`
	FetchedAt        time.Time `json:"fetched_at"`
}

// ValidationError holds information about why a record failed schema validation
type ValidationError struct {
	ProductURL string `json:"product_url"`
	Reason     string `json:"reason"`
}

var priceRegex = regexp.MustCompile(`[0-9]+\.?[0-9]*`)

// NormalizeAndValidate transforms raw string fields into clean typed values and validates schema
func NormalizeAndValidate(raw *RawBookRecord) (*CleanBookRecord, *ValidationError) {
	if raw == nil {
		return nil, &ValidationError{Reason: "raw record is nil"}
	}

	// 1. Validate Product URL
	if strings.TrimSpace(raw.ProductURL) == "" || !strings.HasPrefix(raw.ProductURL, "https://") {
		return nil, &ValidationError{ProductURL: raw.ProductURL, Reason: "product_url must be a valid https absolute URL"}
	}

	// 2. Validate Title
	if strings.TrimSpace(raw.Title) == "" {
		return nil, &ValidationError{ProductURL: raw.ProductURL, Reason: "title cannot be empty"}
	}

	// 3. Normalize Price Text to Float64
	match := priceRegex.FindString(raw.PriceText)
	if match == "" {
		return nil, &ValidationError{ProductURL: raw.ProductURL, Reason: fmt.Sprintf("unable to extract price float from '%s'", raw.PriceText)}
	}

	priceGBP, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return nil, &ValidationError{ProductURL: raw.ProductURL, Reason: fmt.Sprintf("failed to parse price float: %v", err)}
	}

	clean := &CleanBookRecord{
		Title:            raw.Title,
		ProductURL:       raw.ProductURL,
		PriceText:        raw.PriceText,
		PriceGBP:         priceGBP,
		AvailabilityText: raw.AvailabilityText,
		RatingText:       raw.RatingText,
		Description:      raw.Description,
		SourcePage:       raw.SourcePage,
		FetchedAt:        raw.FetchedAt,
	}

	return clean, nil
}
