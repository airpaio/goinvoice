// Copyright 2016 Cory Robinson. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

// repository.go implements functionality for interacting with the
// MongoDB backend. Here, we connect to the MongoDB database, and
// implement functions with CRUD (Create-Read-Update-Delete)
// like functionality.

package main

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "mongodb://localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "dummyInvoice"

// COLLECTION is the name of the collection in DB
const COLLECTION = "invoice"

var invoiceId = 6 // TODO implement current invoiceID based on DB

// GetInvoices returns the list of whole Invoices
func (r Repository) GetInvoices() Invoices {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	results := Invoices{}

	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// GetTableAllView returns values for the invoicesAllTableView
func (r Repository) GetTableAllView() Invoices {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var results Invoices

	if err := c.Find(nil).Select(bson.M{"vendor": 1, "invoiceno": 1, "date": 1, "total": 1, "paid": 1}).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// GetTableVendorView returns values for the invoicesVendorTableView
func (r Repository) GetTableVendorView(name string) Invoices {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var results Invoices

	if err := c.Find(bson.M{"vendor": name}).Select(bson.M{"invoiceno": 1, "date": 1, "total": 1, "paid": 1}).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// GetTableLineItemView returns values for the lineItemTableView
func (r Repository) GetTableLineItemView(num, vendor string) Items {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var results Invoice

	if err := c.Find(bson.M{"invoiceno": num, "vendor": vendor}).Select(bson.M{
		"lineitems": 1}).One(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results.LineItems
}

// GetInvoiceVendors returns the list of vendors out of all of the invoices
func (r Repository) GetInvoiceVendors() []string {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var results []string

	if err := c.Find(nil).Distinct("vendor", &results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// CountVendors counts the total number of distinct vendors.
func (r Repository) CountVendors() int {
	distinctVendors := r.GetInvoiceVendors()
	count := len(distinctVendors)

	return count
}

// GetInvoiceVendorIDs returns the list of unique DB vendor IDs
func (r Repository) GetInvoiceVendorIDs() []int {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var results []int

	if err := c.Find(nil).Distinct("id", &results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// GetInvoiceById returns a unique Invoice queried by ID.
func (r Repository) GetInvoiceById(id int) Invoice {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result Invoice

	fmt.Println("ID in GetInvoiceById", id)

	if err := c.FindId(id).One(&result); err != nil {
		fmt.Println("Failed to write result:", err)
	}

	return result
}

// GetInvoiceByInvoiceNoAndVendor returns a unique Invoice.
func (r Repository) GetInvoiceByInvoiceNoAndVendor(num, vendor string) Invoice {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result Invoice

	if err := c.Find(bson.M{"invoiceno": num, "vendor": vendor}).One(&result); err != nil {
		fmt.Println("Failed to write result:", err)
	}

	return result
}

// GetInvoicesByString takes a search string as input and returns Invoices
func (r Repository) GetInvoiceByString(query string) Invoices {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var results Invoices

	// Logic to create filter
	qs := strings.Split(query, " ")
	and := make([]bson.M, len(qs))
	for i, q := range qs {
		and[i] = bson.M{"vendor": bson.M{
			"$regex": bson.RegEx{Pattern: ".*" + q + ".*", Options: "i"},
		}}
	}
	filter := bson.M{"$and": and}

	if err := c.Find(&filter).Limit(5).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// CountInvoicesByVendorName returns the number of invoices for each unique
// vendor name.
func (r Repository) CountInvoicesByVendorName(name string) int {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result int

	result, err = c.Find(bson.M{"vendor": name}).Count()
	if err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return result
}

// GetLineItemsByVendorID takes the DB vendor ID and returns a list of
// LineItems from the invoice
func (r Repository) GetLineItemsByVendorID(id int) []string {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var results []string

	if err := c.Find(bson.M{"id": id}).Distinct("lineitems", &results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

// AddInvoice adds an Invoice in the DB
func (r Repository) AddInvoice(invoice Invoice) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	invoiceId = r.incrementVendorID()
	invoice.ID = invoiceId
	session.DB(DBNAME).C(COLLECTION).Insert(invoice)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Added New Invoice ID- ", invoice.ID)

	return true
}

// UpdateInvoice updates an Invoice in the DB
func (r Repository) UpdateInvoice(invoice Invoice) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	err = session.DB(DBNAME).C(COLLECTION).UpdateId(invoice.ID, invoice)

	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("Updated Invoice ID - ", invoice.ID)

	return true
}

// DeleteInvoice deletes an Invoice by ID
func (r Repository) DeleteInvoice(id int) string {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	// Remove Invoice
	if err = session.DB(DBNAME).C(COLLECTION).RemoveId(id); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}

	fmt.Println("Deleted Invoice ID - ", id)
	// Write status
	return "OK"
}

// CountPaidTrue returns the number of paid invoices.
func (r Repository) CountPaidTrue() int {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result int

	result, err = c.Find(bson.M{"paid": true}).Count()
	if err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return result
}

// CountPaidFalse returns the number of not paid invoices
func (r Repository) CountPaidFalse() int {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result int

	result, err = c.Find(bson.M{"paid": false}).Count()
	if err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return result
}

// RecordCount returns the total number of records in the DB.
func (r Repository) RecordCount() int {
	session, err := mgo.Dial(SERVER)

	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}

	defer session.Close()

	c := session.DB(DBNAME).C(COLLECTION)
	var result int

	result, err = c.Find(nil).Count()
	if err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return result
}

// maxID returns the largest ID in the DB.
// Note: IDs are incremented sequentially, but there may be gaps in the
// numbering due to Deleted records.
func (r Repository) maxID() int {
	var ids []int

	ids = r.GetInvoiceVendorIDs()
	max := ids[0]
	for _, value := range ids {
		if value > max {
			max = value
		}
	}

	return max
}

// incrementVendorID increments the ID by one for adding new records.
// Is this implemented somewhere else?
func (r Repository) incrementVendorID() int {
	var maxID int

	maxID = r.maxID()
	maxID++

	return maxID
}
