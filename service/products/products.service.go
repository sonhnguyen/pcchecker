package productService

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

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
	router.GET("/autocomplete", AutoComplete)
}

var productCollection, _ = mlabConnector.GetCollection("products")

func GetProducts(c *gin.Context) {
	category := c.Param("category")
	fmt.Println(category)
	var results []PcItemModel.PcItem
	productCollection.Find(bson.M{"category": category}).All(&results)
	c.JSON(200, gin.H{
		"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
		"result": results})
}

func GetProductsV2(c *gin.Context) {
	query := c.Query("query")
	// vendor := c.Query("vendor")
	// category := c.Query("category")
	// fmt.Println(name, vendor, category)
	var results []PcItemModel.PcItem

	err := productCollection.EnsureIndexKey("$text:title", "$text:category", "$text:origin", "$text:gearvn")
	if err != nil {
		c.JSON(500, gin.H{"error": responseService.ResponseError(500, err, "CONNECT_ERROR"), "result": nil})
	} else {
		productCollection.Find(
			bson.M{"$text": bson.M{"$search": query, "$caseSensitive": false, "$diacriticSensitive": false}},
		).Select(
			bson.M{"score": bson.M{"$meta": "textScore"}},
		).Sort("$textScore:score").All(&results)
		if results == nil {
			results = []PcItemModel.PcItem{}
		}
		c.JSON(200, gin.H{
			"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
			"result": results})
	}
}

func AutoComplete(c *gin.Context) {
	querystring := c.Query("querystring")
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	skip, err := strconv.Atoi(c.Query("skip"))
	if err != nil {
		skip = 10
	}
	reg1, err := regexp.Compile("[ \t]+")
	if err != nil {
		c.JSON(500, gin.H{"error": responseService.ResponseError(500, err, "CONNECT_ERROR"), "result": nil})
	} else {
		var regexString = "''"
		for _, element := range reg1.Split(querystring, -1) {
			if len(element) > 1 {
				regexString += "|" + element
			}
		}
		type autocompleteResonse struct {
			Id     bson.ObjectId `json:"id" bson:"_id"`
			Price  int           `json:"price" bson:"price"`
			Title  string        `json:"title" bson:"title"`
			Vendor string        `json:"vendor" bson:"vendor"`
		}
		var results []autocompleteResonse
		productCollection.Find(bson.M{
			"$or": []bson.M{
				bson.M{"title": bson.RegEx{regexString, "i"}},
				bson.M{"vendor": bson.RegEx{regexString, "i"}},
			},
		}).Limit(limit).Skip(skip).Select(bson.M{"price": 1, "title": 1, "vendor": 1}).All(&results)
		if results == nil {
			results = []autocompleteResonse{}
		}
		c.JSON(200, gin.H{"error": responseService.ResponseError(200, errors.New("OK"), "OK"), "result": results})
	}
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")
	var result PcItemModel.PcItem
	productCollection.FindId(bson.ObjectIdHex(id)).One(&result)
	c.JSON(200, gin.H{"error": responseService.ResponseError(200, errors.New("OK"), "OK"), "result": result})
}
