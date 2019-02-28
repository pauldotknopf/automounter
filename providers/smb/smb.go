package smb

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/olebedev/emitter"

	"github.com/pauldotknopf/automounter/helpers"
	"github.com/pauldotknopf/automounter/providers"
)

type smbProvider struct {
	mutex  sync.Mutex
	media  []*smbMedia
	mounts []*smbMount
	emit   *emitter.Emitter
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	result := make([]providers.Media, 0)
	for _, media := range s.media {
		result = append(result, media)
	}
	return result
}

func (s *smbProvider) GetMediaByID(id string) providers.Media {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, media := range s.media {
		if media.ID() == id {
			return media
		}
	}
	return nil
}

func (s *smbProvider) Mount(id string) (providers.MountSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check to see if the device is already mounted
	for _, mount := range s.mounts {
		if mount.id == id {
			return &smbMount{id, mount.mountPath, mount.options, s}, nil
		}
	}

	// Look for the smb media to try to mount it
	for _, media := range s.media {
		if media.id == id {
			// We are trying to mount this smb media
			mount := &smbMount{}
			mount.id = id
			mountPath, err := helpers.GetTmpMountPath()
			if err != nil {
				return nil, err
			}
			mount.mountPath = mountPath
			mount.options = media.options
			mount.provider = s

			output, err := run(media.options.MountCommand(mountPath))
			if err != nil {
				// We couldn't mount the smb connection.
				os.RemoveAll(mountPath)
				logrus.Warningf("couldn't mount smb connection %s: %s: %+v", media.DisplayName(), output, err)
				output = extractErrorsFromMountOutput(output)
				if len(output) == 0 {
					return nil, fmt.Errorf("could not mount")
				}
				return nil, fmt.Errorf(output)
			}

			s.mounts = append(s.mounts, mount)

			return mount, nil
		}
	}

	return nil, providers.ErrIDNotFound
}

func (s *smbProvider) Unmount(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check to see if it is already mounted
	for mountIndex, mount := range s.mounts {
		if mount.id == id {
			output, err := run(mount.options.UnmountCommand(mount.mountPath))
			if err != nil {
				logrus.Errorf("couldn't unmount smb directory %s: %+v: %s", mount.mountPath, err, output)
				return fmt.Errorf("couldn't unmount smb directory")
			}

			// Since it was unmounted, let's remove the mount from our collection
			s.mounts = append(s.mounts[:mountIndex], s.mounts[mountIndex+1:]...)

			err = os.RemoveAll(mount.mountPath)
			if err != nil {
				logrus.Warnf("couldn't remove mount path %s after unmounting: %+v", mount.mountPath, err)
			}
			return nil
		}
	}

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
		logrus.Debugf("error testing mount for %s: %s: %+v", options.FriendlyName(), output, err)
		output = extractErrorsFromMountOutput(output)
		if len(output) == 0 {
			return fmt.Errorf("could not mount")
		}
		return fmt.Errorf(output)
	}

	output, err = run(options.UnmountCommand(tmpMountPath))
	if err != nil {
		logrus.Debugf("error removing mount after test for %s: %s: %+v", options.FriendlyName(), output, err)
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

	if len(mediaID) == 0 {
		return providers.ErrIDNotFound
	}

	// First, let's see if these options resemble a media item
	// that is already present.
	for mediaIndex, media := range s.media {
		if media.id == mediaID {
			s.media = append(s.media[:mediaIndex], s.media[mediaIndex+1:]...)
			s.emit.Emit("mediaRemoved", mediaID)
			// We are choosing to not unmount now,
			// since there may be mounts/leases currently in effect.
			// No harm in letting people unmount smb connections that
			// are no longer present.
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
