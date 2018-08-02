package main

import (
	"os"

	"github.com/therecipe/qt/widgets"
)

var qApp *widgets.QApplication

func main() {
	qApp = widgets.NewQApplication(len(os.Args), os.Args)

	// if !createConnection() {
	// 	return
	// }

	//albumDetails := core.NewQFile2("albumdetails.xml")
	window := NewMainWindow(nil, 0)
	window.initWith(nil)
	window.Show()

	qApp.Exec()
}
