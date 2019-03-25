package udisks

type udisksMountSession struct {
	media     *udisksMedia
	mountPath string
	provider  *udisksProvider
}

func (s *udisksMountSession) Release() error {
	return s.provider.Unmount(s.media.ID())
}

func (s *udisksMountSession) Location() string {
	return s.mountPath
}
