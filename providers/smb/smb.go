package smb

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/olebedev/emitter"

	"github.com/godbus/dbus"
	"github.com/pauldotknopf/automounter/helpers"
	"github.com/pauldotknopf/automounter/providers"
)

type smbProvider struct {
	conn  *dbus.Conn
	mutex sync.Mutex
	media []*smbMedia
	emit  *emitter.Emitter
}

// Provider .
type Provider interface {
	providers.MediaProvider
	TestConnection(options Options) error
	AddMedia(options Options) (providers.Media, error)
	RemoveMedia(mediaID string) error
}

// Create a udisks block device media provider
func Create() (Provider, error) {
	p := &smbProvider{}

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	p.conn = conn
	p.emit = &emitter.Emitter{}
	p.emit.Use("*", emitter.Void)

	return p, nil
}

func (s *smbProvider) Name() string {
	return "smb"
}

func (s *smbProvider) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *smbProvider) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, media := range s.media {
		result = append(result, media)
	}
	return result
}

func (s *smbProvider) Mount(id string) (providers.MountSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return nil, providers.ErrIDNotFound
}

func (s *smbProvider) Unmount(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return providers.ErrIDNotFound
}

func (s *smbProvider) MediaAddded() (<-chan providers.Media, func()) {
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

func (s *smbProvider) MediaRemoved() (<-chan string, func()) {
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

func (s *smbProvider) TestConnection(options Options) error {
	tmpMountPath, err := helpers.GetTmpMountPath()
	if err != nil {
		return err
	}
	defer os.Remove(tmpMountPath)

	output, err := run(options.MountCommand(tmpMountPath))
	if err != nil {
		// We had an error, let's see if we can get the error from the output
		log.Println(output)
		output = extractErrorsFromMountOutput(output)
		if len(output) == 0 {
			return fmt.Errorf("could not mount")
		}
		return fmt.Errorf(output)
	}

	output, err = run(options.UnmountCommand(tmpMountPath))
	if err != nil {
		log.Println(output)
		return fmt.Errorf("couldn't unmount the test directory")
	}

	return nil
}

func (s *smbProvider) AddMedia(options Options) (providers.Media, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// First, let's see if these options resemble a media item
	// that is already present.
	for _, media := range s.media {
		if media.options.Hash == options.Hash {
			// Just act as if we added it.
			return media, nil
		}
	}

	// Add it as a new item.
	media := &smbMedia{}
	media.id = fmt.Sprintf("smb-%s", options.Hash)
	media.options = options
	s.media = append(s.media, media)
	s.emit.Emit("mediaAdded", media)

	return media, nil
}

func (s *smbProvider) RemoveMedia(mediaID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// First, let's see if these options resemble a media item
	// that is already present.
	for mediaIndex, media := range s.media {
		if media.id == mediaID {
			s.media = append(s.media[:mediaIndex], s.media[mediaIndex+1:]...)
			s.emit.Emit("mediaRemoved", mediaID)
			return nil
		}
	}

	return providers.ErrIDNotFound
}

func extractErrorsFromMountOutput(output string) string {
	var result bytes.Buffer
	regex := regexp.MustCompile(`mount error(\(.*\))?: (.*)`)
	matches := regex.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		result.WriteString(match[2])
	}
	return result.String()
}
