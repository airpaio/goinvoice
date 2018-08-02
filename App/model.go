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

type Invoices []Invoice
