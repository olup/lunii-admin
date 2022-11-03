package lunii

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Uuid        uuid.UUID `yaml:"uuid" json:"uuid"`
	Ref         string    `yaml:"ref" json:"ref"`
	Title       string    `yaml:"title" json:"title"`
	Description string    `yaml:"description" json:"description"`
	// undefined or lunii or custom
	PackType string `yaml:"packType" json:"packType"`
}

func (device *Device) GetPacks() ([]Metadata, error) {
	var packs []Metadata

	uuids, err := device.ReadGlobalIndexFile()
	if err != nil {
		return nil, err
	}

	for _, storyUuid := range uuids { // Read md file
		metadata, _ := device.GetPacksMetadata(storyUuid)
		if metadata == nil {
			metadata = &Metadata{
				Uuid:        storyUuid,
				Ref:         GetRefFromUUid(storyUuid),
				Title:       "",
				Description: "",
				PackType:    "undefined",
			}
		}
		packs = append(packs, *metadata)
	}
	return packs, nil
}

func (device *Device) GetPacksMetadata(uuid uuid.UUID) (*Metadata, error) {
	mdFilePath := filepath.Join(device.MountPoint, ".content", GetRefFromUUid(uuid), "md")
	metadataFile, err := os.ReadFile(mdFilePath)
	if err != nil {
		return nil, err
	}
	metadata := Metadata{}
	err = yaml.Unmarshal(metadataFile, &metadata)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &metadata, nil
}

func (device *Device) SyncMetadataFromDb(uuid uuid.UUID, db *Db) error {
	story := db.GetStoryById(uuid)

	if story == nil {
		return errors.New("Could not found this uuid in DB")
	}

	md := Metadata{
		Uuid:        uuid,
		Ref:         GetRefFromUUid(uuid),
		Title:       story.Title,
		Description: story.Description,
		PackType:    story.PackType,
	}

	mdYaml, err := yaml.Marshal(&md)
	if err != nil {
		return err
	}

	mdFilePath := filepath.Join(device.MountPoint, ".content", GetRefFromUUid(uuid), "md")

	err = os.WriteFile(mdFilePath, mdYaml, 0777)
	if err != nil {
		return err
	}

	return nil
}
