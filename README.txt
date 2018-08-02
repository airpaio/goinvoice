Changes are coming. This is an initial version that I am open sourcing
as an example of building a GUI with Go. There is some code cleanup
coming as well as some better documentation and how-to build/use 
explanations and screenshots.

This application is meant to serve as a demo for automation purposes.
The app interacts with data from a MongoDB backend, and a Qt5 frontend.

NOTES:
The MongoDB server needs to be running before the app is launched.
	To start the MongoDB server (localhost) enter the following into a console:
		$(linux):sudo service mongod start
		$(windows):%MONGODPATH% mongod start

	To stop the server:
		$(linux):sudo service mongod stop
		$(windows):%MONGODPATH% mongod stop


To build the app, enter the following into a console:
	$(linux):export CGO_LDFLAGS_ALLOW=".*"
	$(linux):qtmoc desktop
	$(linux):qtrcc desktop  # currently not needed since there are no resources
	$(linux):go build -o InvoiceViewer.lex

	$(windows):qtdeploy build desktop


There exists a createDummyData.go to insert some initial data into the MongoDB. The
MongoDB server will need to be running.

There exists a play directory for messing around and getting code to work right before
putting it into the app.