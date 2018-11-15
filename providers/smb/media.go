package smb

type smbMedia struct {
	id       string
	server   string
	share    string
	folder   string
	secure   bool
	username string
	password string
}

func (s *smbMedia) ID() string {
	return s.id
}

func (s *smbMedia) DisplayName() string {
	// TODO: return network/share
	return s.id
}

func (s *smbMedia) Provider() string {
	return "smb"
}

func (s *smbMedia) Properties() map[string]string {
	result := make(map[string]string, 0)

	result["server"] = s.server
	result["share"] = s.share
	result["folder"] = s.folder
	result["username"] = s.username
	result["password"] = s.password

	return result
}
