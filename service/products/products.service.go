package productService

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sonhnguyen/pcchecker/mlabConnector"
	"github.com/sonhnguyen/pcchecker/model"
	"github.com/sonhnguyen/pcchecker/service/response"
	"gopkg.in/mgo.v2/bson"
)

// func GetAllProducts(c *gin.Context) {
// 	results, err := api.GetAllDocs()
// 	if err == nil {
// 		c.JSON(200, gin.H{
// 			"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
// 			"result": results})
// 	} else {
// 		c.JSON(400, gin.H{"error": responseService.ResponseError(400, err, "ERROR"), "result": nil})
// 	}
// }

func RegisterAPI(router *gin.Engine) {

	router.GET("/getProducts/:category/", GetProducts)
	router.GET("/getProducts", GetProductsV2)
	router.GET("/product/:id/", GetProductById)
}

var productCollection, _ = mlabConnector.GetCollection("products")

func GetProducts(c *gin.Context) {
	category := c.Param("category")
	var results []PcItemModel.PcItem
	productCollection.Find(bson.M{"category": category}).All(&results)
	c.JSON(200, gin.H{
		"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
		"result": results})
}

func GetProductsV2(c *gin.Context) {
	name := c.Query("name")
	vendor := c.Query("vendor")
	category := c.Query("category")
	var results []PcItemModel.PcItem

	productCollection.Find(bson.M{"title": bson.RegEx{name, ""}, "vendor": bson.RegEx{vendor, ""}, "category": bson.RegEx{category, ""}}).All(&results)

	c.JSON(200, gin.H{
		"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
		"result": results})
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")
	var result PcItemModel.PcItem
	productCollection.FindId(bson.ObjectIdHex(id)).One(&result)
	c.JSON(200, gin.H{"error": responseService.ResponseError(200, errors.New("OK"), "OK"), "result": result})
}
