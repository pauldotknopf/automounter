package udevil

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type deviceInfoPartitionTable struct {
	partitionTable string
	schema         string
	count          string
}

type deviceInfoDrive struct {
	drive      string
	vendor     string
	model      string
	revision   string
	serial     string
	WWN        string
	detachable string
	ejectable  string
	media      *deviceInfoDriveMedia
	inter      string
	ifSpeed    string
}

type deviceInfoDriveMedia struct {
	media  string
	compat string
}

type deviceInfoFile struct {
	file         string
	presentation string
	byID         string
}

type deviceInfoPartition struct {
	partition       string
	scheme          string
	number          string
	partitionType   string
	flags           string
	offset          string
	alignmentOffset string
	size            string
	label           string
	uuid            string
}

type deviceInfo struct {
	nativePath           string
	device               string
	deviceFile           *deviceInfoFile
	systemInternal       string
	removable            string
	hasMedia             string
	isReadOnly           string
	isMounted            string
	mountedPaths         string
	presentationHide     string
	presentationNoPolicy string
	presentationName     string
	presentationIcon     string
	autoMountHint        string
	size                 string
	blockSize            string
	usage                string
	deviceType           string
	version              string
	uuid                 string
	label                string
	partitionTable       *deviceInfoPartitionTable
	drive                *deviceInfoDrive
	partition            *deviceInfoPartition
}

func monitorDevices(ctx context.Context,
	added func(device string) error,
	changed func(device string) error,
	removed func(device string) error) error {
	cmd := exec.Command("udevil", "--monitor")
	cmd.Stderr = os.Stderr
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		cmd.Process.Kill()
	}()

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
			changed(device)
		} else if action == "added" {
			added(device)
		} else if action == "removed" {
			removed(device)
		}
	}

	cmd.Wait() // TODO: check error code

	return nil
}

func getDeviceInfo(device string) (deviceInfo, error) {
	var result deviceInfo
	cmd := exec.Command("udevil", "--show-info", device)

	b, err := cmd.Output()
	if err != nil {
		return result, err
	}

	output := strings.Split(string(b), "\n")

	r := regexp.MustCompile(`([\w -]*):\s*([\w-\.\/:]*)`)
	getKeyValue := func(value string) (string, string, error) {
		match := r.FindStringSubmatch(value)
		if match == nil {
			return "", "", fmt.Errorf("no match")
		}
		return strings.TrimSpace(match[1]), strings.TrimSpace(match[2]), nil
	}

	for i := 0; i < len(output); i++ {
		line := output[i]
		key, value, err := getKeyValue(line)
		if err != nil {
			continue
		}

		switch key {
		case "native-path":
			result.nativePath = value
			break
		case "device":
			result.device = value
			break
		case "device-file":
			result.deviceFile = &deviceInfoFile{}
			result.deviceFile.file = value
			i++
			for ; i < len(output); i++ {
				deviceLine := output[i]
				if !strings.HasPrefix(deviceLine, "    ") {
					i--
					break
				}
				key, value, err = getKeyValue(deviceLine)
				if err != nil {
					log.Println(err)
					break
				}
				switch key {
				case "presentation":
					result.deviceFile.presentation = value
					break
				case "by-id":
					result.deviceFile.byID = value
					break
				default:
					log.Println("invalid key: " + key)
					break
				}
			}
			break
		case "system internal":
			result.systemInternal = value
			break
		case "removable":
			result.removable = value
			break
		case "has media":
			result.hasMedia = value
			break
		case "is read only":
			result.isReadOnly = value
			break
		case "is mounted":
			result.isMounted = value
			break
		case "mount paths":
			result.mountedPaths = value
			break
		case "presentation hide":
			result.presentationHide = value
			break
		case "presentation nopolicy":
			result.presentationNoPolicy = value
			break
		case "presentation name":
			result.presentationName = value
			break
		case "presentation icon":
			result.presentationIcon = value
			break
		case "automount hint":
			result.autoMountHint = value
			break
		case "size":
			result.size = value
			break
		case "block size":
			result.blockSize = value
			break
		case "usage":
			result.usage = value
			break
		case "type":
			result.deviceType = value
			break
		case "version":
			result.version = value
			break
		case "uuid":
			result.uuid = value
			break
		case "label":
			result.label = value
			break
		case "partition table":
			result.partitionTable = &deviceInfoPartitionTable{}
			result.partitionTable.partitionTable = value
			i++
			for ; i < len(output); i++ {
				partitionLine := output[i]
				if !strings.HasPrefix(partitionLine, "    ") {
					i--
					break
				}
				key, value, err = getKeyValue(partitionLine)
				if err != nil {
					log.Println(err)
					break
				}
				switch key {
				case "scheme":
					result.partitionTable.schema = value
					break
				case "count":
					result.partitionTable.count = value
					break
				default:
					log.Println("invalid key: " + key)
					break
				}
			}
			break
		case "drive":
			result.drive = &deviceInfoDrive{}
			result.drive.drive = value
			i++
			for ; i < len(output); i++ {
				driveLine := output[i]
				if !strings.HasPrefix(driveLine, "    ") {
					i--
					break
				}
				key, value, err = getKeyValue(driveLine)
				if err != nil {
					log.Println(err)
					break
				}
				switch key {
				case "vendor":
					result.drive.vendor = value
					break
				case "model":
					result.drive.model = value
					break
				case "revision":
					result.drive.revision = value
					break
				case "serial":
					result.drive.serial = value
					break
				case "WWN":
					result.drive.WWN = value
					break
				case "detachable":
					result.drive.detachable = value
					break
				case "ejectable":
					result.drive.ejectable = value
					break
				case "media":
					result.drive.media = &deviceInfoDriveMedia{}
					result.drive.media.media = value
					i++
					mediaLine := output[i]
					if !strings.HasPrefix(mediaLine, "      ") {
						i--
						break
					}
					key, value, err = getKeyValue(mediaLine)
					if err != nil {
						log.Println(err)
						break
					}
					switch key {
					case "compat":
						result.drive.media.compat = value
						break
					default:
						log.Println("invalid key: " + key)
						break
					}
					break
				case "interface":
					result.drive.inter = value
					break
				case "if speed":
					result.drive.ifSpeed = value
					break
				default:
					log.Println("invalid key: " + key)
					break
				}
			}
			break
		case "partition":
			result.partition = &deviceInfoPartition{}
			result.partition.partition = value
			i++
			for ; i < len(output); i++ {
				partitionLine := output[i]
				if !strings.HasPrefix(partitionLine, "    ") {
					i--
					break
				}
				key, value, err = getKeyValue(partitionLine)
				if err != nil {
					log.Println(err)
					break
				}
				switch key {
				case "scheme":
					result.partition.scheme = value
					break
				case "number":
					result.partition.number = value
					break
				case "type":
					result.partition.partitionType = value
					break
				case "flags":
					result.partition.flags = value
					break
				case "offset":
					result.partition.offset = value
					break
				case "alignment offset":
					result.partition.alignmentOffset = value
					break
				case "size":
					result.partition.size = value
					break
				case "label":
					result.partition.label = value
					break
				case "uuid":
					result.partition.uuid = value
					break
				default:
					log.Println("invalid key: " + key)
					break
				}
			}
			break
		default:
			log.Println("unknown key " + key)
			break
		}
	}

	return result, nil
}

func mountDevice(device string) error {

}
