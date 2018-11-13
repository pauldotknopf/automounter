package udevil

import (
	"github.com/pauldotknopf/automounter/providers"
)

type udevilMedia struct {
	deviceInfo   deviceInfo
	mountSession providers.MountSession
}

func (s udevilMedia) ID() string {
	return s.deviceInfo.deviceFile.file
}

func (s udevilMedia) DisplayName() string {
	return s.deviceInfo.label
}
