// Copyright 2016 Cory Robinson. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

// mainWindow.go implements the frontend layout and functionality
// of the Invoice Demo GUI.

package main

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type MainWindow struct {
	widgets.QMainWindow

	_ func() `constructor:"init"`

	_ func()                        `slot:"about"`
	_ func()                        `slot:"addInvoice"`
	_ func()                        `slot:"showAllVendorsProfile"`
	_ func(index *core.QModelIndex) `slot:"showInvoiceProfile"`
	_ func(text string)             `slot:"changeVendor"`

	tableCase string

	vendorView              *widgets.QComboBox
	invoicesAllTableView    *widgets.QTableView
	invoicesVendorTableView *widgets.QTableView
	invoicesTableView       *widgets.QTableView
	lineItemTableView       *widgets.QTableView
	tableModel              *core.QAbstractTableModel

	headerView *widgets.QHeaderView

	vendorLabel             *widgets.QLabel
	invoiceCountVendorLabel *widgets.QLabel
	invoiceDetailsLabel     *widgets.QLabel
	allVendorsLabel         *widgets.QLabel

	model Repository
}

// init() initializes the app connecting some default slots
func (w *MainWindow) init() {
	w.ConnectAbout(w.about)
	w.ConnectAddInvoice(w.addInvoice)
	w.ConnectChangeVendor(w.changeVendor)
	w.ConnectShowAllVendorsProfile(w.showAllVendorsProfile)
}

// initWith() initializes the layout views
func (w *MainWindow) initWith(parent *widgets.QWidget) {
	//w.setVendorView()
	vendor := w.createVendorGroupBox()
	invoices := w.createInvoicesGroupBox()
	details := w.createDetailsGroupBox()
	lineItems := w.createLineItemsGroupBox()
	w.showLineItemsTableView("", "", "change")

	//w.vendorView.SetCurrentIndex(0)
	w.vendorView.SetCurrentText("<all invoices>")
	w.tableCase = "all"
	w.showInvoicesTableView("<all invoices>")

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(vendor, 0, 0, 0)
	layout.AddWidget(invoices, 1, 0, 0)
	layout.AddWidget(lineItems, 2, 0, 0)
	layout.AddWidget3(details, 0, 1, 2, 1, 0)
	layout.SetColumnStretch(1, 1)
	layout.SetColumnMinimumWidth(0, 575)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	w.SetCentralWidget(widget)
	w.createMenuBar()

	w.Resize2(950, 600)
	w.SetMinimumSize2(950, 600)
	w.SetWindowTitle("Invoice Viewer (Demo)")

	w.showAllVendorsProfile()
}

// changeVendor() a slot that changes the view depending on the
// invoice that is selected in the combobox.
func (w *MainWindow) changeVendor(text string) {
	if text == "<all invoices>" {
		w.showAllVendorsProfile()
		w.tableCase = "all"
		w.showInvoicesTableView(text)
		w.showLineItemsTableView("", "", "change")
	} else {
		//name := w.vendorView.CurrentText()
		w.showVendorProfile(text)
		w.tableCase = "individual"
		w.showInvoicesTableView(text)
		w.showLineItemsTableView("", "", "change")
		//w.lineItemTableView.Reset()
	}
}

// showVendorProfile() renders the display of vendor information
// on the right hand side of the app grid.
func (w *MainWindow) showVendorProfile(name string) {
	record := w.model.GetInvoiceByString(name)[0]

	vendor := record.Vendor
	address := record.Address
	street := address.Street
	city := address.City
	state := address.State
	zipcode := address.Zipcode

	numInvoices := w.model.CountInvoicesByVendorName(vendor)

	w.vendorLabel.SetText(fmt.Sprintf("Vendor: \t%v \n\nAddress: %v\n\t%v, %v %v",
		vendor, street, city, state, zipcode))
	//w.addressLabel.SetText(fmt.Sprintf("Address: %v\n\t%v, %v %v", street, city, state, zipcode))
	w.invoiceCountVendorLabel.SetText(fmt.Sprintf("Number of Invoices: %d", numInvoices))

	w.vendorLabel.Show()
	//w.addressLabel.Show()
	w.invoiceCountVendorLabel.Show()

	w.allVendorsLabel.Hide()
	w.invoiceDetailsLabel.Hide()

	//w.invoiceDetailsLabel.Hide()
}

