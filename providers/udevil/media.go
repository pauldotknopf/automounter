package udevil

type udevilMedia struct {
	deviceInfo deviceInfo
}

func (s udevilMedia) ID() string {
	return "test"
}
