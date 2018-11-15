package smb

import (
	"bytes"
	"fmt"
	"strings"
)

// Options .
type Options struct {
	server   string
	share    string
	folder   string
	security string
	secure   bool
	domain   string
	username string
	password string
}

// CreateOptions .
func CreateOptions(server string, share string, folder string, security string, secure bool, domain string, username string, password string) (Options, error) {
	var result Options
	result.server = server
	result.share = share
	result.folder = folder
	result.security = security
	result.secure = secure
	result.domain = domain
	result.username = username
	result.password = password

	if len(result.server) == 0 {
		return result, fmt.Errorf("server is requierd")
	}

	if len(result.folder) == 0 {
		return result, fmt.Errorf("share is required")
	}

	if len(result.security) > 0 {
		switch security {
		case "none":
		case "krb5":
		case "krb5i":
		case "ntlm":
		case "ntlmi":
		case "ntlmv2":
		case "ntlmv2i":
		case "ntlmssp":
		case "ntlmsspi":
			break
		default:
			return result, fmt.Errorf("invalid security value")
		}
	}

	if result.secure {
		if len(username) == 0 {
			return result, fmt.Errorf("a secured connection cannot be made without a password")
		}
	} else {
		username = ""
		password = ""
		domain = ""
	}

	return result, nil
}

// MountCommand The shell command to mount these options
func (s Options) MountCommand(mountPoint string) string {
	var opts bytes.Buffer

	if s.secure {
		// escape single quotes in password character as it will be quoted in command line
		opts.WriteString(fmt.Sprintf("username='%s',", strings.Replace(s.username, "'", "\\'", -1)))
		if len(s.password) > 0 {
			opts.WriteString(fmt.Sprintf("password='%s',", strings.Replace(s.password, "'", "\\'", -1)))
		}
		if len(s.domain) > 0 {
			opts.WriteString(fmt.Sprintf("domain=%s,", s.domain))
		}
	} else {
		opts.WriteString("guest,")
	}

	if len(s.security) > 0 {
		opts.WriteString(fmt.Sprintf("sec=%s,", s.security))
	}

	opts.WriteString("rw ")

	opts.WriteString(fmt.Sprintf("//%s/%s %s", s.server, s.share, mountPoint))

	return fmt.Sprintf("sudo mount -t cifs -o %s", opts.String())
}

func (s Options) UnmountCommand(mountPoint string) string {
	return fmt.Sprintf("sudo umount %s", mountPoint)
}
