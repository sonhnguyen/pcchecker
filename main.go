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

type Items []*PcItem
func (m Items) Len() int {
	return len(m)
}


func ScrapeTanDoanh() ([]*PcItem, error) {
	ROOT_URL := "http://tandoanh.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/baogia.php")
	if err != nil {
		return nil, err
	}
	res := []*PcItem{}
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
			item := &PcItem{Title: col2, Price: price, Guarantee: col4, Available: col5, Vendor: "tandoanh", Category: category}
			res = append(res, item)
		}
	})

	// catPage.Find(".title_pro_home a").Each(func(i int, s *goquery.Selection) {
	// 	iTitle := s.Text()
	// 	iLink, exists := s.Attr("href")
	// 	if !exists {
	// 		iLink = ""
	// 	} else {
	// 		iLink = "http://tandoanh.vn" + iLink 
	// 	}
	// 	item := &PcItem{Title: iTitle, Link: iLink, Category: category}
	// 	res = append(res, item)
	// })
	// catPage.Find("tbody tr b font").Each(func(i int, s *goquery.Selection) {
	// 	item :=res[i]
	// 	iPrice := s.Text()
	// 	if(iPrice != "") {
	// 		iPrice = strings.Replace(iPrice, " ", "", -1)
	// 		iPrice = strings.Replace(iPrice, ",", "", -1)
	// 	}
	// 	item.Price = iPrice
	// })

	// doc.Find(".athing").Each(func(i int, s *goquery.Selection) {
	// 	el := s.Find(".title a")
	// 	// If there's more than one found, reduce it to the first element
	// 	if el.Size() > 1 {
	// 		el = el.Slice(0, 1)
	// 	}

	// 	title := el.Text()
	// 	if err != nil {
	// 		fmt.Printf("error grabbing html: %s\n", err)
	// 		return
	// 	}
	// 	link, exists := s.Find(".title a").Attr("href")
	// 	if !exists {
	// 		link = ""
	// 	}

	// 	if strings.HasPrefix(link, "item?") {
	// 		link = ROOT_URL + "/" + link
	// 	}
	// 	item := &NewsItem{Title: title, Link: link}
	// 	res = append(res, item)
	// 	//fmt.Printf("%v - %v\n", title, link)
	// })

	// doc.Find(".subtext").Each(func(i int, s *goquery.Selection) {
	// 	pString := s.Find(".score").Text()
	// 	cString := s.Find("a").Last().Text()
	// 	cLink, exists := s.Find("a").Last().Attr("href")
	// 	if !exists {
	// 		cLink = ""
	// 	} else {
	// 		cLink = ROOT_URL + "/" + cLink
	// 	}
	// 	points := 0
	// 	comments := 0

	// 	if pString != "" {
	// 		pSt := strings.Fields(pString)[0]
	// 		points, err = strconv.Atoi(pSt)
	// 		if err != nil {
	// 			points = 0
	// 		}
	// 	}

	// 	if cString != "" && cString != "discuss" {
	// 		cSt := strings.Fields(cString)[0]
	// 		comments, err = strconv.Atoi(cSt)
	// 		if err != nil {
	// 			comments = 0
	// 		}
	// 	}

	// 	item := res[i]
	// 	item.Points = points
	// 	item.Comments = comments
	// 	item.CommentsLink = cLink
	// })

	return res, nil
}

func insertMlab(items []*PcItem ) {
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

	pcItems, err := ScrapeTanDoanh()
	if err != nil {
		fmt.Println(err)
	}
	insertMlab(pcItems)
	router.Run(":" + port)

}
