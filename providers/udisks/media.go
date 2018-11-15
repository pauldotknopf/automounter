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

func (s *udisksMedia) Provider() string {
	return "udisks"
}

func (s *udisksMedia) Properties() map[string]string {
	result := make(map[string]string, 0)

	if block, ok := s.object["org.freedesktop.UDisks2.Block"]; ok {
		result["fsType"] = block["IdType"].Value().(string)
		result["size"] = string(block["Size"].Value().(uint64))
	}

	return result
}
