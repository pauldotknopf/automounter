package ios

type iosMountPoint struct {
	uuid     string
	path     string
	provider *iosProvider
}

func (s *iosMountPoint) Release() error {
	return s.provider.Unmount(s.uuid)
}

func (s *iosMountPoint) Location() string {
	return s.path
}
