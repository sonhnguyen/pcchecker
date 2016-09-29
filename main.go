package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)


type PcItem struct {
	Title       string
	Link        string
	Price       int
	Guarantee   string
	Desc 		string
	Origin		string
	Available	string
	Status 		string
	Category	string
	Image 		[]string
	Vendor		string
}

func ScrapeTanDoanh (res []PcItem) ([]PcItem, error) {
	ROOT_URL := "http://tandoanh.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/baogia.php")
	if err != nil {
		return nil, err
	}
	doc.Find("div.accordion table tbody tr").Each(func (i int, s *goquery.Selection) {
		category := s.Closest(".accordion").PrevFiltered("h3").Text()
		category = category[14:]
		if(s.Find("td:nth-child(1)").Text()!="STT"){

			price := 0
			col2 := s.Find("td:nth-child(2)").Text()
			col3 := s.Find("td:nth-child(3)").Text()
			col3 = col3[:len(col3)-3]
			col3 = strings.Replace(col3, ".", "", -1)	

			price, err = strconv.Atoi(col3)
			if err != nil {
				price = 0
			}
			col4 := s.Find("td:nth-child(4)").Text()
			col5 := s.Find("td:nth-child(5)").Text()
			item := PcItem{Title: col2, Price: price, Guarantee: col4, Available: col5, Vendor: "tandoanh", Category: category}
			res = append(res, item)
		}
	})

	return res, nil
}

func ScrapeHH(res []PcItem) ([]PcItem, error) {
	res = append(res, PcItem{Title: "hello"})
	return res, nil
}
func insertMlab(items []PcItem ) {
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
        collection := sess.DB("heroku_tr3z0r48").C("godata")
        //remove all before insert
        collection.RemoveAll(nil)

        //prepare bulk insert
        docs := make([]interface{}, len(items))
		for i := 0; i < len(items); i++ {
			docs[i] = items[i]
		}
		x := collection.Bulk()
		x.Unordered() //magic! :)
		x.Insert(docs...)
		res, err := x.Run()
		if (err!=nil) {
			panic(err)
		}

		fmt.Printf("%v", res)
}

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
	pcItems := []PcItem{}

	pcItems, err := ScrapeTanDoanh(pcItems)
	if err != nil {
		fmt.Println(err, pcItems)
	}
	fmt.Printf("length of pcitems 1: %v", len(pcItems))

	pcItems, err = ScrapeHH(pcItems)
	if err != nil {
		fmt.Println(err, pcItems)
	}

	fmt.Printf("length of pcitems 2: %v", len(pcItems))
	//insertMlab(pcItems)
	router.Run(":" + port)

}
