package main

import (
	frm "forum/logic"
	database "forum/web/database"
)

func main() {

	println("GO: Running main.go..")

	frm.Init()

	database.DefineRequests()
	database.Connect()

	frm.InitWebsite()
}
