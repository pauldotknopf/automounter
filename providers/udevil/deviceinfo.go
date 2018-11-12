package udevil

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type deviceInfo struct {
	nativePath               string
	device                   string
	deviceFile               string
	devicePresentation       string
	deviceByID               string
	systemInternal           string
	removable                string
	hasMedia                 string
	isReadOnly               string
	isMounted                string
	mountedPaths             string
	presentationHide         string
	presentationNoPolicy     string
	presentationName         string
	presentationIcon         string
	autoMountHint            string
	size                     string
	blockSize                string
	usage                    string
	deviceType               string
	version                  string
	uuid                     string
	label                    string
	partitionSchema          string
	partitionNumber          string
	partitionType            string
	partitionFlags           string
	partitionOffset          string
	partitionAlignmentOffset string
	partitionSize            string
	partitionLabel           string
	partitionUUID            string
}

func getDeviceInfo(device string) (deviceInfo, error) {
	var result deviceInfo
	cmd := exec.Command("udevil", "--show-info", device)

	b, err := cmd.Output()
	if err != nil {
		return result, err
	}

	r := regexp.MustCompile(`([\w-\S]*):\s*([\w-\.\/:]*)`)
	reader := strings.NewReader(string(b))
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Println(scanner.Text())
		match := r.FindStringSubmatch(scanner.Text())
		if match != nil {
			_, err = reader.Seek(0, io.SeekStart)
			if err != nil {
				return result, err
			}
			scanner = bufio.NewScanner(reader)
			key := match[1]
			value := match[2]
			log.Println(key)
			log.Println(value)
			// switch key {
			// case "native-path":
			// 	result.nativePath = value
			// case "device":
			// 	result.device = value
			// case "device-file":
			// 	result.deviceFile = value
			// case "presentation":
			// 	result.devicePresentation = value
			// case "by-id":
			// 	result.deviceByID = value
			// case "system internal":
			// 	result.systemInternal = value
			// case "removable":
			// 	result.removable = value
			// case "has media":
			// 	result.hasMedia = value
			// case "is read only":
			// 	result.isReadOnly = value
			// case "is mounted":
			// 	result.isMounted = value
			// case "mount paths":
			// 	result.mountedPaths = value
			// case "presentation hide":
			// 	result.presentationHide = value
			// case "presentation nopolicy":
			// 	result.presentationNoPolicy = value
			// case "presentation name":
			// 	result.presentationName = value
			// case "automount hint":
			// 	result.autoMountHint = value
			// case "size":
			// 	result.size = value
			// case "block size":
			// 	result.blockSize = value
			// case "usage":
			// 	result.usage = value
			// case "version":
			// 	result.version = value
			// case "uuid":
			// 	result.uuid = value
			// case "label":
			// 	result.label = value
			// case "partition":
			// 	for scanner.Scan() {

			// 	}
			// case "scheme":
			// 	result.
			// default:
			// 	log.Println("unknown key " + key)
			// }
		}
	}

	// matches := r.FindAllStringSubmatch(output, -1)
	// for _, match := range matches {

	// }

	// scanner := bufio.NewScanner(strings.NewReader(output))
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }

	return result, nil
}
