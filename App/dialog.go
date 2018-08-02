package main

import (
	"fmt"
	"strconv"
	"strings"
	//"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Dialog struct {
	widgets.QDialog

	_ func() `constructor:"init"`

	_ func() `slot:"reset"`
	_ func() `slot:"submit"`

	dateTimeLabel      *widgets.QLabel
	hashLabel          *widgets.QLabel
	counterIdLabel     *widgets.QLabel
	vendorLabel        *widgets.QLabel
	addressLabel       *widgets.QLabel
	lineItemsLabel     *widgets.QLabel
	invoiceNoLabel     *widgets.QLabel
	dateLabel          *widgets.QLabel
	purchaseOrderLabel *widgets.QLabel
	totalLabel         *widgets.QLabel
	currencyLabel      *widgets.QLabel

	vendorEditor        *widgets.QLineEdit
	streetEditor        *widgets.QLineEdit
	cityEditor          *widgets.QLineEdit
	stateEditor         *widgets.QLineEdit
	zipcodeEditor       *widgets.QLineEdit
	invoiceNoEditor     *widgets.QLineEdit
	dateEditor          *widgets.QLineEdit
	purchaseOrderEditor *widgets.QLineEdit
	totalEditor         *widgets.QLineEdit
	currencyEditor      *widgets.QLineEdit

	lineItemsTable *widgets.QTableWidget

	submitButton *widgets.QPushButton
	resetButton  *widgets.QPushButton
	closeButton  *widgets.QPushButton

	mwin MainWindow
}

func (d *Dialog) init() {
	d.ConnectReset(d.reset)
	d.ConnectSubmit(d.submit)
}

func (d *Dialog) initWith(parent *widgets.QWidget) {
	counterBox := d.createCounterGroupBox()
	vendorBox := d.createVendorGroupBox()
	addressBox := d.createAddressGroupBox()
	lineItemsBox := d.createLineItemsGroupBox()
	rightBox := d.createRightGroupBox()
	buttonBox := d.createButtonBox()

	d.lineItemsTable.ConnectKeyPressEvent(d.tableKeyPressEvent)

	layout := widgets.NewQGridLayout2()
	//layout.AddWidget3(widget, fromRow, fromColumn, rowSpan, columnSpan, alignment)
	layout.AddWidget(counterBox, 0, 0, 0)
	layout.AddWidget(vendorBox, 1, 0, 0)
	layout.AddWidget(addressBox, 2, 0, 0)
	layout.AddWidget(lineItemsBox, 3, 0, 0)
	layout.AddWidget3(rightBox, 0, 1, 4, 1, 0)
	layout.AddWidget3(buttonBox, 4, 0, 4, 2, 0)
	layout.SetColumnStretch(0, 1)
	layout.SetColumnMinimumWidth(0, 500)
	d.SetLayout(layout)

	d.SetWindowTitle("Add Invoice")

}

func (d *Dialog) lineItemsTableAddRow() {
	rowCount := d.lineItemsTable.RowCount()
	if d.lineItemsTable.CurrentRow() == rowCount-1 {
		d.lineItemsTable.InsertRow(rowCount)
		d.lineItemsTable.SetCurrentCell(rowCount, 0)
	}
}

func (d *Dialog) tableKeyPressEvent(e *gui.QKeyEvent) {
	rowCount := d.lineItemsTable.RowCount()
	colCount := d.lineItemsTable.ColumnCount()

	if e.Key() == int(core.Qt__Key_Enter) || e.Key() == int(core.Qt__Key_Return) { // add row
		d.lineItemsTableAddRow()
	} else if e.Key() == int(core.Qt__Key_Backtab) { // tab to prev widget
		if d.lineItemsTable.CurrentRow() == 0 && d.lineItemsTable.CurrentColumn() == 0 {
			d.FocusNextPrevChild(false)
		} else {
			d.lineItemsTable.KeyPressEventDefault(e)
		}
	} else if e.Key() == int(core.Qt__Key_Tab) { // tab to next widget
		if d.lineItemsTable.CurrentRow() == rowCount-1 &&
			d.lineItemsTable.CurrentColumn() == colCount-1 {
			d.FocusNextPrevChild(true)
		} else {
			d.lineItemsTable.KeyPressEventDefault(e)
		}
	} else { // continue with default of key press
		d.lineItemsTable.KeyPressEventDefault(e)
	}
}

