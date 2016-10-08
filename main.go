package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/crawler"
)
var pcItems = []crawler.PcItem{}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})


	router.GET("/queryAll", func(c *gin.Context) {
		pcItems, err:= crawler.GetMLab()
		fmt.Printf("length of pcitems 1: %v", pcItems)

		if err != nil {
			fmt.Println(err, pcItems)
		}
		c.JSON(200, pcItems)

	})

	router.Run(":" + port)


}
