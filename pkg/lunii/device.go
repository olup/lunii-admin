package lunii

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ricochet2200/go-disk-usage/du"
)

type DiskUsage struct {
	Free uint64 `json:"free"`
	Used uint64 `json:"used"`
	Size uint64 `json:"total"`
}

type Device struct {
	MountPoint           string     `json:"mountPoint"`
	Uuid                 []byte     `json:"uuid"`
	UuidHex              string     `json:"uuidHex"`
	SpecificKey          []byte     `json:"specificKey"`
	SerialNumber         string     `json:"serialNumber"`
	FirmwareVersionMajor int16      `json:"firmwareVersionMajor"`
	FirmwareVersionMinor int16      `json:"firmwareVersionMinor"`
	SdCardSize           int        `json:"sdCardSize"`
	SdCardUsed           int        `json:"sdCardUsed"`
	DiskUsage            *DiskUsage `json:"diskUsage"`
}

func skip(reader *bufio.Reader, count int64) {
	io.CopyN(ioutil.Discard, reader, count)
}

func DetectDeviceMountPoint() (string, error) {

	var mountPoints []string
	if runtime.GOOS == "windows" {
		mountPoints = []string{"A:\\", "B:\\", "C:\\", "D:\\", "E:\\", "F:\\", "G:\\", "H:\\", "I:\\", "J:\\", "K:\\", "L:\\", "M:\\"}
	} else {
		out, err := exec.Command("bash", "-c", "df | sed 1d | awk '{print $NF}'").Output()

		if err != nil {
			return "", err
		}
		mountPoints = strings.Split(string(out), "\n")
	}

	for _, mountPoint := range mountPoints {
		_, err := os.Open(filepath.Join(mountPoint, ".md"))
		if err == nil {
			return mountPoint, nil
		}
	}

	return "", errors.New("No device found")
}

func GetDevice() (*Device, error) {
	var device Device

	mountPoint, err := DetectDeviceMountPoint()
	if err != nil {
		return nil, err
	}
	device.MountPoint = mountPoint

	data, err := os.Open(filepath.Join(mountPoint, ".md"))
	if err != nil {
		return nil, err
	}
	defer data.Close()

	reader := bufio.NewReader(data)

	var short int16

	binary.Read(reader, binary.LittleEndian, &short)

	skip(reader, 4)

	binary.Read(reader, binary.LittleEndian, &short)

	device.FirmwareVersionMajor = short

	binary.Read(reader, binary.LittleEndian, &short)

	device.FirmwareVersionMinor = short

	var long int64
	binary.Read(reader, binary.BigEndian, &long)
	device.SerialNumber = fmt.Sprintf("%014d", long)

	skip(reader, 238)

	slice := make([]byte, 256)
	reader.Read(slice)

	device.Uuid = slice
	device.UuidHex = hex.EncodeToString(slice)
	device.SpecificKey = computeSpecificKeyFromUUID(slice)

	diskUsage := du.NewDiskUsage(device.MountPoint)
	device.DiskUsage = &DiskUsage{
		Free: diskUsage.Free(),
		Used: diskUsage.Used(),
		Size: diskUsage.Size(),
	}

	return &device, nil
}
