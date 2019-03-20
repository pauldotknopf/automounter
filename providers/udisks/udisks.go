package udisks

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/olebedev/emitter"

	"github.com/godbus/dbus"
	"github.com/pauldotknopf/automounter/providers"
)

type udisksProvider struct {
	conn  *dbus.Conn
	mutex sync.Mutex
	media []*udisksMedia
	emit  *emitter.Emitter
}

// Create a udisks block device media provider
func Create() (providers.MediaProvider, error) {
	p := &udisksProvider{}

	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	p.conn = conn
	p.emit = &emitter.Emitter{}
	p.emit.Use("*", emitter.Void)

	return p, nil
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	result := make([]providers.Media, 0)
	for _, media := range s.media {
		result = append(result, media)
	}
	return result
}

func (s *udisksProvider) GetMediaByID(id string) providers.Media {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, media := range s.media {
		if media.ID() == id {
			return media
		}
	}
	return nil
}

func (s *udisksProvider) Mount(id string) (providers.MountSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, media := range s.media {
		if media.ID() == id {
			obj := s.conn.Object("org.freedesktop.UDisks2", media.path)
			var params map[string]dbus.Variant
			var location string
			err := obj.Call("org.freedesktop.UDisks2.Filesystem.Mount", 0, params).Store(&location)
			if err != nil {
				if dbusError, ok := err.(dbus.Error); ok {
					if dbusError.Name == "org.freedesktop.UDisks2.Error.AlreadyMounted" {
						v, err := getPropertyStringArray(s.conn, media.path, "org.freedesktop.UDisks2.Filesystem.MountPoints")
						if err != nil {
							return nil, err
						}
						if len(v) == 0 {
							return nil, fmt.Errorf("mount indicated it was already mounted, but couldn't find the mount")
						}
						session := &udisksMountSession{}
						session.path = media.path
						session.mountPath = v[0]
						session.provider = s
						return session, nil
					}
				}
				return nil, err
			}
			session := &udisksMountSession{}
			session.path = media.path
			session.mountPath = location
			session.provider = s
			return session, nil
		}
	}

	return nil, providers.ErrIDNotFound
}

func (s *udisksProvider) Unmount(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, media := range s.media {
		if media.ID() == id {
			obj := s.conn.Object("org.freedesktop.UDisks2", media.path)
			var params = map[string]dbus.Variant{
				"force": dbus.MakeVariant(true),
			}
			err := obj.Call("org.freedesktop.UDisks2.Filesystem.Unmount", 0, params).Store()
			if err != nil {
				if dbusError, ok := err.(dbus.Error); ok {
					if dbusError.Name == "org.freedesktop.UDisks2.Error.NotMounted" {
						return nil
					}
				}
				return err
			}
			return nil
		}
	}

	return providers.ErrIDNotFound
}

func (s *udisksProvider) MediaAddded() (<-chan providers.Media, func()) {
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

func (s *udisksProvider) MediaRemoved() (<-chan string, func()) {
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

func (s *udisksProvider) deviceAdded(path dbus.ObjectPath, dBusObject map[string]map[string]dbus.Variant) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

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
						m := &udisksMedia{path, dBusObject}
						s.media = append(s.media, m)
						s.emit.Emit("mediaAdded", m)
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for mediaIndex, media := range s.media {
		if media.path == path {
			s.media = append(s.media[:mediaIndex], s.media[mediaIndex+1:]...)
			s.emit.Emit("mediaRemoved", string(path))
			return nil
		}
	}

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
			return s.media[i]
		}
	}
	return nil
}
