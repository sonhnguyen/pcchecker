package mlabConnector

import (
        "os"
        "fmt"
        "gopkg.in/mgo.v2/bson"      
        "gopkg.in/mgo.v2"
        . "github.com/pcchecker/model"
)

func GetMLab() ([]PcItem, error) {
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
        var results []PcItem
		err = collection.Find(bson.M{"category": "HDD / SSD"}).Sort("-timestamp").All(&results)
		if err != nil {
			panic(err)
		}
		return results, err
}
func InsertMlab(items []PcItem ) {
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
		fmt.Println("Inserting into mongodb")
		x := collection.Bulk()
		x.Unordered() //magic! :)
		x.Insert(docs...)
		res, err := x.Run()
		if (err!=nil) {
			panic(err)
		}

	fmt.Printf("done inserting into mongodb %v", res)
}
