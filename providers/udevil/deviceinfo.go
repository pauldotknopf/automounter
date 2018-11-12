package udevil

import (
	"fmt"
	"log"
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
	media      deviceInfoDriveMedia
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
	deviceFile           deviceInfoFile
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
	partitionTable       deviceInfoPartitionTable
	drive                deviceInfoDrive
	partition            deviceInfoPartition
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
			var deviceFile deviceInfoFile
			deviceFile.file = value
			i++
			for ; i < len(output); i++ {
				deviceLine := output[i]
				if !strings.HasPrefix(deviceLine, "    ") {
					i--
					result.deviceFile = deviceFile
					break
				}
				key, value, err = getKeyValue(deviceLine)
				if err != nil {
					log.Println(err)
					break
				}
				switch key {
				case "presentation":
					deviceFile.presentation = value
					break
				case "by-id":
					deviceFile.byID = value
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
			var partitionTable deviceInfoPartitionTable
			partitionTable.partitionTable = value
			i++
			for ; i < len(output); i++ {
				deviceLine := output[i]
				if !strings.HasPrefix(deviceLine, "    ") {
					i--
					result.partitionTable = partitionTable
					break
				}
				key, value, err = getKeyValue(deviceLine)
				if err != nil {
					log.Println(err)
					break
				}
				switch key {
				case "scheme":
					partitionTable.schema = value
					break
				case "count":
					partitionTable.count = value
					break
				default:
					log.Println("invalid key: " + key)
					break
				}
			}
			break
		case "drive":
			var drive deviceInfoDrive
			drive.drive = value
			i++
			for ; i < len(output); i++ {
				driveLine := output[i]
				if !strings.HasPrefix(driveLine, "    ") {
					i--
					result.drive = drive
					break
				}
				key, value, err = getKeyValue(driveLine)
				if err != nil {
					log.Println(err)
					break
				}
				switch key {
				case "vendor":
					drive.vendor = value
					break
				case "model":
					drive.model = value
					break
				case "revision":
					drive.revision = value
					break
				case "serial":
					drive.serial = value
					break
				case "WWN":
					drive.WWN = value
					break
				case "detachable":
					drive.detachable = value
					break
				case "ejectable":
					drive.ejectable = value
					break
				case "media":
					var media deviceInfoDriveMedia
					media.media = value
					i++
					mediaLine := output[i]
					if !strings.HasPrefix(mediaLine, "      ") {
						i--
						drive.media = media
						break
					}
					key, value, err = getKeyValue(mediaLine)
					if err != nil {
						log.Println(err)
						break
					}
					switch key {
					case "compat":
						media.compat = value
						break
					default:
						log.Println("invalid key: " + key)
						break
					}
					break
				case "interface":
					drive.inter = value
					break
				case "if speed":
					drive.ifSpeed = value
					break
				default:
					log.Println("invalid key: " + key)
					break
				}
			}
			break
		case "partition":
			var partition deviceInfoPartition
			for ; i < len(output); i++ {
				partitionLine := output[i]
				if !strings.HasPrefix(partitionLine, "    ") {
					result.partition = partition
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
					partition.scheme = value
					break
				case "number":
					partition.number = value
					break
				case "type":
					partition.partitionType = value
					break
				case "flags":
					partition.flags = value
					break
				case "offset":
					partition.offset = value
					break
				case "alignment offset":
					partition.alignmentOffset = value
					break
				case "size":
					partition.size = value
					break
				case "label":
					partition.label = value
					break
				case "uuid":
					partition.uuid = value
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
