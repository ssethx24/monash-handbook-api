package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"handbook-scraper/scrapers/area_of_study"
	"handbook-scraper/scrapers/common"
	"handbook-scraper/scrapers/courses"
	"handbook-scraper/scrapers/units"
	"handbook-scraper/utils/databases"
	"handbook-scraper/utils/log"
	"net/http"
	"time"
)

// HandbookHandler is a generic handler for handbook data
// urlKey could be "courses", "aos", or "units"
func HandbookHandler(c *gin.Context, collector *colly.Collector, urlKey string) {

	year := c.Param("year")
	code := c.Param("code")

	if year == "current" {
		year = fmt.Sprintf("%d", time.Now().Year())
	}

	baseURL := fmt.Sprintf("https://handbook.monash.edu/%s/%s/%s", year, urlKey, code)

	log.Infof("[START] Scraping %s", baseURL)

	// Call the reusable scraping function
	final, err := ScrapeAndCache(baseURL, collector, urlKey)

	if err != nil {
		log.Errorf("[ERROR] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, final)
}

// ScrapeAndCache is a reusable function for scraping and caching data
func ScrapeAndCache(baseURL string, collector *colly.Collector, urlKey string) (interface{}, error) {

	dbHandler := databases.GetDatabaseHandler()

	// HandbookCache retrieval
	var cached interface{}
	err := dbHandler.Retrieve(databases.Handbook, baseURL, &cached)

	if cached != nil {
		log.Successf("[CACHE HIT] Success for %s", baseURL)
		return cached, nil
	}

	log.Infof("[CACHE MISS] %s", baseURL)

	// If cache miss, scrape
	data, err := common.ExtractRawJSON(baseURL, collector)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON: %w", err)
	}

	if data == nil {
		return nil, fmt.Errorf("failed to find JSON data in the HTML")
	}

	// Scrape data based on urlKey
	scraped, err := scrapeData(urlKey, data, baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape data: %w", err)
	}

	// Wrap the data and save to cache
	if err := dbHandler.Store(databases.Handbook, baseURL, scraped, time.Hour*144); err != nil {
		log.Errorf("Error saving to cache: %v", err)
	}

	log.Infof("[CACHE SAVE] %s", baseURL)

	log.Successf("[SUCCESS] Finished scraping %s", baseURL)

	return scraped, nil
}

// scrapeData handles the scraping logic based on the urlKey
func scrapeData(urlKey string, data map[string]interface{}, baseURL string) (interface{}, error) {
	switch urlKey {
	case "courses":
		return courses.Scrape(data, baseURL)
	case "aos":
		return area_of_study.Scrape(data, baseURL)
	case "units":
		return units.Scrape(data, baseURL)
	default:
		return nil, fmt.Errorf("invalid URL key: %s", urlKey)
	}
}
