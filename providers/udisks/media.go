package udisks

import (
	"github.com/godbus/dbus"
)

type udisksMedia struct {
	path   dbus.ObjectPath
	object map[string]map[string]dbus.Variant
}

func (s *udisksMedia) ID() string {
	return string(s.path)
}

func (s *udisksMedia) DisplayName() string {
	if block, ok := s.object["org.freedesktop.UDisks2.Block"]; ok {
		if label, ok := block["IdLabel"]; ok {
			return label.Value().(string)
		}
	}
	return s.ID()
}
