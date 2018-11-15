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

func (s *iosMedia) Provider() string {
	return "ios"
}

func (s *iosMedia) Properties() map[string]string {
	return make(map[string]string, 0)
}
