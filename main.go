package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
    "github.com/russross/blackfriday"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

)
type Person struct {
        Name string
        Email string
}

func mlab(c *gin.Context) {
        // Do the following:
        // In a command window:
        // set MONGOLAB_URL=mongodb://IndianGuru:dbpassword@ds051523.mongolab.com:51523/godata
        // IndianGuru is my username, replace the same with yours. Type in your password.
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

        err = collection.Insert(&Person{"Stefan Klaste", "klaste@posteo.de"},
	                        &Person{"Nishant Modak", "modak.nishant@gmail.com"},
	                        &Person{"Prathamesh Sonpatki", "csonpatki@gmail.com"},
	                        &Person{"murtuza kutub", "murtuzafirst@gmail.com"},
	                        &Person{"aniket joshi", "joshianiket22@gmail.com"},
	                        &Person{"Michael de Silva", "michael@mwdesilva.com"},
	                        &Person{"Alejandro Cespedes Vicente", "cesal_vizar@hotmail.com"})
        if err != nil {
                log.Fatal("Problem inserting data: ", err)
                return
        }

        result := Person{}
        err = collection.Find(bson.M{"name": "Prathamesh Sonpatki"}).One(&result)
        if err != nil {
                log.Fatal("Error finding record: ", err)
                return
        }

        fmt.Println("Email Id:", result.Email)

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
	router.GET("/mark", func(c *gin.Context) {
  		c.String(http.StatusOK, string(blackfriday.MarkdownBasic([]byte("**hi!**"))))
	})
    router.GET("/mongodb", mlab)

	router.Run(":" + port)
}
