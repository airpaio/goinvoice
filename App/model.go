// Copyright 2016 Cory Robinson. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

// model.go defines the model that our MongoDB data repository will follow.

package main

// Invoice represents parts of an invoice
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

// Location is a subfield containing address information.
type Location struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zipcode string `json:"zipcode"`
}

// Item is a subfield containing invoice line-item information.
type Item struct {
	ProductID   string `json:"productid"`
	Description string `json:"description"`
	Quantity    uint16 `json:"quantity"`
	Amount      int64  `json:"amount"`
}

// Items is an array of Item
type Items []Item

// Invoices is an array of Invoice
type Invoices []Invoice
