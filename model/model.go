package PcItemModel

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type PriceToday struct {
	Datetime time.Time `json:"datetime" bson:"datetime"`
	Price    int       `json:"price" bson:"price"`
}

// PcItem
type PcItem struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Title     string        `json:"title" bson:"title"`
	Link      string        `json:"link" bson:"link"`
	Price     PriceToday    `json:"price_today" bson:"price_today"`
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

type Build struct {
	Id             bson.ObjectId  `json:"id" bson:"_id"`
	DatetimeCreate time.Time      `json:"datetimeCreate" bson:"datetimeCreate"`
	By             *bson.ObjectId `json:"by,omitempty" bson:"by,omitempty"`
	EncodedURL     string         `json:"encodedurl" bson:"encodedurl"`
	Detail         []string       `json:"detail" bson:"detail"`
}
