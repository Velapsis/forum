package main

import (
	frm "forum/logic"
	database "forum/web/database"
	"time"
)

func main() {

	time.Sleep(1 * time.Second)

	println("GO: Running main.go..")

	frm.Init()

	database.DefineRequests()
	database.Connect()

	frm.InitWebsite()
}
