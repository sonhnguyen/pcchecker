package PcItemModel

import "gopkg.in/mgo.v2/bson"

// PcItem asdklajl
type PcItem struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Title     string        `json:"title" bson:"title"`
	Link      string        `json:"link" bson:"link"`
	Price     int           `json:"price" bson:"price"`
	Guarantee string        `json:"guarantee" bson:"guarantee"`
	ShortDesc string        `json:"shortDesc" bson:"shortDesc"`
	Desc      string        `json:"desc" bson:"desc"`
	Origin    string        `json:"origin" bson:"origin"`
	Available string        `json:"available" bson:"available"`
	Status    string        `json:"status" bson:"status"`
	Category  string        `json:"category" bson:"category"`
	Image     []string      `json:"image" bson:"image"`
	Vendor    string        `json:"vendor" bson:"vendor"`
}
