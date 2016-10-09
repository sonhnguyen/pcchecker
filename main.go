package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sonhnguyen/pcchecker/service/products"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/getAllDocs", productService.GetAllProducts)
	router.Run(":" + port)
}
