package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/pauldotknopf/automounter/leaser"

	"github.com/pauldotknopf/automounter/providers/ios"
	"github.com/pauldotknopf/automounter/providers/muxer"
	"github.com/pauldotknopf/automounter/providers/smb"
	"github.com/pauldotknopf/automounter/providers/udisks"
	"github.com/pauldotknopf/automounter/web"
	"golang.org/x/sync/errgroup"
)

var eg errgroup.Group

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(5 * time.Second)
		//cancel()
	}()

	udisksProvider, err := udisks.Create()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	iosProvider, err := ios.Create()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	smbProvider, err := smb.Create()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	mediaProvider := muxer.Create(udisksProvider, iosProvider, smbProvider)
	leaser := leaser.Create(mediaProvider)

	// Start the processing of leases.
	eg.Go(func() error {
		leaseErr := leaser.Process(ctx)
		if leaseErr != nil {
			cancel()
			return leaseErr
		}
		return nil
	})

	// Start the monitoring of media.
	eg.Go(func() error {
		startErr := mediaProvider.Start(ctx)
		if startErr != nil {
			cancel()
			return startErr
		}
		return nil
	})

	// Start the web API.
	eg.Go(func() error {
		server := web.Create(leaser, smbProvider)
		serverErr := server.Listen(ctx, 3000)
		if serverErr != nil {
			cancel()
			return serverErr
		}
		return nil
	})

	// When quit signal comes in, stop everything
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
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
