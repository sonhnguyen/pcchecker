package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/sonhnguyen/pcchecker/crawler"
	"github.com/sonhnguyen/pcchecker/service/build"
	"github.com/sonhnguyen/pcchecker/service/products"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	port := os.Getenv("PORT")

	router := gin.New()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/runCrawler", func(c *gin.Context) {
		c.JSON(200, nil)
		go crawler.Run()
	})
	router.GET("/wakemydyno.txt", func(c *gin.Context) {
		c.JSON(200, "helloooooo")
	})

	productService.RegisterAPI(router)
	buildService.RegisterAPI(router)

	//router.POST("/createBuild", buildService.CreateBuild)
	router.Run(":" + port)
}
