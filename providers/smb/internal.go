package smb

import (
	"fmt"
	"os"

	"github.com/pauldotknopf/automounter/helpers"
	"github.com/sirupsen/logrus"
)

func (s *smbProvider) buildMedia(options Options) *smbMedia {
	// Add it as a new item.
	media := &smbMedia{}
	media.id = fmt.Sprintf("smb-%s", options.Hash)
	media.options = options
	return media
}

func (s *smbMount) unmount() error {
	output, err := run(s.options.UnmountCommand(s.mountPath))
	if err != nil {
		logrus.Errorf("couldn't unmount smb directory %s: %+v: %s", s.mountPath, err, output)
		return fmt.Errorf("couldn't unmount smb directory")
	}

	err = os.RemoveAll(s.mountPath)
	if err != nil {
		logrus.Warnf("couldn't remove mount path %s after unmounting: %+v", s.mountPath, err)
	}

	return nil
}

func (s *smbProvider) mount(media *smbMedia) (*smbMount, error) {
	mount := &smbMount{}
	mount.id = media.ID()
	mountPath, err := helpers.GetTmpMountPath()
	if err != nil {
		return nil, err
	}
	mount.mountPath = mountPath
	mount.options = media.options
	mount.provider = s

	output, err := run(media.options.MountCommand(mountPath))
	if err != nil {
		// We couldn't mount the smb connection.
		os.RemoveAll(mountPath)
		logrus.Warningf("couldn't mount smb connection %s: %s: %+v", media.DisplayName(), output, err)
		output = extractErrorsFromMountOutput(output)
		if len(output) == 0 {
			return nil, fmt.Errorf("could not mount")
		}
		return nil, fmt.Errorf(output)
	}

	return mount, nil
}
