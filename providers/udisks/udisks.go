package udisks

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/godbus/dbus"
	"github.com/pauldotknopf/automounter/providers"
)

var (
	interfaceAdded = ""
)

func init() {
	providers.AddProvider(&udisksProvider{})
}

type udisksProvider struct {
	conn  *dbus.Conn
	mutex sync.Mutex
	media []udisksMedia
}

func (s *udisksProvider) Initialize() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *udisksProvider) Name() string {
	return "udisks"
}

func (s *udisksProvider) Start(ctx context.Context) error {

	udisks := s.conn.Object("org.freedesktop.UDisks2", "/org/freedesktop/UDisks2")

	udisks.AddMatchSignal("org.freedesktop.DBus.ObjectManager", "InterfacesAdded")
	udisks.AddMatchSignal("org.freedesktop.DBus.ObjectManager", "InterfacesRemoved")
	ch := make(chan *dbus.Signal, 5)
	s.conn.Signal(ch)

	// Before we start processing events, let's process drives that may already be plugged in.
	var result map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	o := udisks.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0)
	err := o.Store(&result)
	if err != nil {
		return err
	}

	for path := range result {
		err := s.deviceAdded(path, result[path])
		if err != nil {
			log.Println(err)
		}
	}

	go func() {
		<-ctx.Done()
		s.conn.RemoveSignal(ch)
		udisks.AddMatchSignal("org.freedesktop.DBus.ObjectManager", "InterfacesAdded")
		udisks.AddMatchSignal("org.freedesktop.DBus.ObjectManager", "InterfacesRemoved")
		close(ch)
	}()

	for {
		sig := <-ch
		if sig == nil {
			// Channel was closed
			break
		}

		path := sig.Body[0].(dbus.ObjectPath)

		switch sig.Name {
		case "org.freedesktop.DBus.ObjectManager.InterfacesAdded":
			obj, _ := sig.Body[1].(map[string]map[string]dbus.Variant)
			err = s.deviceAdded(path, obj)
			if err != nil {
				log.Println(err)
			}
			break
		case "org.freedesktop.DBus.ObjectManager.InterfacesRemoved":
			err = s.deviceRemoved(path)
			if err != nil {
				log.Println(err)
			}
			break
		}
	}

	return nil
}

func (s *udisksProvider) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, media := range s.media {
		result = append(result, &media)
	}
	return result
}

func (s *udisksProvider) Mount(media providers.Media) (providers.MountSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return nil, fmt.Errorf("not implemented")
}

func (s *udisksProvider) deviceAdded(path dbus.ObjectPath, dBusObject map[string]map[string]dbus.Variant) error {
	if _, ok := dBusObject["org.freedesktop.UDisks2.Filesystem"]; ok {
		if block, ok := dBusObject["org.freedesktop.UDisks2.Block"]; ok {
			if hintIgnore, ok := block["HintIgnore"]; ok {
				if hintIgnore.Value() == true {
					return nil
				}
			} else {
				return nil
			}
			if hintIgnore, ok := block["HintAuto"]; ok {
				if hintIgnore.Value() == true {
					// Add this device
					if !s.hasObject(path) {
						s.media = append(s.media, udisksMedia{path, dBusObject})
					}
				}
			} else {
				return nil
			}
		}
	}
	return nil
}

func (s *udisksProvider) deviceRemoved(path dbus.ObjectPath) error {
	fmt.Println(path)
	return nil
}

func (s *udisksProvider) hasObject(path dbus.ObjectPath) bool {
	for _, media := range s.media {
		if media.path == path {
			return true
		}
	}
	return false
}

func (s *udisksProvider) removeObject(path dbus.ObjectPath) {
	for i := 0; i < len(s.media); i++ {
		if s.media[i].path == path {
			s.media = append(s.media[:i], s.media[i+1:]...)
			i--
		}
	}
}

func (s *udisksProvider) getObject(path dbus.ObjectPath) *udisksMedia {
	for i := 0; i < len(s.media); i++ {
		if s.media[i].path == path {
			return &s.media[i]
		}
	}
	return nil
}
