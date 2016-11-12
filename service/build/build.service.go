package buildService

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sonhnguyen/pcchecker/mlabConnector"
	. "github.com/sonhnguyen/pcchecker/model"
	"github.com/sonhnguyen/pcchecker/service/response"
)

func RegisterAPI(router *gin.Engine) {
	router.POST("/createBuild", CreateBuild)
	router.GET("/getBuildById", GetBuildById)
	router.GET("/build/:encodedURL", GetBuildByEncodedURL)
	router.GET("/getBuildRecent", GetBuildRecent)
}

type CreateBuildPostData struct {
	Items []string `json:"items" binding:"required"`
}

var buildCollection, _ = mlabConnector.GetCollection("build")
var productCollection, _ = mlabConnector.GetCollection("products")

func CreateBuild(c *gin.Context) {
	var data CreateBuildPostData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(500, gin.H{"error": responseService.ResponseError(500, err, "CONNECT_ERROR"), "result": nil})
	} else {
		oids := make([]bson.ObjectId, len(data.Items))
		for i := range data.Items {
			if bson.IsObjectIdHex(data.Items[i]) {
				oids[i] = bson.ObjectIdHex(data.Items[i])
			} else {
				err = errors.New(data.Items[i] + " is not a valid data for CreateBuild")
			}
		}
		if err != nil {
			c.JSON(400, gin.H{"error": responseService.ResponseError(400, err, "INVALID_DATA"), "result": nil})
		} else {
			count, err := productCollection.Find(bson.M{"_id": bson.M{"$in": oids}}).Count()
			fmt.Println(count)
			if err != nil {
				c.JSON(500, gin.H{"error": responseService.ResponseError(500, err, "CONNECT_ERROR"), "result": nil})
			} else {
				if count != len(data.Items) {
					c.JSON(400, gin.H{"error": responseService.ResponseError(400, err, "INVALID_DATA"), "result": nil})
				} else {
					var datetimeNow time.Time
					encodedString := ""
					for len(encodedString) != 11 {
						datetimeNow = time.Now()
						randomIntWithDateNow := strconv.Itoa(rand.Intn(1e2)) + strconv.FormatInt(datetimeNow.UnixNano()/int64(time.Microsecond), 10)
						if n, err := strconv.ParseInt(randomIntWithDateNow, 10, 64); err == nil {
							encodedString = encode(n)
						}
					}
					err = buildCollection.Insert(&Build{Id: bson.NewObjectId(), DatetimeCreate: datetimeNow, Detail: data.Items, EncodedURL: encodedString})

					if err != nil {
						c.JSON(500, gin.H{"error": responseService.ResponseError(500, err, "CONNECT_ERROR"), "result": nil})
					} else {
						c.JSON(200, gin.H{
							"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
							"result": encodedString})
					}
				}
			}
		}
	}
}

func GetBuildById(c *gin.Context) {
	buildId := c.Query("id")
	var buildResult Build
	buildCollection.FindId(bson.ObjectIdHex(buildId)).One(&buildResult)
	var responseData BuildResponse
	responseData.Id = buildResult.Id
	responseData.DatetimeCreate = buildResult.DatetimeCreate
	responseData.By = buildResult.By
	responseData.EncodedURL = buildResult.EncodedURL

	oids := make([]bson.ObjectId, len(buildResult.Detail))
	for i := range buildResult.Detail {
		oids[i] = bson.ObjectIdHex(buildResult.Detail[i])
	}
	productCollection.Find(bson.M{"_id": bson.M{"$in": oids}}).All(&responseData.Detail)
	c.JSON(200, gin.H{
		"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
		"result": responseData})
}

func GetBuildByEncodedURL(c *gin.Context) {
	encodedURL := c.Param("encodedURL")

	var buildResult Build
	var responseData BuildResponse
	buildCollection.Find(bson.M{"encodedurl": encodedURL}).One(&buildResult)

	responseData.Id = buildResult.Id
	responseData.DatetimeCreate = buildResult.DatetimeCreate
	responseData.By = buildResult.By
	responseData.EncodedURL = buildResult.EncodedURL

	oids := make([]bson.ObjectId, len(buildResult.Detail))
	for i := range buildResult.Detail {
		oids[i] = bson.ObjectIdHex(buildResult.Detail[i])
	}
	productCollection.Find(bson.M{"_id": bson.M{"$in": oids}}).All(&responseData.Detail)
	c.JSON(200, gin.H{
		"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
		"result": responseData})

}

func GetBuildRecent(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	fmt.Println(limit)
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	var responseData []BuildResponse
	buildCollection.Find(nil).Limit(limitInt).Sort("-timestamp").All(&responseData)

	c.JSON(200, gin.H{
		"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
		"result": responseData})
}

type BuildResponse struct {
	Id             bson.ObjectId  `json:"id" bson:"_id"`
	DatetimeCreate time.Time      `json:"datetimeCreate" bson:"datetimeCreate"`
	By             *bson.ObjectId `json:"by,omitempty" bson:"by,omitempty"`
	EncodedURL     string         `json:"encodedurl" bson:"encodedurl"`
	Detail         []PcItem       `json:"detail" bson:"detail"`
}

func encode(num int64) string {
	var encoded = ""
	var alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	var base = len(alphabet) // base is the length of the alphabet (58 in this case)
	divisionResult := float64(num)
	for divisionResult > 0 {
		var remainder = int(math.Remainder(divisionResult, float64(base)))
		if remainder < 0 {
			remainder += base
		}
		divisionResult = math.Floor(divisionResult / float64(base))
		encoded = string(alphabet[remainder]) + encoded
	}
	return encoded
}
