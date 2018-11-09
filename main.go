package main

import (
	_ "github.com/pauldotknopf/automounter/providers/udevil"

	"github.com/pauldotknopf/automounter/web"
)

func main() {
	server := web.Create()
	server.Listen()
}
