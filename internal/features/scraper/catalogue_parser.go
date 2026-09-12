package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ParseCataloguePage extracts book product detail URLs and the relative "next" page URL.
func ParseCataloguePage(body []byte, baseURLStr string) ([]string, string, error) {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid base URL: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var bookURLs []string
	var nextPageURL string

	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string
			var isBookLink bool
			var isNextLink bool

			// Check element hierarchy and attributes
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
				}
			}

			// Identify if link belongs to a book container (inside <h3> inside <article class="product_pod">)
			if n.Parent != nil && n.Parent.Data == "h3" {
				isBookLink = true
			}

			// Identify if link is pagination "next" button (inside <li class="next">)
			if n.Parent != nil && n.Parent.Data == "li" {
				for _, attr := range n.Parent.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "next") {
						isNextLink = true
					}
				}
			}

			if isBookLink && href != "" {
				relURL, err := url.Parse(href)
				if err == nil {
					absURL := baseURL.ResolveReference(relURL).String()
					bookURLs = append(bookURLs, absURL)
				}
			}

			if isNextLink && href != "" {
				relURL, err := url.Parse(href)
				if err == nil {
					nextPageURL = baseURL.ResolveReference(relURL).String()
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(doc)

	return bookURLs, nextPageURL, nil
}
