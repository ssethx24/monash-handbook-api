package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"handbook-scraper/scrapers/common"
	"handbook-scraper/server/handlers"
	"handbook-scraper/utils/databases"
	"handbook-scraper/utils/log"
)

func StartServer() {
	databases.GetDatabaseHandler()

	collector := common.SetupCollyCollector("handbook.monash.edu")
	router := SetupRouter(collector)

	log.Infof("Server started on port 8080")
	err := router.Run(":8080")
	if err != nil {
		return
	}
}

func SetupRouter(c *colly.Collector) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(corsMiddleware())

	err := router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	if err != nil {
		log.Fatal(err.Error())
	}
	SetupRoutes(router, c)
	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SetupRoutes(router *gin.Engine, collector *colly.Collector) {
	router.GET("v1/:year/units/:code", func(c *gin.Context) {
		handlers.HandbookHandler(c, collector, "units")
	})
	router.GET("v1/:year/courses/:code", func(c *gin.Context) {
		handlers.HandbookHandler(c, collector, "courses")
	})
	router.GET("v1/:year/aos/:code", func(c *gin.Context) {
		handlers.HandbookHandler(c, collector, "aos")
	})
	router.POST("v1/:year/units/:code/check", func(c *gin.Context) {
		handlers.UnitCheckHandler(c, collector)
	})
	router.GET("v1/handbook/search_url", func(c *gin.Context) {
		handlers.GetHandbookSearchAPI(c, collector)
	})
	router.GET("v1/health", handlers.HealthCheckHandler)
}
