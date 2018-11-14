package ios

import (
	"context"
	"sync"

	"github.com/olebedev/emitter"
	"github.com/pauldotknopf/automounter/providers"
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
	<-ctx.Done()
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
