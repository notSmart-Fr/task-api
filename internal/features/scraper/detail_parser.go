package scraper

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// RawBookRecord holds uncleaned extracted DOM strings alongside provenance metadata
type RawBookRecord struct {
	Title            string    `json:"title"`
	ProductURL       string    `json:"product_url"`
	PriceText        string    `json:"price_text"`
	AvailabilityText string    `json:"availability_text"`
	RatingText       string    `json:"rating_text"`
	Description      *string   `json:"description"` // Pointer allows null in JSON if missing
	SourcePage       string    `json:"source_page"`
	FetchedAt        time.Time `json:"fetched_at"`
}

// ParseBookDetailPage extracts raw record fields from a detail page HTML body
func ParseBookDetailPage(body []byte, productURLStr string, sourcePageURL string) (*RawBookRecord, error) {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse detail HTML: %w", err)
	}

	record := &RawBookRecord{
		ProductURL: productURLStr,
		SourcePage: sourcePageURL,
		FetchedAt:  time.Now().UTC(),
	}

	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Title & Price & Rating & Availability inside <div class="col-sm-6 product_main">
			if n.Data == "div" && hasClass(n, "product_main") {
				extractProductMainInfo(n, record)
			}

			// Description inside <div id="product_description"> -> sibling <p>
			if n.Data == "p" && isDescriptionPara(n) {
				text := extractTextContent(n)
				if strings.TrimSpace(text) != "" {
					desc := strings.TrimSpace(text)
					record.Description = &desc
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(doc)
	return record, nil
}

func extractProductMainInfo(n *html.Node, record *RawBookRecord) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c.Data == "h1" {
				record.Title = strings.TrimSpace(extractTextContent(c))
			}
			if c.Data == "p" && hasClass(c, "price_color") {
				record.PriceText = strings.TrimSpace(extractTextContent(c))
			}
			if c.Data == "p" && hasClass(c, "instock") {
				record.AvailabilityText = strings.TrimSpace(extractTextContent(c))
			}
			if c.Data == "p" && hasClass(c, "star-rating") {
				record.RatingText = extractRatingClass(c)
			}
		}
	}
}

func extractRatingClass(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, class := range classes {
				if class != "star-rating" {
					return class
				}
			}
		}
	}
	return ""
}

func isDescriptionPara(n *html.Node) bool {
	// The product description text is the <p> immediately following <div id="product_description">
	if n.Parent == nil {
		return false
	}
	prev := n.PrevSibling
	for prev != nil {
		if prev.Type == html.ElementNode && prev.Data == "div" {
			for _, attr := range prev.Attr {
				if attr.Key == "id" && attr.Val == "product_description" {
					return true
				}
			}
		}
		prev = prev.PrevSibling
	}
	return false
}

func hasClass(n *html.Node, className string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, className) {
			return true
		}
	}
	return false
}

func extractTextContent(n *html.Node) string {
	var buf strings.Builder
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(n)
	return buf.String()
}
