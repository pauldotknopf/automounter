package smb

type smbMount struct {
	id        string
	mountPath string
	options   Options
	provider  *smbProvider
	isDynamic bool
}

func (s *smbMount) Release() error {
	if s.isDynamic {
		return s.unmount()
	}
	return s.provider.Unmount(s.id)
}

func (s *smbMount) Location() string {
	return s.mountPath
}
