// Copyright 2016 Cory Robinson. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

// createDummyData.go is a script that will insert dummy data into a MongoDB database.
// You should have an instance of MongoDB running and listening on the default port 27017
// before running this script. To run the script simply run the command
// `go run createDummyData.go` in your console.
//
// Feel free to modify this script to meet your needs. You can easily add more data
// with the err = c.Insert() code in the main() function below.

package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "mongodb://localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "dummyInvoice"

// COLLECTION is the name of the collection in DB
const COLLECTION = "invoice"

type Invoice struct {
	ID            int      `bson:"id"`
	Vendor        string   `json:"vendor"`
	Address       Location `json:"address"`
	LineItems     Items    `json:"items"`
	InvoiceNo     string   `json:"invoiceno"`
	Date          string   `json:"date"`
	PurchaseOrder string   `json:"purchaseorder"`
	Total         int64    `json:"total"` // hold the values as integer cents, i.e. 7420 --> $74.20
	Currency      string   `json:"currency"`
	Paid          bool     `json:"paid"`
}

type Location struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zipcode string `json:"zipcode"`
}

type Item struct {
	ProductID   string `json:"productid"`
	Description string `json:"description"`
	Quantity    uint16 `json:"quantity"`
	Amount      int64  `json:"amount"`
}

type Items []Item

func main() {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	err = c.Insert(&Invoice{1, "Right Company", Location{"123 Right Way Dr.", "Smalltown", "TX", "77336"}, Items{Item{"90d-p", "Right angle pencils", 5, 758}}, "123456", "02/24/2018", "1200364", 3790, "USD", true},
		&Invoice{2, "Wrong Company", Location{"199 Wrong Way Dr.", "Bigtown", "TX", "63377"}, Items{Item{"rw-297", "Wrong side out shirts", 7, 1499}, Item{"rw-3041", "Wrong way sttreet signs", 1, 17989}}, "15647", "03/08/2018", "1200372", 28482, "USD", true},
		&Invoice{3, "Niche Electronics", Location{"777 Electric Blvd.", "Teslatown", "IN", "77117"}, Items{Item{"el-459-h", "Electric hammers", 8, 1785}}, "143356", "12/19/2017", "1200031", 14280, "USD", false},
		&Invoice{4, "The Broken Company", Location{"173 Crooked Rd.", "Broken City", "GA", "10035"}, Items{Item{"77256103", "Broken computers", 3, 65211}}, "326679", "04/30/2018", "1200499", 195633, "USD", true},
		&Invoice{5, "Ozz", Location{"987 Yellow Brick Rd.", "Knowhere", "KS", "33665"}, Items{Item{"d-9128", "Red shoes", 1, 9999}}, "4552367", "05/01/2018", "1200506", 9999, "USD", false},
		&Invoice{6, "Niche Electronics", Location{"777 Electric Blvd.", "Teslatown", "IN", "77117"}, Items{Item{"el-376-b", "Power back scratchers", 4, 976}}, "143512", "01/27/2018", "1200126", 3904, "USD", true})
	if err != nil {
		log.Fatal(err)
	}
}
