package udevil

type udevilMedia struct {
	deviceInfo deviceInfo
}

func (s udevilMedia) ID() string {
	return s.deviceInfo.deviceFile.file
}

func (s udevilMedia) DisplayName() string {
	return s.deviceInfo.label
}
