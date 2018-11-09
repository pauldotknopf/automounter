package udevil

import (
	"bufio"
	"context"
	"os/exec"
	"regexp"

	"github.com/pauldotknopf/automounter/providers"
)

func init() {
	providers.AddProvider(&udevil{})
}

type udevil struct {
}

func (s *udevil) Name() string {
	return "udevil"
}

func (s *udevil) Start(ctx context.Context) error {
	cmd := exec.Command("udevil", "--monitor")

	stdout, _ := cmd.StdoutPipe()

	scanner := bufio.NewScanner(stdout)
	//scanner.Split(bufio.ScanWords)

	go func() {
		r := regexp.MustCompile(`(changed|removed|added):\s*/org/freedesktop/UDisks/devices/(.*)`)
		for scanner.Scan() {
			m := scanner.Text()
			matches := r.FindStringSubmatch(m)
			if matches == nil {
				continue
			}
			action := matches[1]
			device := "/device/" + matches[2]

			if action == "changed" {
				s.deviceChanged(device)
			} else if action == "added" {
				s.deviceAdded(device)
			} else if action == "removed" {
				s.deviceRemoved(device)
			}
		}
	}()

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		cmd.Process.Kill()
	}()

	cmd.Wait()

	return nil
}

func (s *udevil) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	result = append(result, udevilMedia{})
	return result
}

func (s *udevil) deviceAdded(device string) {

}

func (s *udevil) deviceChanged(device string) {

}

func (s *udevil) deviceRemoved(device string) {

}
