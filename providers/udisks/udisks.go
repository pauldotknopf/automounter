package udisks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/pauldotknopf/automounter/providers"
)

func init() {
	providers.AddProvider(&udisksProvider{})
}

type udisksProvider struct {
	conn  dbus.Conn
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

	// Get the current objects
	var result map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	o := udisks.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0)
	err := o.Store(&result)
	if err != nil {
		return err
	}
	for k := range result {
		err = s.deviceAdded(k)
		if err != nil {
			return err
		}
		// for k2, v2 := range v {
		// 	fmt.Println(k2)
		// 	for k3, v3 := range v2 {
		// 		fmt.Println(k3)
		// 		fmt.Println(v3.Value())
		// 	}
		// }
		// if strings.HasPrefix(string(k), "") {
		// 	// This is a block device
		// 	s.deviceAdded(k)
		// 	o2 := conn.Object("org.freedesktop.UDisks2", k)
		// 	node, err := introspect.Call(o2)
		// 	if err != nil {
		// 		log.Println(err)
		// 		panic(err)
		// 	}
		// 	data, _ := json.MarshalIndent(node, "", "    ")
		// 	os.Stdout.Write(data)
		// 	fmt.Println(string(k))
		// }
	}

	node, err := introspect.Call(conn.Object("org.freedesktop.DBus.ObjectManager", "/org/freedesktop/UDisks2"))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	data, _ := json.MarshalIndent(node, "", "    ")
	os.Stdout.Write(data)
	return nil
	// conn, err := dbus.SessionBus()
	// if err != nil {
	// 	return err
	// }

	// conn.BusObject()

	// conn.AddMatchSignal("/org/freedesktop/UDisks2", "object-added")

	// c := make(chan *dbus.Signal, 10)
	// conn.Signal(c)

	// for v := range c {
	// 	fmt.Println(v)
	// }

	// fmt.Println(obj)
	// return nil
}

func (s *udisksProvider) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	return result
}

func (s *udisksProvider) Mount(media providers.Media) (providers.MountSession, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return nil, fmt.Errorf("not implemented")
}

func (s *udisksProvider) deviceAdded(devicePath dbus.ObjectPath) error {
	return nil
}
