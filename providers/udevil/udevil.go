package udevil

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"sync"

	"github.com/pauldotknopf/automounter/providers"
	"golang.org/x/sync/errgroup"
)

func init() {
	providers.AddProvider(&udevil{})
}

type udevil struct {
	mutex   sync.Mutex
	devices []udevilMedia
}

func (s *udevil) Name() string {
	return "udevil"
}

func (s *udevil) Start(ctx context.Context) error {

	g, ctx := errgroup.WithContext(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g.Go(func() error {
		return monitorDevices(ctx, s.deviceAdded, s.deviceChanged, s.deviceRemoved)
	})

	// Look for the initially plugged in drives.
	pluggedIndevices, err := getPluggedInDevices()
	if err != nil {
		cancel()
		g.Wait()
		return err
	}

	for _, device := range pluggedIndevices {
		err = s.deviceAdded(device)
		if err != nil {
			cancel()
			g.Wait()
			return err
		}
	}

	return g.Wait()
}

func (s *udevil) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, media := range s.devices {
		result = append(result, media)
	}
	return result
}

func (s *udevil) deviceAdded(device string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	deviceInfo, err := getDeviceInfo(device)
	if err != nil {
		return nil
	}

	if deviceInfo.systemInternal != "1" &&
		deviceInfo.partition != nil {
		// This looks like a valid device, let's add it if it isn't already present.
		if !s.deviceExists(deviceInfo.deviceFile.file) {
			s.devices = append(s.devices, udevilMedia{deviceInfo})
		}
	}

	return nil
}

func (s *udevil) deviceChanged(device string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return nil
}

func (s *udevil) deviceRemoved(device string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	deviceInfo, err := getDeviceInfo(device)
	if err != nil {
		return err
	}

	if s.deviceExists(deviceInfo.deviceFile.file) {
		s.removeDevice(deviceInfo.deviceFile.file)
	}

	return nil
}

func getPluggedInDevices() ([]string, error) {
	b, err := ioutil.ReadFile("/proc/partitions")
	if err != nil {
		return nil, err
	}
	r := regexp.MustCompile("[ms]d[a-z0-9]*")
	matches := r.FindAllString(string(b), -1)
	if matches == nil {
		return make([]string, 0), nil
	}
	var result = make([]string, 0)
	for _, match := range matches {
		result = append(result, fmt.Sprintf("/dev/%s", match))
	}
	return result, nil
}

func (s *udevil) deviceExists(deviceFile string) bool {
	for _, device := range s.devices {
		if device.deviceInfo.deviceFile.file == deviceFile {
			return true
		}
	}
	return false
}

func (s *udevil) removeDevice(deviceFile string) {
	for i := 0; i < len(s.devices); i++ {
		if s.devices[i].deviceInfo.deviceFile.file == deviceFile {
			s.devices = append(s.devices[:i], s.devices[i+1:]...)
			i--
		}
	}
}
