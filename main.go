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
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

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

	router.GET("/crawTanDoanh", func(c *gin.Context) {
		c.JSON(200, nil)
		go crawler.Run()
	})
	router.GET("/getProducts/:category/", productService.GetProducts)
	router.GET("/product/:id/", productService.GetProduct)
	router.GET("/getAllDocs", productService.GetAllProducts)

	router.POST("/createBuild", buildService.CreateBuild)
	router.Run(":" + port)
}
