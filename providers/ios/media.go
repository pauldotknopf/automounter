package ios

type iosMedia struct {
	uuid       string
	deviceName string
}

func (s *iosMedia) ID() string {
	return s.uuid
}

func (s *iosMedia) DisplayName() string {
	return s.deviceName
}
