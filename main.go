package main

import (
	"context"
	"sync"
	"time"

	"github.com/pauldotknopf/automounter/providers"
	"github.com/pauldotknopf/automounter/providers/muxer"
	_ "github.com/pauldotknopf/automounter/providers/udevil"
	"github.com/pauldotknopf/automounter/web"
)

var wg sync.WaitGroup

func main() {
	//ctx := context.Background()
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		//cancel()
	}()

	mediaProvider := muxer.Create(providers.GetProviders())

	// Start the monitoring of media.
	wg.Add(1)
	go func() {
		mediaProvider.Start(ctx)
		wg.Done()
	}()

	// Start the web API.
	wg.Add(1)
	go func() {
		server := web.Create(mediaProvider)
		server.Listen(ctx, 3000)
		wg.Done()
	}()

	wg.Wait()
}
