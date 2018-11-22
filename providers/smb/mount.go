package smb

type smbMount struct {
	id        string
	mountPath string
	options   Options
	provider  *smbProvider
}

func (s *smbMount) Release() error {
	return s.provider.Unmount(s.id)
}

func (s *smbMount) Location() string {
	return s.mountPath
}
