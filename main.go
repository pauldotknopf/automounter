package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/pauldotknopf/automounter/providers"
	"github.com/pauldotknopf/automounter/providers/muxer"
	_ "github.com/pauldotknopf/automounter/providers/udisks"
	"github.com/pauldotknopf/automounter/web"
	"golang.org/x/sync/errgroup"
)

var eg errgroup.Group

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// For testing
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	cancel()
	// }()

	mediaProvider := muxer.Create(providers.GetProviders())

	err := mediaProvider.Initialize()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Start the monitoring of media.
	eg.Go(func() error {
		startErr := mediaProvider.Start(ctx)
		if err != nil {
			cancel()
			return startErr
		}
		return nil
	})

	// Start the web API.
	eg.Go(func() error {
		server := web.Create(mediaProvider)
		serverErr := server.Listen(ctx, 3000)
		if serverErr != nil {
			cancel()
			return serverErr
		}
		return nil
	})

	// When quit signal comes in, stop everything
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	err = eg.Wait()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
