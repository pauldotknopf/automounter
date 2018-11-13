package udisks

import (
	"github.com/godbus/dbus"
)

type udisksMountSession struct {
	path      dbus.ObjectPath
	mountPath string
	provider  *udisksProvider
}

func (s *udisksMountSession) Release() error {
	return s.provider.Unmount(string(s.path))
}

func (s *udisksMountSession) Location() string {
	return s.mountPath
}
