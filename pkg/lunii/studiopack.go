package lunii

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
)

type StudioPack struct {
	PackType     string      `json:"_"`
	OriginalPath string      `json:"_"`
	Format       string      `json:"format"` // enum ?
	Title        string      `json:"title"`
	Version      int         `json:"version"` // enum ?
	Description  string      `json:"description"`
	StageNodes   []StageNode `json:"stageNodes"`
	ListNodes    []ListNode  `json:"actionNodes"`
	Uuid         uuid.UUID
	Ref          string
}

type StageNode struct {
	Uuid            uuid.UUID        `json:"uuid"`
	Type            string           `json:"type"` //enum ?
	Name            string           `json:"name"`
	Image           string           `json:"image"`
	Audio           string           `json:"audio"`
	OkTransition    *Transition      `json:"okTransition"`
	HomeTransition  *Transition      `json:"homeTransition"`
	ControlSettings *ControlSettings `json:"controlSettings"`
	SquareOne       bool             `json:"squareOne"`
}

type ListNode struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Options []uuid.UUID `json:"options"`
}

type ControlSettings struct {
	Wheel    bool `json:"wheel"`
	Ok       bool `json:"ok"`
	Home     bool `json:"home"`
	Pause    bool `json:"pause"`
	Autoplay bool `json:"autoplay"`
}

type Transition struct {
	ActionNode  string `json:"actionNode"`
	OptionIndex int    `json:"optionIndex"`
}

func ReadStudioPack(path string) (*StudioPack, error) {

	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, errors.New("package could not be found")
	}
	defer reader.Close()

	file, err := reader.Open("story.json")
	if err != nil {
		log.Fatal("story.json could not be found")
	}
	fileAsBytes, _ := ioutil.ReadAll(file)

	var pack StudioPack

	pack.OriginalPath = path
	json.Unmarshal(fileAsBytes, &pack)

	pack.Uuid = pack.StageNodes[0].Uuid
	pack.Ref = GetRefFromUUid(pack.Uuid)

	return &pack, nil
}
