package productService

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sonhnguyen/pcchecker/model"
	"github.com/sonhnguyen/pcchecker/service/response"
	"gopkg.in/mgo.v2"
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

	router.GET("/product/:id/", GetProductById)
}

func GetProducts(c *gin.Context) {
	category := c.Param("category")

	fmt.Println(category)
	fmt.Println("Hehehe")
	uri := os.Getenv("MONGODB_URI")
	fmt.Println(uri)
	if uri == "" {
		c.JSON(400, gin.H{"error": responseService.ResponseError(400, errors.New("Connection string error"), "CONNECT_ERROR"), "result": nil})
	} else {
		sess, err := mgo.Dial(uri)
		if err != nil {
			c.JSON(400, gin.H{"error": responseService.ResponseError(400, err, "CONNECT_ERROR"), "result": nil})
		} else {
			defer sess.Close()
			sess.SetSafe(&mgo.Safe{})
			collection := sess.DB("heroku_tr3z0r48").C("products")

			var results []PcItemModel.PcItem
			collection.Find(bson.M{"category": category}).All(&results)

			c.JSON(200, gin.H{
				"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
				"result": results})
		}
	}
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")

	fmt.Println(id)
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		c.JSON(400, gin.H{"error": responseService.ResponseError(400, errors.New("Connection string error"), "CONNECT_ERROR"), "result": nil})
	} else {
		sess, err := mgo.Dial(uri)
		if err != nil {
			c.JSON(400, gin.H{"error": responseService.ResponseError(400, err, "CONNECT_ERROR"), "result": nil})
		} else {
			defer sess.Close()
			sess.SetSafe(&mgo.Safe{})
			collection := sess.DB("heroku_tr3z0r48").C("products")

			var result PcItemModel.PcItem

			collection.FindId(bson.ObjectIdHex(id)).One(&result)

			c.JSON(200, gin.H{"error": responseService.ResponseError(200, errors.New("OK"), "OK"), "result": result})
		}
	}
}