//(e.Key() == int(core.Qt__Key_Tab))
func (d *Dialog) createCounterGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox(nil)

	d.counterIdLabel = widgets.NewQLabel2("7", nil, 0)
	d.dateTimeLabel = widgets.NewQLabel2("05/25/2018 7:01 AM CST", nil, 0)
	d.hashLabel = widgets.NewQLabel2("hncw98e57towg4fn", nil, 0)

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(d.dateTimeLabel, 0, 0, 0)
	layout.AddWidget(d.hashLabel, 1, 0, 0)
	layout.AddWidget(d.counterIdLabel, 0, 1, 0)
	box.SetLayout(layout)

	return box
}

func (d *Dialog) createVendorGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox2("VENDOR:", nil)

	d.vendorEditor = widgets.NewQLineEdit(nil)
	d.vendorEditor.SetPlaceholderText("Vendor Name")

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(d.vendorEditor, 0, 0, 0)
	box.SetLayout(layout)

	return box
}

func (d *Dialog) createAddressGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox2("ADDRESS:", nil)

	d.streetEditor = widgets.NewQLineEdit(nil)
	d.cityEditor = widgets.NewQLineEdit(nil)
	d.stateEditor = widgets.NewQLineEdit(nil)
	d.zipcodeEditor = widgets.NewQLineEdit(nil)

	d.stateEditor.SetMaximumWidth(22)
	d.stateEditor.SetFixedWidth(27)
	d.zipcodeEditor.SetMaximumWidth(45)
	d.zipcodeEditor.SetFixedWidth(50)

	d.streetEditor.SetPlaceholderText("123 Main St.")
	d.cityEditor.SetPlaceholderText("Vendortown")
	d.stateEditor.SetPlaceholderText("TX")
	d.zipcodeEditor.SetPlaceholderText("12345")

	layout := widgets.NewQGridLayout2()
	//layout.AddWidget3(widget, fromRow, fromColumn, rowSpan, columnSpan, alignment)
	layout.AddWidget(d.streetEditor, 0, 0, 0)
	layout.AddWidget(d.cityEditor, 1, 0, 0)
	layout.AddWidget(d.stateEditor, 1, 1, 0)
	layout.AddWidget(d.zipcodeEditor, 1, 2, 0)
	box.SetLayout(layout)

	return box
}

func (d *Dialog) createLineItemsGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox2("LINE ITEMS:", nil)

	d.lineItemsTable = widgets.NewQTableWidget2(1, 4, nil)
	d.lineItemsTable.SetHorizontalHeaderLabels(
		[]string{"Product ID", "Description", "Quantity", "Unit Price"})

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(d.lineItemsTable, 0, 0, 0)

	d.lineItemsTable.HorizontalHeader().SetSectionResizeMode2(1, widgets.QHeaderView__Stretch)
	d.lineItemsTable.ResizeColumnToContents(0)
	d.lineItemsTable.ResizeColumnToContents(1)
	d.lineItemsTable.ResizeColumnToContents(2)
	d.lineItemsTable.ResizeColumnToContents(3)

	box.SetLayout(layout)

	return box
}

func (d *Dialog) createRightGroupBox() *widgets.QGroupBox {
	box := widgets.NewQGroupBox(nil)

	d.invoiceNoLabel = widgets.NewQLabel2("INVOICE NO:", nil, 0)
	d.dateLabel = widgets.NewQLabel2("DATE:", nil, 0)
	d.purchaseOrderLabel = widgets.NewQLabel2("PURCHASE ORDER:", nil, 0)
	d.totalLabel = widgets.NewQLabel2("TOTAL:", nil, 0)
	d.currencyLabel = widgets.NewQLabel2("CURRENCY:", nil, 0)

	d.invoiceNoEditor = widgets.NewQLineEdit(nil)
	d.dateEditor = widgets.NewQLineEdit(nil)
	d.purchaseOrderEditor = widgets.NewQLineEdit(nil)
	d.totalEditor = widgets.NewQLineEdit(nil)
	d.currencyEditor = widgets.NewQLineEdit2("USD", nil)

	d.invoiceNoEditor.SetPlaceholderText("123456789")
	d.dateEditor.SetPlaceholderText("MM/DD/YYY")
	d.purchaseOrderEditor.SetPlaceholderText("ab-987654321-yz")
	d.totalEditor.SetPlaceholderText("123.45")
	d.currencyEditor.SetText("USD")

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(d.invoiceNoLabel, 0, 0, 0)
	layout.AddWidget(d.invoiceNoEditor, 0, 1, 0)
	layout.AddWidget(d.dateLabel, 1, 0, 0)
	layout.AddWidget(d.dateEditor, 1, 1, 0)
	layout.AddWidget(d.purchaseOrderLabel, 2, 0, 0)
	layout.AddWidget(d.purchaseOrderEditor, 2, 1, 0)
	layout.AddWidget(d.totalLabel, 3, 0, 0)
	layout.AddWidget(d.totalEditor, 3, 1, 0)
	layout.AddWidget(d.currencyLabel, 4, 0, 0)
	layout.AddWidget(d.currencyEditor, 4, 1, 0)
	box.SetLayout(layout)

	return box
}

