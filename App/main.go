// Copyright 2016 Cory Robinson. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

// main.go is the script that starts the GUI application and keeps it running.

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