// showInvoiceProfile renders the display of invoice information
// on the right hand side of the app grid.
func (w *MainWindow) showInvoiceProfile(index *core.QModelIndex) {
	var record Invoice
	var invNo, vend *core.QVariant
	var statusStr, vendStr string
	//index := w.invoicesTableView.SelectionModel().CurrentIndex()
	switch w.tableCase {
	case "all":
		invNo = index.Sibling(index.Row(), 1).Data(0)
		vend = index.Sibling(index.Row(), 0).Data(0)
		record = w.model.GetInvoiceByInvoiceNoAndVendor(invNo.ToString(), vend.ToString())
		w.showVendorProfile(vend.ToString())
		w.showLineItemsTableView(invNo.ToString(), vend.ToString(), "nochange")
	case "individual":
		invNo = index.Sibling(index.Row(), 0).Data(0)
		vendStr = w.vendorView.CurrentText()
		record = w.model.GetInvoiceByInvoiceNoAndVendor(invNo.ToString(), vendStr)
		w.showLineItemsTableView(invNo.ToString(), vendStr, "nochange")
	}

	date := record.Date
	invoiceno := record.InvoiceNo // same as value.ToString()
	purchaseorder := record.PurchaseOrder
	total := record.Total
	totalStr := strconv.FormatInt(total, 10)
	totalStr = "$" + totalStr[:len(totalStr)-2] + "." + totalStr[len(totalStr)-2:]
	status := record.Paid
	if status == true {
		statusStr = "Paid"
	} else {
		statusStr = "Not Paid"
	}
	currency := record.Currency

	w.invoiceDetailsLabel.SetText(fmt.Sprintf(
		"Date: \t\t%v \nInvoice No.: \t%v \nPurchase Order: \t%v \nTotal: \t\t%v \nStatus: \t\t%v \nCurrency: \t%v",
		date, invoiceno, purchaseorder, totalStr, statusStr, currency))

	w.allVendorsLabel.Hide()

	w.vendorLabel.Show()
	w.invoiceDetailsLabel.Show()
	w.invoiceCountVendorLabel.Show()
}

// showAllVendorsProfile() renders the display of general stats
// e.g. total number of invoices, etc. on the right hand side of the app grid.
func (w *MainWindow) showAllVendorsProfile() {
	countAllVendors := w.model.CountVendors()
	countAllInvoices := w.model.RecordCount()
	countPaid := w.model.CountPaidTrue()
	countNotPaid := w.model.CountPaidFalse()

	w.allVendorsLabel.SetText(fmt.Sprintf("Vendor Count: %d \nInvoice Count: %d \nPaid/Not Paid Count: %d/%d",
		countAllVendors, countAllInvoices, countPaid, countNotPaid))

	w.allVendorsLabel.Show()

	w.vendorLabel.Hide()
	//w.addressLabel.Hide()
	w.invoiceCountVendorLabel.Hide()
	w.invoiceDetailsLabel.Hide()
}

