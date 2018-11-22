package smb

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"
)

// Options .
type Options struct {
	Server   string
	Share    string
	Folder   string
	Security string
	Secure   bool
	Domain   string
	Username string
	Password string
	Hash     string
}

// CreateOptions .
func CreateOptions(server string, share string, folder string, security string, secure bool, domain string, username string, password string) (Options, error) {
	var result Options
	result.Server = server
	result.Share = share
	result.Folder = folder
	result.Security = security
	result.Secure = secure
	result.Domain = domain
	result.Username = username
	result.Password = password

	if len(result.Server) == 0 {
		return result, fmt.Errorf("server is requierd")
	}

	if len(result.Share) == 0 {
		return result, fmt.Errorf("share is required")
	}

	if len(result.Security) > 0 {
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

	if result.Secure {
		if len(result.Username) == 0 {
			return result, fmt.Errorf("a secured connection cannot be made without a password")
		}
	} else {
		result.Username = ""
		result.Password = ""
		result.Domain = ""
	}

	// Build a hash of all the parameters
	var hashBytes bytes.Buffer
	hashBytes.Write([]byte(result.Server))
	hashBytes.Write([]byte(result.Share))
	hashBytes.Write([]byte(result.Folder))
	hashBytes.Write([]byte(result.Security))
	if result.Secure {
		hashBytes.WriteByte(1)
	} else {
		hashBytes.WriteByte(0)
	}
	hashBytes.Write([]byte(result.Domain))
	hashBytes.Write([]byte(result.Username))
	hashBytes.Write([]byte(result.Password))

	result.Hash = fmt.Sprintf("%x", md5.Sum(hashBytes.Bytes()))

	return result, nil
}

// FriendlyName .
func (s Options) FriendlyName() string {
	return fmt.Sprintf("//%s/%s", s.Server, s.Share)
}

// MountCommand The shell command to mount these options
func (s Options) MountCommand(mountPoint string) string {
	var opts bytes.Buffer

	if s.Secure {
		// escape single quotes in password character as it will be quoted in command line
		opts.WriteString(fmt.Sprintf("username='%s',", strings.Replace(s.Username, "'", "\\'", -1)))
		if len(s.Password) > 0 {
			opts.WriteString(fmt.Sprintf("password='%s',", strings.Replace(s.Password, "'", "\\'", -1)))
		}
		if len(s.Domain) > 0 {
			opts.WriteString(fmt.Sprintf("domain=%s,", s.Domain))
		}
	} else {
		opts.WriteString("guest,")
	}

	if len(s.Security) > 0 {
		opts.WriteString(fmt.Sprintf("sec=%s,", s.Security))
	}

	opts.WriteString("rw ")

	opts.WriteString(fmt.Sprintf("//%s/%s %s", s.Server, s.Share, mountPoint))

	return fmt.Sprintf("sudo mount -t cifs -o %s", opts.String())
}

// UnmountCommand .
func (s Options) UnmountCommand(mountPoint string) string {
	return fmt.Sprintf("sudo umount %s", mountPoint)
}
