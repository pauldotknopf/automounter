package main

import (
	"context"
	"log"
	"sync"

	"github.com/pauldotknopf/automounter/providers"
	"github.com/pauldotknopf/automounter/providers/muxer"
	_ "github.com/pauldotknopf/automounter/providers/udevil"
	"github.com/pauldotknopf/automounter/web"
)

var wg sync.WaitGroup

func main() {

	mediaProvider := muxer.Create(providers.GetProviders())
	server := web.Create(mediaProvider)

	err := server.Listen(context.Background(), 3000)
	if err != nil {
		log.Println(err)
	}
}
