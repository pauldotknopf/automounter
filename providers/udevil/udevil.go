package udevil

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"sync"

	"github.com/pauldotknopf/automounter/providers"
)

func init() {
	providers.AddProvider(&udevil{})
}

type udevil struct {
	mutex   sync.Mutex
	devices []udevilMedia
}

func (s *udevil) Name() string {
	return "udevil"
}

func (s *udevil) Start(ctx context.Context) error {
	cmd := exec.Command("udevil", "--monitor")

	stdout, _ := cmd.StdoutPipe()

	scanner := bufio.NewScanner(stdout)

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

	pluggedIndevices, err := getPluggedInDevices()
	if err != nil {
		cmd.Process.Kill()
		return err
	}

	for _, device := range pluggedIndevices {
		s.deviceAdded(device)
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	deviceInfo, err := getDeviceInfo(device)
	if err != nil {
		log.Println(err)
		return
	}

	s.devices = append(s.devices, udevilMedia{deviceInfo})
}

func (s *udevil) deviceChanged(device string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
}

func (s *udevil) deviceRemoved(device string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
}

func getPluggedInDevices() ([]string, error) {
	b, err := ioutil.ReadFile("/proc/partitions")
	if err != nil {
		return nil, err
	}
	r := regexp.MustCompile("[ms]d[a-z0-9]*")
	matches := r.FindAllString(string(b), -1)
	if matches == nil {
		return make([]string, 0), nil
	}
	var result = make([]string, 0)
	for _, match := range matches {
		result = append(result, fmt.Sprintf("/dev/%s", match))
	}
	return result, nil
}