// showInvoicesTableView() generates the table displaying invoices
// for either all vendors or each individual vendor. This is the top
// table on the left hand side of the app grid.
func (w *MainWindow) showInvoicesTableView(vendor string) {
	var table [][]string
	//var w.tableModel *core.QAbstractTableModel

	switch w.tableCase {
	case "all":
		table = w.tableForInvoicesAllTableView()

		w.tableModel = core.NewQAbstractTableModel(nil)
		w.tableModel.ConnectRowCount(func(parent *core.QModelIndex) int {
			return len(table[0]) // row-col transposed - counts cols
		})
		w.tableModel.ConnectColumnCount(func(parent *core.QModelIndex) int {
			return len(table) // row-col transposed - counts rows
		})
		w.tableModel.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
			if role != int(core.Qt__DisplayRole) {
				return core.NewQVariant()
			}
			return core.NewQVariant14(table[index.Column()][index.Row()]) // row-col transposed
		})
		w.tableModel.ConnectHeaderData(w.headerdataAll)

		//w.tableModel.Index(row, column, parent).Data(role) // see about changing paid/not paid colors.
	case "individual":
		table = w.tableForInvoicesVendorTableView(vendor)

		w.tableModel = core.NewQAbstractTableModel(nil)
		w.tableModel.ConnectRowCount(func(parent *core.QModelIndex) int {
			return len(table[0]) // row-col transposed - counts cols
		})
		w.tableModel.ConnectColumnCount(func(parent *core.QModelIndex) int {
			return len(table) // row-col transposed - counts rows
		})
		w.tableModel.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
			if role != int(core.Qt__DisplayRole) {
				return core.NewQVariant()
			}
			return core.NewQVariant14(table[index.Column()][index.Row()]) // row-col transposed
		})
		w.tableModel.ConnectHeaderData(w.headerdataIndividual)

	default:
		table = w.tableForInvoicesAllTableView()

		w.tableModel = core.NewQAbstractTableModel(nil)
		w.tableModel.ConnectRowCount(func(parent *core.QModelIndex) int {
			return len(table[0]) // row-col transposed - counts cols
		})
		w.tableModel.ConnectColumnCount(func(parent *core.QModelIndex) int {
			return len(table) // row-col transposed - counts rows
		})
		w.tableModel.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
			if role != int(core.Qt__DisplayRole) {
				return core.NewQVariant()
			}
			return core.NewQVariant14(table[index.Column()][index.Row()]) // row-col transposed
		})
		w.tableModel.ConnectHeaderData(w.headerdataAll)

	}

	w.invoicesTableView.SetModel(w.tableModel)
	//w.invoicesTableView.SetHorizontalHeader(header)
	w.adjustHeader()

	w.invoicesTableView.Show()
}

// showLineItemsTableView() generates the table displaying line items
// for the invoice selected from the top table or showInvoicesTableView.
// This is the bottom table on the left hand side of the app grid.
func (w *MainWindow) showLineItemsTableView(invoiceNo, ven, changeCase string) {
	var table [][]string
	switch changeCase {
	case "change":
		table = [][]string{{}}
	case "nochange":
		table = w.tableForLineItemsTableView(invoiceNo, ven)
	}

	Model := core.NewQAbstractTableModel(nil)
	Model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return len(table[0]) // row-col transposed - counts cols
	})
	Model.ConnectColumnCount(func(parent *core.QModelIndex) int {
		return len(table) // row-col transposed - counts rows
	})
	Model.ConnectData(func(index *core.QModelIndex, role int) *core.QVariant {
		if role != int(core.Qt__DisplayRole) {
			return core.NewQVariant()
		}
		return core.NewQVariant14(table[index.Column()][index.Row()]) // row-col transposed
	})

	switch changeCase {
	case "change":
		Model.ConnectHeaderData(w.headerdataLineItemsSelectInvoice)

	case "nochange":
		Model.ConnectHeaderData(w.headerdataLineItems)
	}

	w.lineItemTableView.SetModel(Model)
	w.adjustLineItemsHeader()

	w.lineItemTableView.Show()
}

// headerdataAll() displays the header in the showInvoicesTableView().
// NOTE: These "header" function probably could have been implemented
// better, they are kinda brute-forced, but it got the job done for now.
func (w *MainWindow) headerdataAll(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if section == 0 && orientation == core.Qt__Horizontal && role == 0 { // Qt__Horizontal, Qt__DisplayRole
		return core.NewQVariant14("Vendor")
	} else if section == 1 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Invoice No.")
	} else if section == 2 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Date")
	} else if section == 3 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Total")
	} else if section == 4 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Status")
	}
	return core.NewQVariant()
}

// headerdataIndividual() displays the header in the showInvoicesTableView().
// NOTE: These "header" function probably could have been implemented
// better, they are kinda brute-forced, but it got the job done for now.
func (w *MainWindow) headerdataIndividual(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if section == 0 && orientation == core.Qt__Horizontal && role == 0 { // Qt__Horizontal, Qt__DisplayRole
		return core.NewQVariant14("Invoice No.")
	} else if section == 1 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Date")
	} else if section == 2 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Total")
	} else if section == 3 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Status")
	}
	return core.NewQVariant()
}

