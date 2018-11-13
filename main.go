package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/pauldotknopf/automounter/providers"
	"github.com/pauldotknopf/automounter/providers/muxer"
	_ "github.com/pauldotknopf/automounter/providers/udevil"
	"github.com/pauldotknopf/automounter/web"
	"golang.org/x/sync/errgroup"
)

var eg errgroup.Group

func main() {
	//ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mediaProvider := muxer.Create(providers.GetProviders())

	// Start the monitoring of media.
	eg.Go(func() error {
		err := mediaProvider.Start(ctx)
		if err != nil {
			cancel()
			return err
		}
		return nil
	})

	// Start the web API.
	eg.Go(func() error {
		server := web.Create(mediaProvider)
		err := server.Listen(ctx, 3000)
		if err != nil {
			cancel()
			return err
		}
		return nil
	})

	// When quit signal comes in, stop everything
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("got signal")
		time.Sleep(time.Second * 2)
		cancel()
	}()

	err := eg.Wait()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
