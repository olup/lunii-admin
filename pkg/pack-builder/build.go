package studiopackbuilder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/olup/lunii-admin/pkg/lunii"
	"gopkg.in/yaml.v3"
)

func CreateStudioPack(directoryPath string, outputPath string) (*lunii.StudioPack, error) {
	var metadata lunii.Metadata
	metadataPath := filepath.Join(directoryPath, "md.yaml")
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, err
	}

	// get ref from uuid
	metadata.Ref = lunii.GetRefFromUUid(metadata.Uuid)

	tempOutputPath := filepath.Join(os.TempDir(), "build", metadata.Ref)
	tempOutputAssetPath := filepath.Join(tempOutputPath, "assets")

	err = os.MkdirAll(tempOutputPath, 0700)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(tempOutputAssetPath, 0700)
	if err != nil {
		return nil, err
	}

	// start node grabbing
	stageNodes, listNodes, err := getTitleNode(nil, directoryPath, tempOutputAssetPath)
	if err != nil {
		return nil, err
	}

	studioPack := &lunii.StudioPack{
		Uuid:        metadata.Uuid,
		Title:       metadata.Title,
		Ref:         metadata.Ref,
		Description: metadata.Description,

		StageNodes: stageNodes,
		ListNodes:  listNodes,
		PackType:   "",
		Version:    2,
	}

	// set first node's uuid to the pack one
	studioPack.StageNodes[0].Uuid = studioPack.Uuid
	studioPack.StageNodes[0].Type = "cover"

	// convert to JSON
	jsonBytes, err := json.Marshal(&studioPack)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filepath.Join(tempOutputPath, "story.json"), jsonBytes, 0777)
	if err != nil {
		return nil, err
	}

	// create archive
	err = zipDir(tempOutputPath, outputPath)
	if err != nil {
		return nil, err
	}

	// clean
	err = os.RemoveAll(tempOutputPath)
	if err != nil {
		fmt.Println("Warning: Could not remove the temporary directory")
	}

	return studioPack, nil
}
