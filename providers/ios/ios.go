package ios

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/olebedev/emitter"
	"github.com/pauldotknopf/automounter/providers"
	"github.com/pauldotknopf/goidevice/idevice"
	"github.com/pauldotknopf/goidevice/lockdown"
)

type iosProvider struct {
	mutex sync.Mutex
	emit  *emitter.Emitter
}

// Create a media provider for iOS devices
func Create() (providers.MediaProvider, error) {
	p := &iosProvider{}
	p.emit = &emitter.Emitter{}
	p.emit.Use("*", emitter.Void)
	return p, nil
}

func (s *iosProvider) Name() string {
	return "ios"
}

func (s *iosProvider) Start(ctx context.Context) error {
	// Attach an event handler to monitor for iOS events.
	events, eventsCancel := idevice.AddEvent()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range events {
			fmt.Println(event.UUID)
			device, err := idevice.New(event.UUID)
			if err != nil {
				log.Println(err)
				continue
			}
			defer device.Close()
			uuid, err := device.UUID()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(uuid)
			lockdown, err := lockdown.NewClient(device, "lockdown")
			if err != nil {
				log.Println(err)
			}
			defer lockdown.Close()
			t, err := lockdown.Type()
			if err != nil {
				log.Println(err)
			}
			fmt.Println(t)
			err = lockdown.Pair()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	// Start raising events.
	idevice.Subscribe()

	<-ctx.Done()

	// Remove our event handler and stop monitoring for events.
	eventsCancel()
	idevice.Unsubscribe()
	wg.Wait()

	return nil
}

func (s *iosProvider) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	return result
}

func (s *iosProvider) Mount(id string) (providers.MountSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return nil, providers.ErrIDNotFound
}

func (s *iosProvider) Unmount(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return providers.ErrIDNotFound
}

func (s *iosProvider) MediaAddded() (<-chan providers.Media, func()) {
	out := make(chan providers.Media)
	in := s.emit.On("mediaAdded", func(event *emitter.Event) {
		out <- event.Args[0].(providers.Media)
	})
	cancel := func() {
		s.emit.Off("mediaAdded", in)
		close(out)
	}
	return out, cancel
}

func (s *iosProvider) MediaRemoved() (<-chan string, func()) {
	out := make(chan string)
	in := s.emit.On("mediaRemoved", func(event *emitter.Event) {
		out <- event.String(0)
	})
	cancel := func() {
		s.emit.Off("mediaRemoved", in)
		close(out)
	}
	return out, cancel
}
