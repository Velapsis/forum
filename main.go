package main

import (
	"main/logic"
	"main/web"
)

func main() {
	logic.Init()
	web.CreateWebsite()
}