func (d *Dialog) createButtonBox() *widgets.QDialogButtonBox {
	box := widgets.NewQDialogButtonBox(nil)

	d.submitButton = widgets.NewQPushButton2("&Submit", nil)
	d.resetButton = widgets.NewQPushButton2("&Reset", nil)
	d.closeButton = widgets.NewQPushButton2("&Close", nil)

	d.closeButton.SetDefault(true)

	d.closeButton.ConnectClicked(func(bool) { d.Close() })
	d.resetButton.ConnectClicked(func(bool) { d.reset() })
	d.submitButton.ConnectClicked(func(bool) { d.submit() })

	box.AddButton(d.resetButton, widgets.QDialogButtonBox__ResetRole)
	box.AddButton(d.submitButton, widgets.QDialogButtonBox__AcceptRole)
	box.AddButton(d.closeButton, widgets.QDialogButtonBox__RejectRole)

	return box
}

func (d *Dialog) submit() {
	var (
		liProductID   string
		liDescription string
		liQuantity    uint16
		//liPrice       int64
		total  int64
		paid   bool
		price  string
		stotal string

		item  Item
		items Items

		r Repository
	)

	rowCount := d.lineItemsTable.RowCount()
	for i := 0; i < rowCount; i++ {
		liProductID = d.lineItemsTable.Item(i, 0).Data(0).ToString()
		liDescription = d.lineItemsTable.Item(i, 1).Data(0).ToString()
		quant := d.lineItemsTable.Item(i, 2).Data(0).ToString()
		iquant, err := strconv.ParseInt(quant, 10, 32)
		if err != nil {
			fmt.Println("Failed to convert string to int64:", err)
		}
		liQuantity = uint16(iquant)
		price = d.lineItemsTable.Item(i, 3).Data(0).ToString()
		splitPrice := strings.Split(price, ".")
		price = strings.Join(splitPrice, "")
		liPrice, err := strconv.ParseInt(price, 10, 64)
		if err != nil {
			fmt.Println("Failed to convert string to int64: ", err)
		}

		item.ProductID = liProductID
		item.Description = liDescription
		item.Quantity = liQuantity
		item.Amount = liPrice

		items = append(items, item)
	}

	address := Location{d.streetEditor.Text(),
		d.cityEditor.Text(),
		d.stateEditor.Text(),
		d.zipcodeEditor.Text(),
	}

	stotal = d.totalEditor.Text()
	splitTotal := strings.Split(stotal, ".")
	stotal = strings.Join(splitTotal, "")
	total, err := strconv.ParseInt(stotal, 10, 64)
	if err != nil {
		fmt.Println("Failed to convert string to int64: ", err)
	}

	paid = false

	invoice := Invoice{r.maxID(),
		d.vendorEditor.Text(),
		address,
		items,
		d.invoiceNoEditor.Text(),
		d.dateEditor.Text(),
		d.purchaseOrderEditor.Text(),
		total,
		d.currencyEditor.Text(),
		paid,
	}

	// add invoice to db then update MainWindow data items and reset the dialog
	r.AddInvoice(invoice)
	d.mwin.setVendorView()
	d.reset()
	d.Accepted()
}

func (d *Dialog) reset() {
	d.vendorEditor.Clear()
	d.streetEditor.Clear()
	d.cityEditor.Clear()
	d.stateEditor.Clear()
	d.zipcodeEditor.Clear()
	d.lineItemsTable.ClearContentsDefault()
	d.lineItemsTable.SetRowCount(1)
	d.invoiceNoEditor.Clear()
	d.dateEditor.Clear()
	d.purchaseOrderEditor.Clear()
	d.totalEditor.Clear()
	d.currencyEditor.Clear()

	d.vendorEditor.SetPlaceholderText("Vendor Name")
	d.streetEditor.SetPlaceholderText("123 Main St.")
	d.cityEditor.SetPlaceholderText("Vendortown")
	d.stateEditor.SetPlaceholderText("TX")
	d.zipcodeEditor.SetPlaceholderText("12345")
	d.invoiceNoEditor.SetPlaceholderText("123456789")
	d.dateEditor.SetPlaceholderText("MM/DD/YYY")
	d.purchaseOrderEditor.SetPlaceholderText("ab-987654321-yz")
	d.totalEditor.SetPlaceholderText("123.45")
	d.currencyEditor.SetText("USD")
}
