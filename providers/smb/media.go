package smb

type smbMedia struct {
	id      string
	options Options
}

func (s *smbMedia) ID() string {
	return s.id
}

func (s *smbMedia) DisplayName() string {
	return s.options.FriendlyName()
}

func (s *smbMedia) Provider() string {
	return "smb"
}

func (s *smbMedia) Properties() map[string]string {
	result := make(map[string]string, 0)

	result["server"] = s.options.Server
	result["share"] = s.options.Share
	result["folder"] = s.options.Folder
	result["security"] = s.options.Security
	if s.options.Secure {
		result["secure"] = "true"
	} else {
		result["secure"] = "false"
	}
	result["domain"] = s.options.Domain
	result["username"] = s.options.Username
	result["password"] = s.options.Password

	return result
}