// headerdataLineItems() displays the header in the showLineItemsTableView().
// NOTE: These "header" function probably could have been implemented
// better, they are kinda brute-forced, but it got the job done for now.
func (w *MainWindow) headerdataLineItems(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if section == 0 && orientation == core.Qt__Horizontal && role == 0 { // Qt__Horizontal, Qt__DisplayRole
		return core.NewQVariant14("Product ID")
	} else if section == 1 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Description")
	} else if section == 2 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Quantity")
	} else if section == 3 && orientation == core.Qt__Horizontal && role == 0 {
		return core.NewQVariant14("Amount")
	}
	return core.NewQVariant()
}

// headerdataLineItemsSelectInvoice() on the bottom table, simply displays
// 'Select An Invoice' if no invoice has been selected from the top table.
func (w *MainWindow) headerdataLineItemsSelectInvoice(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if section == 0 && orientation == core.Qt__Horizontal && role == 0 { // Qt__Horizontal, Qt__DisplayRole
		return core.NewQVariant14("Select An Invoice")
	}
	return core.NewQVariant()
}

// setVendorView() sets the vendorView model from a string array.
func (w *MainWindow) setVendorView() {
	stringList := w.model.GetInvoiceVendors()
	stringList = append([]string{"<all invoices>"}, stringList...) // prepend to stringList
	// a better prepend might be:
	//     stringList = append(stringList, "")
	//     copy(stringList[1:], stringList)
	//     stringList[0] = "<all invoices>"
	w.vendorView.SetModel(core.NewQStringListModel2(stringList, nil))
}

// createVendorGroupBox() sets up the layout for the combobox section of the
// app grid, where you can select which vendor from the combobox.
func (w *MainWindow) createVendorGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox2("Vendor", nil)

	w.vendorView = widgets.NewQComboBox(nil)

	stringList := w.model.GetInvoiceVendors()
	stringList = append([]string{"<all invoices>"}, stringList...) // prepend to stringList
	w.vendorView.SetModel(core.NewQStringListModel2(stringList, nil))

	w.vendorView.ConnectCurrentTextChanged(w.changeVendor)

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(w.vendorView, 0, 0, 0)
	box.SetLayout(layout)

	return box
}

// createInvoicesGroupBox() sets up the layout for the List of invoices table,
// i.e. the top table on the left hand side of the app grid.
func (w *MainWindow) createInvoicesGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox2("Invoices", nil)

	w.invoicesTableView = widgets.NewQTableView(nil)
	w.invoicesTableView.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	w.invoicesTableView.SetSortingEnabled(true)
	w.invoicesTableView.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)
	w.invoicesTableView.SetSelectionMode(widgets.QAbstractItemView__SingleSelection)
	w.invoicesTableView.SetShowGrid(false)
	w.invoicesTableView.VerticalHeader().Show()
	w.invoicesTableView.HorizontalHeader().Show()
	w.invoicesTableView.SetAlternatingRowColors(true)

	locale := w.invoicesTableView.Locale()
	locale.SetNumberOptions(core.QLocale__OmitGroupSeparator)
	w.invoicesTableView.SetLocale(locale)

	w.invoicesTableView.ConnectClicked(w.showInvoiceProfile)
	w.invoicesTableView.ConnectActivated(w.showInvoiceProfile)

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(w.invoicesTableView, 0, 0)
	box.SetLayout(layout)

	return box
}

// createLineItemsGroupBox() sets up the layout for the line-items table,
// i.e. the bottom table on the left hand side of the app grid.
func (w *MainWindow) createLineItemsGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox2("Line Items", nil)

	w.lineItemTableView = widgets.NewQTableView(nil)
	w.lineItemTableView.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)
	w.lineItemTableView.SetSortingEnabled(true)
	w.lineItemTableView.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)
	//w.lineItemTableView.SetSelectionMode(widgets.QAbstractItemView__NoSelection)
	w.lineItemTableView.SetSelectionMode(widgets.QAbstractItemView__SingleSelection)
	w.lineItemTableView.SetShowGrid(false)
	w.lineItemTableView.VerticalHeader().Show()
	w.lineItemTableView.HorizontalHeader().Show()
	w.lineItemTableView.SetAlternatingRowColors(true)

	locale := w.lineItemTableView.Locale()
	locale.SetNumberOptions(core.QLocale__OmitGroupSeparator)
	w.lineItemTableView.SetLocale(locale)

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(w.lineItemTableView, 0, 0)
	box.SetLayout(layout)

	return box
}

