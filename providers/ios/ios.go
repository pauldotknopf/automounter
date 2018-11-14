package ios

import (
	"context"
	"fmt"
	"sync"

	"github.com/olebedev/emitter"
	"github.com/pauldotknopf/automounter/providers"
	"github.com/pauldotknopf/goidevice/idevice"
	"github.com/pauldotknopf/goidevice/installation"
	"github.com/pauldotknopf/goidevice/lockdown"
	"github.com/pauldotknopf/goidevice/plist"
)

type iosProvider struct {
	mutex   sync.Mutex
	emit    *emitter.Emitter
	devices []*iosMedia
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
			switch event.Event {
			case idevice.DeviceAdded:
				s.deviceAdded(event.UUID)
				break
			case idevice.DeviceRemoved:
				s.deviceRemoved(event.UUID)
				break
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
	for _, device := range s.devices {
		result = append(result, device)
	}
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

func (s *iosProvider) deviceAdded(uuid string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device, err := idevice.New(uuid)
	if err != nil {
		return err
	}
	defer device.Close()

	lockdown, err := lockdown.NewClientWithHandshake(device, "automounter")
	if err != nil {
		return err
	}
	defer lockdown.Close()

	deviceName, err := lockdown.DeviceName()
	if err != nil {
		return err
	}

	instProxy, err := installation.NewClientStartService(device, "automounter")
	if err != nil {
		return err
	}
	defer instProxy.Close()

	// Make sure the device has our app installed before we show it.
	// NOTE: THE FOLLOWING ISN'T USEFUL TO ANYONE OTHER THAN MEDXCHANGE.

	options := plist.Create()
	defer options.Free()
	options.SetItem("ApplicationType", "User")
	returnValues := plist.CreateArray()
	defer returnValues.Free()
	returnValues.Append("CFBundleIdentifier")
	options.SetItem("ReturnAttributes", returnValues)

	apps, err := instProxy.Browse(options)
	if err != nil {
		return err
	}
	defer apps.Free()

	arraySize := apps.ArraySize()
	for i := 0; i < arraySize; i++ {
		item := apps.ArrayItem(i)
		fmt.Println(item.Type())
		bundleID := item.GetItem("CFBundleIdentifier")
		bundleIDString := bundleID.String()
		if bundleIDString == "com.medxchange.ackbar" {
			if !s.hasDevice(uuid) {
				media := &iosMedia{}
				media.deviceName = deviceName
				media.uuid = uuid
				s.devices = append(s.devices, media)
				s.emit.Emit("mediaAdded", media)
				return nil
			}
		}
	}

	return nil
}

func (s *iosProvider) deviceRemoved(uuid string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for deviceIndex, device := range s.devices {
		if device.uuid == uuid {
			s.devices = append(s.devices[:deviceIndex], s.devices[deviceIndex+1:]...)
			s.emit.Emit("mediaRemoved", uuid)
			return nil
		}
	}

	return nil
}

func (s *iosProvider) hasDevice(uuid string) bool {
	for _, device := range s.devices {
		if device.uuid == uuid {
			return true
		}
	}
	return false
}

func (s *iosProvider) removeDevice(uuid string) {
	for i := 0; i < len(s.devices); i++ {
		if s.devices[i].uuid == uuid {
			s.devices = append(s.devices[:i], s.devices[i+1:]...)
			i--
		}
	}
}
