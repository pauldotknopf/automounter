package helpers

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetTmpMountPath returns a tmp path, suitable for mounting
func GetTmpMountPath() (string, error) {
	exists, err := PathExists("/run/mount")
	if err != nil {
		return "", err
	}
	if !exists {
		err = os.Mkdir("/run/mount", 0755)
		if err != nil {
			return "", err
		}
	}

	// Make our tmp directory
	path := fmt.Sprintf("/run/mount/%s", RandString(5))
	exists, err = PathExists(path)
	if err != nil {
		return "", err
	}
	for exists {
		path = fmt.Sprintf("/run/mount/%s", RandString(5))
		exists, err = PathExists(path)
		if err != nil {
			return "", err
		}
	}
	err = os.Mkdir(path, 0755)
	if err != nil {
		return "", err
	}
	return path, nil
}