// func (w *MainWindow) hideLineItemTableColumns() {
// 	w.lineItemTableView.HideColumn(0)
// 	w.lineItemTableView.HideColumn(1)
// 	w.lineItemTableView.HideColumn(2)
// 	w.lineItemTableView.HideColumn(3)
// }

// createDetailsGroupBox() sets up the layout for the details display
// on the right hand side of the app grid.
func (w *MainWindow) createDetailsGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox2("Details", nil)

	w.allVendorsLabel = widgets.NewQLabel(nil, 0)
	w.allVendorsLabel.SetWordWrap(true)
	w.allVendorsLabel.SetAlignment(core.Qt__AlignTop)

	w.vendorLabel = widgets.NewQLabel(nil, 0)
	w.vendorLabel.SetWordWrap(false)
	w.vendorLabel.SetAlignment(core.Qt__AlignTop)

	// w.addressLabel = widgets.NewQLabel(nil, 0)
	// w.addressLabel.SetWordWrap(true)
	// w.addressLabel.SetAlignment(core.Qt__AlignBottom)

	w.invoiceCountVendorLabel = widgets.NewQLabel(nil, 0)
	w.invoiceCountVendorLabel.SetWordWrap(false)
	w.invoiceCountVendorLabel.SetAlignment(core.Qt__AlignTop)

	w.invoiceDetailsLabel = widgets.NewQLabel(nil, 0)
	w.invoiceDetailsLabel.SetWordWrap(true)
	w.invoiceDetailsLabel.SetAlignment(core.Qt__AlignTop)

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(w.allVendorsLabel, 0, 0, 0)
	layout.AddWidget(w.invoiceCountVendorLabel, 0, 1, 0)
	layout.AddWidget(w.vendorLabel, 0, 0, 0)
	layout.AddWidget(w.invoiceDetailsLabel, 1, 0, 0)
	//layout.AddWidget(w.addressLabel, 1, 0, 0)
	box.SetLayout(layout)

	return box
}

// tableForInvoicesAllTableView() queries data and sets up the table model
// for the InvoicesAllTableView.
func (w *MainWindow) tableForInvoicesAllTableView() [][]string {
	r := w.model.GetTableAllView()

	var invoiceno, vendors, dates, totalsStr, status []string
	var paid string

	for _, vens := range r {
		vendors = append(vendors, vens.Vendor)
		invoiceno = append(invoiceno, vens.InvoiceNo)
		dates = append(dates, vens.Date)
		totalStr := strconv.FormatInt(vens.Total, 10)
		totalStr = "$" + totalStr[:len(totalStr)-2] + "." + totalStr[len(totalStr)-2:]
		totalsStr = append(totalsStr, totalStr)
		if vens.Paid {
			paid = "Paid"
		} else {
			paid = "Not Paid"
		}
		status = append(status, paid)
	}

	table := [][]string{
		0: vendors,
		1: invoiceno,
		2: dates,
		3: totalsStr,
		4: status,
	}

	return table
}

// tableForInvoicesVendorTableView() queries data and sets up the table model
// for the individual vendors InvoicesVendorTableView.
func (w *MainWindow) tableForInvoicesVendorTableView(vendor string) [][]string {
	r := w.model.GetTableVendorView(vendor)

	var invoiceno, dates, totalsStr, status []string
	var paid string

	for _, vens := range r {
		invoiceno = append(invoiceno, vens.InvoiceNo)
		dates = append(dates, vens.Date)
		totalStr := strconv.FormatInt(vens.Total, 10)
		totalStr = "$" + totalStr[:len(totalStr)-2] + "." + totalStr[len(totalStr)-2:]
		totalsStr = append(totalsStr, totalStr)
		if vens.Paid {
			paid = "Paid"
		} else {
			paid = "Not Paid"
		}
		status = append(status, paid)
	}

	table := [][]string{
		0: invoiceno,
		1: dates,
		2: totalsStr,
		3: status,
	}

	return table
}

