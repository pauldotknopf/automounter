package udisks

import (
	"strconv"

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
			v := label.Value().(string)
			if len(v) > 0 {
				return v
			}
		}
		if uuid, ok := block["IdUUID"]; ok {
			v := uuid.Value().(string)
			if len(v) > 0 {
				return v
			}
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
		result["fsVersion"] = block["IdVersion"].Value().(string)
		result["size"] = strconv.FormatUint(block["Size"].Value().(uint64), 10)
		result["uuid"] = block["IdUUID"].Value().(string)
	}

	return result
}
