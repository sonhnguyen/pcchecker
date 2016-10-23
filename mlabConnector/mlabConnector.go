package mlabConnector

import (
	"fmt"
	"os"

	. "github.com/sonhnguyen/pcchecker/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	collection := sess.DB("heroku_tr3z0r48").C("products")
	var results []PcItem
	err = collection.Find(bson.M{"category": "HDD / SSD"}).Sort("-timestamp").All(&results)
	if err != nil {
		panic(err)
	}
	return results, err
}

func InsertMlab(items []PcItem) {
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
	//remove all before insert
	collection.RemoveAll(nil)

	//prepare bulk insert
	docs := make([]interface{}, len(items))
	for i := 0; i < len(items); i++ {
		items[i].Id = bson.NewObjectId()
		docs[i] = items[i]
	}
	fmt.Println("Inserting into mongodb", len(docs))
	x := collection.Bulk()
	x.Unordered() //magic! :)
	x.Insert(docs...)
	fmt.Println("Inserting into mongodb")
	res, err := x.Run()
	if err != nil {
		panic(err)
	}
	fmt.Printf("done inserting into mongodb %v", res)
}
