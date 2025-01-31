package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"handbook-scraper/scrapers/common"
	"handbook-scraper/utils"
	"handbook-scraper/utils/databases"
	"time"
)

func GetHandbookSearchAPI(c *gin.Context, collector *colly.Collector) {

	dbHandler := databases.GetDatabaseHandler()

	// Check cache
	var cachedData string
	err := dbHandler.Retrieve(databases.Cache, "handbook_search_url", &cachedData)
	if cachedData != "" {
		c.JSON(200, gin.H{"url": cachedData})
		return
	}

	// Get the handbook search URL
	result, err := common.ExtractRawJSON("https://handbook.monash.edu/search", collector)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Navigate the result
	url := utils.GetTypedValue[string](result, "props.envConfig.API_DOMAIN")
	if url == "" {
		c.JSON(500, gin.H{"error": "could not find handbook search URL"})
		return
	}

	// Store the URL in cache
	if err := dbHandler.Store(databases.Cache, "handbook_search_url", url, time.Hour*24); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Return the URL
	c.JSON(200, gin.H{"url": url})
}
