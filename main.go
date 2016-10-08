package main

import (
	"log"
	"net/http"
	"os"
    "github.com/gin-gonic/gin"
    "github.com/pcchecker/api"
)


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
	router.GET("/getAllDocs", func(c *gin.Context) {
	    result, err := api.GetAllDocs()
		if(err==nil) {
			c.JSON(200, result)
		}	else {
			c.JSON(400, err)
		}
	})
	router.Run(":" + port)

}
