package productService

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sonhnguyen/pcchecker/api"
	"github.com/sonhnguyen/pcchecker/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	fmt.Println(category)
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("no connection string provided")
		os.Exit(1)
	}
	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("heroku_tr3z0r48").C("products")

	var results []PcItemModel.PcItem

	collection.Find(bson.M{"category": category}).All(&results)

	c.JSON(200, results)
}
