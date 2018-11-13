package udisks

import (
	"github.com/godbus/dbus"
)

type udisksMountSession struct {
	path      dbus.ObjectPath
	mountPath string
}

func (s *udisksMountSession) Release() error {
	return nil
}

func (s *udisksMountSession) Location() string {
	return s.mountPath
}

// func createMount(path, dbus.ObjectPath) (*udisksMountSession, error) {

// }
