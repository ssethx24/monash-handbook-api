package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"handbook-scraper/scrapers/common"
	"handbook-scraper/scrapers/units"
	"net/http"
	"time"
)

func UnitCheckHandler(c *gin.Context, collector *colly.Collector) {
	year := c.Param("year")
	code := c.Param("code")

	if year == "current" {
		year = fmt.Sprintf("%d", time.Now().Year())
	}

	baseURL := fmt.Sprintf("https://handbook.monash.edu/%s/units/%s", year, code)

	data, err := ScrapeAndCache(baseURL, collector, "units")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	unitData, ok := data.(units.UnitData)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cast scraped data to UnitData"})
		return
	}

	var completedUnits []common.Unit
	if err := c.BindJSON(&completedUnits); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format for completed units"})
		return
	}

	met, unmetRequisites, err := units.CheckRequisites(unitData, completedUnits)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	enrolmentRulesString := ""
	for _, rule := range unitData.EnrolmentRules {
		enrolmentRulesString += rule.Description + " "
	}

	c.JSON(http.StatusOK, gin.H{"met_requisites": met, "message": unmetRequisites, "warning": enrolmentRulesString})
}
