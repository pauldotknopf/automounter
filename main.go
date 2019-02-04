package main

import (
	"context"
	"log"
	"os"

	"github.com/coreos/go-systemd/daemon"
	"github.com/coreos/go-systemd/journal"
	"github.com/sirupsen/logrus"
	"github.com/wercker/journalhook"

	"github.com/pauldotknopf/automounter/leaser"
	"github.com/pauldotknopf/automounter/utils/appcontext"

	"github.com/pauldotknopf/automounter/providers/ios"
	"github.com/pauldotknopf/automounter/providers/muxer"
	"github.com/pauldotknopf/automounter/providers/smb"
	"github.com/pauldotknopf/automounter/providers/udisks"
	"github.com/pauldotknopf/automounter/web"
	"golang.org/x/sync/errgroup"
)

var eg errgroup.Group

func main() {

	if journal.Enabled() {
		logrus.AddHook(&journalhook.JournalHook{})
	}

	ctx, cancel := context.WithCancel(appcontext.Context())

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
		serverErr := server.Listen(ctx, 3000, func() {
			// We have started listening for requests.
			daemon.SdNotify(false, "READY=1")
		})
		if serverErr != nil {
			cancel()
			return serverErr
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