// tableForLineItemsTableView() queries data and sets up the table model
// for the individual vendor invoices line items LineItemsTableView.
func (w *MainWindow) tableForLineItemsTableView(invoiceNo, vendor string) [][]string {
	r := w.model.GetTableLineItemView(invoiceNo, vendor)

	var prodId, description, quantityStr, amountsStr []string

	for _, vens := range r {
		prodId = append(prodId, vens.ProductID)
		description = append(description, vens.Description)
		quantityStr = append(quantityStr, strconv.FormatInt(int64(vens.Quantity), 10))
		amountStr := strconv.FormatInt(vens.Amount, 10)
		amountStr = "$" + amountStr[:len(amountStr)-2] + "." + amountStr[len(amountStr)-2:]
		amountsStr = append(amountsStr, amountStr)
	}

	table := [][]string{
		0: prodId,
		1: description,
		2: quantityStr,
		3: amountsStr,
	}

	return table
}

// createMenuBar() sets up the menu bar in the main window.
func (w *MainWindow) createMenuBar() {
	addAction := widgets.NewQAction2("&Add Invoice...", w)
	quitAction := widgets.NewQAction2("&Quit", w)
	aboutAction := widgets.NewQAction2("&About", w)
	aboutQtAction := widgets.NewQAction2("About &Qt", w)

	addAction.SetShortcut(gui.QKeySequence_FromString("Ctrl+A", 0))
	quitAction.SetShortcuts2(gui.QKeySequence__Quit)

	fileMenu := w.MenuBar().AddMenu2("&File")
	fileMenu.AddActions([]*widgets.QAction{addAction})
	fileMenu.AddSeparator()
	fileMenu.AddActions([]*widgets.QAction{quitAction})

	helpMenu := w.MenuBar().AddMenu2("&Help")
	helpMenu.AddActions([]*widgets.QAction{aboutAction, aboutQtAction})

	addAction.ConnectTriggered(func(bool) { w.addInvoice() })
	quitAction.ConnectTriggered(func(bool) { w.Close() })
	aboutAction.ConnectTriggered(func(bool) { w.about() })
	aboutQtAction.ConnectTriggered(func(bool) { qApp.AboutQt() })

}

// addInvoice() slot to open the addInvoice dialog.
func (w *MainWindow) addInvoice() {
	dialog := NewDialog(nil, 0)
	//dialog.initWith(w.QWidget_PTR())
	dialog.initWith(w.QWidget_PTR())
	dialog.Exec()
	// The commented stuff below has been taken care of in the dialog code.
	// accepted := dialog.Exec()
	// fmt.Println("accepted = ", accepted)
	// if accepted == 0 {
	// 	w.setVendorView()
	// }

	//dialog.Show()
}

// adjustHeader() will adjust the table headers in the QTableViews
func (w *MainWindow) adjustHeader() {
	switch w.tableCase {
	case "all":
		//w.invoicesAllTableView.HideColumn(0)
		w.invoicesTableView.HorizontalHeader().SetSectionResizeMode2(1, widgets.QHeaderView__Stretch)
		//w.invoicesAllTableView.ResizeColumnToContents(0)
		w.invoicesTableView.ResizeColumnToContents(0)
		w.invoicesTableView.ResizeColumnToContents(2)
		w.invoicesTableView.ResizeColumnToContents(3)
		w.invoicesTableView.ResizeColumnToContents(4)
	case "individual":
		w.invoicesTableView.HorizontalHeader().SetSectionResizeMode2(0, widgets.QHeaderView__Stretch)
		w.invoicesTableView.ResizeColumnToContents(1)
		w.invoicesTableView.ResizeColumnToContents(2)
		w.invoicesTableView.ResizeColumnToContents(3)
	}
}

// adjustLineItemsHeader() will adjust the table headers in the QTableViews
func (w *MainWindow) adjustLineItemsHeader() {
	w.lineItemTableView.HorizontalHeader().SetSectionResizeMode2(0, widgets.QHeaderView__Stretch)
	w.lineItemTableView.ResizeColumnToContents(1)
	w.lineItemTableView.ResizeColumnToContents(2)
	w.lineItemTableView.ResizeColumnToContents(3)
}

// about() slot to launch dialog displaying info about the app.
func (w *MainWindow) about() {
	widgets.QMessageBox_About(w, "About Invoice Viewer",
		`<p> The <b>Invoice Viewer</b> example shows how invoices can be presented from a database. 
		The functionality includes adding, deleting, and updating invoices to a database, so as to 
		provide a fully auditable solution for storing invoice info and paying invoices.</p>`)
}
