package lunii

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func GetRefFromUUid(uuid uuid.UUID) string {
	uuidString := uuid.String()
	return strings.ToUpper(strings.ReplaceAll(uuidString[len(uuidString)-8:], "_", ""))
}

func (device *Device) ReadGlobalIndexFile() ([]uuid.UUID, error) {

	// read .pi file and get
	data, err := os.Open(filepath.Join(device.MountPoint, ".pi"))
	if err != nil {
		return nil, errors.New("Could not reach the pack index file")
	}
	defer data.Close()

	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(data)
	slice := make([]byte, 16)
	var ids []uuid.UUID

	for {
		_, err = reader.Read(slice)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("There was an error reading the pack index file")
		}

		uuid, err := uuid.FromBytes(slice)
		if err != nil {
			return nil, errors.New("There was an error getting UUID from the pack index file")
		}
		ids = append(ids, uuid)
	}
	return ids, nil

}

func (device *Device) WriteGlobalIndexFile(stories []uuid.UUID) error {
	var buf []byte
	for _, storyUuid := range stories {
		buf = append(buf, storyUuid[:]...)
	}
	err := os.WriteFile(filepath.Join(device.MountPoint, ".pi"), buf, 0777)
	return err
}

func (device *Device) AddPackToIndex(uuid uuid.UUID) error {
	uuids, err := device.ReadGlobalIndexFile()
	if err != nil {
		return err
	}

	// if the story is already in the index, exit
	for _, storyUuid := range uuids {
		if storyUuid == uuid {
			return nil
		}
	}

	uuids = append(uuids, uuid)
	err = device.WriteGlobalIndexFile(uuids)
	return err
}

func (device *Device) RemovePackFromIndex(thisUuid uuid.UUID) error {
	var uuids []uuid.UUID

	deviceStories, err := device.ReadGlobalIndexFile()
	if err != nil {
		return err
	}

	// filter out the specified UUID
	for _, storyUuid := range deviceStories {
		if storyUuid != thisUuid {
			uuids = append(uuids, storyUuid)
		}
	}

	err = device.WriteGlobalIndexFile(uuids)

	return err
}

func (device *Device) ChangePackOrder(thisUuid uuid.UUID, index int) error {
	var newStoriesUuids []uuid.UUID

	deviceStoriesUuids, err := device.ReadGlobalIndexFile()
	if err != nil {
		return err
	}

	// filter out the specified UUID
	for _, storyUuid := range deviceStoriesUuids {
		if storyUuid != thisUuid {
			newStoriesUuids = append(newStoriesUuids, storyUuid)
		}
	}

	// re-add the specified uuid to the proper index
	newStoriesUuids = insert(newStoriesUuids, thisUuid, index)

	err = device.WriteGlobalIndexFile(newStoriesUuids)

	return err
}
