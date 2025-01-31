package common

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"handbook-scraper/utils/log"
)

// SetupCollyCollector sets up a colly collector with shared error handling
func SetupCollyCollector(baseDomain string) *colly.Collector {
	log.Info("Setting up colly collector for handbook scraping")

	collector := colly.NewCollector(
		colly.AllowedDomains(baseDomain),
	)

	// Set shared error handling
	collector.OnError(func(r *colly.Response, err error) {
		log.Errorf("Request to %s failed with %v", r.Request.URL, err)
	})
	return collector
}

// ExtractRawJSON extracts raw JSON data from a URL
func ExtractRawJSON(URL string, c *colly.Collector) (map[string]interface{}, error) {
	var parsedData map[string]interface{}

	log.Logf("Extracting raw JSON data from URL: %s", URL)

	// Set the new OnHTML callback
	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		if err := json.Unmarshal([]byte(e.Text), &parsedData); err != nil {
			log.Errorf("Failed parsing JSON data: %v", err)
		}
	})

	// Start the scrape
	err := c.Visit(URL)

	// Detach the callback
	c.OnHTMLDetach("script#__NEXT_DATA__")
	if err != nil {
		return nil, fmt.Errorf("failed to visit URL: %w", err)
	}

	log.Infof("Successfully visited URL %s", URL)
	
	// Check if data is parsed
	if parsedData == nil {
		return nil, fmt.Errorf("failed to find JSON data in the HTML")
	}

	log.Log("Successfully extracted raw JSON data")
	return parsedData, nil
}
