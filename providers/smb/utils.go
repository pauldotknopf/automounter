package smb

import (
	"fmt"
	"os/exec"
)

func run(cmd string) (string, error) {
	fmt.Println(cmd)
	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return string(out), err
	}
	return "", nil
}
