package productService

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sonhnguyen/pcchecker/api"
)

// Test asd
func Test() {
	fmt.Println("test")
}

func GetAllProducts(c *gin.Context) {
	result, err := api.GetAllDocs()
	if err == nil {
		c.JSON(200, result)
	} else {
		c.JSON(400, err)
	}
}

func GetProducts(c *gin.Context) {
	category := c.Param("category")
	name := c.Param("name")
	c.JSON(200, category+name)
}
