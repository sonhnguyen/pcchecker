package buildService

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sonhnguyen/pcchecker/mlabConnector"
	"github.com/sonhnguyen/pcchecker/service/response"
)

type CreateBuildPostData struct {
	C  string `json:"userID" binding:"required"`
	So int    `json:"number" binding:"required"`
}

var buildCollection, _ = mlabConnector.GetCollection("build")

func CreateBuild(c *gin.Context) {
	var data CreateBuildPostData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": responseService.ResponseError(400, err, "CONNECT_ERROR"), "result": nil})
	}
	var datetimeNow time.Time
	encodedString := ""
	for len(encodedString) != 11 {
		datetimeNow = time.Now()
		randomIntWithDateNow := strconv.Itoa(rand.Intn(1e2)) + strconv.FormatInt(datetimeNow.UnixNano()/int64(time.Microsecond), 10)
		if n, err := strconv.ParseInt(randomIntWithDateNow, 10, 64); err == nil {
			encodedString = encode(n)
		}
	}
	c.JSON(200, gin.H{
		"error":  responseService.ResponseError(200, errors.New("OK"), "OK"),
		"result": encodedString})
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
