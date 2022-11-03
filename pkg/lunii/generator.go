package lunii

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
)

type Asset struct {
	index      int
	SourceName string
	NodeUuid   uuid.UUID
}

type ListNodeIndexed struct {
	Position int
	Length   int
	Index    int
	Node     ListNode
}

func GetImageAssetListFromPack(archive *StudioPack) *[]Asset {
	var imageAssetsList []Asset
	index := 0

	for i := 0; i < len(archive.StageNodes); i++ {
		imageAssetFileName := archive.StageNodes[i].Image
		if imageAssetFileName == "" {
			continue
		}

		imageAssetsList = append(imageAssetsList, Asset{
			index:      index,
			SourceName: imageAssetFileName,
			NodeUuid:   archive.StageNodes[i].Uuid,
		})
		index++
	}
	return &imageAssetsList
}

func GetAudioAssetListFromPack(archive *StudioPack) *[]Asset {
	var soundAssetsList []Asset
	index := 0

	for i := 0; i < len(archive.StageNodes); i++ {
		soundAssetFileName := archive.StageNodes[i].Audio

		if soundAssetFileName == "" {
			// in case that the audio is null, we reference a special name
			// and backfill it to our package representation

			archive.StageNodes[i].Audio = "EMPTY_SOUND"
			soundAssetFileName = "EMPTY_SOUND"
		}

		soundAssetsList = append(soundAssetsList, Asset{
			index:      index,
			SourceName: soundAssetFileName,
			NodeUuid:   archive.StageNodes[i].Uuid,
		})
		index++
	}
	return &soundAssetsList
}

func getAssetByName(name string, assets *[]Asset) *Asset {
	for _, asset := range *assets {
		if name == asset.SourceName {
			return &asset
		}
	}
	return nil
}

func getStageNodeIndexByUuid(uuid uuid.UUID, nodes *[]StageNode) int {
	for i, node := range *nodes {
		if uuid == node.Uuid {
			return i
		}
	}
	return -1
}

func intTo8Chars(i int) string {
	return fmt.Sprintf("%08d", i)
}

func GenerateBinaryFromAssetIndex(assets *[]Asset) []byte {
	var bin []byte
	for i := range *assets {
		path := "000\\" + intTo8Chars(i)
		bin = append(bin, path...)
	}
	return bin
}

func GetListNodeIndex(listNodes *[]ListNode) *[]ListNodeIndexed {
	var listNodeIndex []ListNodeIndexed
	pos := 0
	for i, node := range *listNodes {
		listNodeIndex = append(listNodeIndex, ListNodeIndexed{
			Index:    i,
			Position: pos,
			Length:   len(node.Options),
			Node:     node,
		})
		pos += len(node.Options)
	}
	return &listNodeIndex
}

func GetLisNodeIndexedById(id string, nodes *[]ListNodeIndexed) *ListNodeIndexed {
	for _, node := range *nodes {
		if id == node.Node.Id {
			return &node
		}
	}
	return nil
}

func generateLiBinary(nodes *[]ListNodeIndexed, stageNodes *[]StageNode) []byte {
	buf := new(bytes.Buffer)
	for _, node := range *nodes {
		// for each node
		for _, option := range node.Node.Options {
			// for each option, write stage node index in the buffer
			// todo What happen if we can't find the stage node ?
			optionIndex := getStageNodeIndexByUuid(option, stageNodes)
			binary.Write(buf, binary.LittleEndian, uint32(optionIndex))
		}
	}
	return buf.Bytes()
}

func GenerateNiBinary(pack *StudioPack, stageNodes *[]StageNode, listNodeIndex *[]ListNodeIndexed, imageIndex *[]Asset, soundIndex *[]Asset) []byte {
	buf := new(bytes.Buffer)

	// Nodes index file format version (1)
	binary.Write(buf, binary.LittleEndian, uint16(1))
	// Story pack version (1)
	binary.Write(buf, binary.LittleEndian, int16(pack.Version))

	// Start of actual nodes list in this file (0x200 / 512)
	binary.Write(buf, binary.LittleEndian, int32(512))

	// Size of a stage node in this file (0x2C / 44)
	binary.Write(buf, binary.LittleEndian, int32(44))

	// Number of stage nodes in this file
	binary.Write(buf, binary.LittleEndian, int32(len(*stageNodes)))
	// Number of images (in RI file and rf/ folder)
	binary.Write(buf, binary.LittleEndian, int32(len(*imageIndex)))

	// Number of sounds (in SI file and sf/ folder)
	binary.Write(buf, binary.LittleEndian, int32(len(*soundIndex)))

	// Is factory pack : byte to one to avoid pack inspection by official Luniistore application
	binary.Write(buf, binary.LittleEndian, int8(1))

	// Jump to address 0x200 (512) for actual list of nodes (already written 25 bytes)
	buf.Write(make([]byte, 512-25))

	// write each stage node
	for _, node := range *stageNodes {

		// image might be empty
		if node.Image == "" {
			binary.Write(buf, binary.LittleEndian, int32(-1))
		} else {
			binary.Write(buf, binary.LittleEndian, int32(getAssetByName(node.Image, imageIndex).index))
		}

		binary.Write(buf, binary.LittleEndian, int32(getAssetByName(node.Audio, soundIndex).index))

		// okTransition might be empty

		okTransition := node.OkTransition

		if okTransition == nil {
			binary.Write(buf, binary.LittleEndian, int32(-1))
			binary.Write(buf, binary.LittleEndian, int32(-1))
			binary.Write(buf, binary.LittleEndian, int32(-1))
		} else {
			actionNode := GetLisNodeIndexedById(okTransition.ActionNode, listNodeIndex)
			binary.Write(buf, binary.LittleEndian, int32(actionNode.Position))
			binary.Write(buf, binary.LittleEndian, int32(actionNode.Length))
			binary.Write(buf, binary.LittleEndian, int32(okTransition.OptionIndex))
		}

		// hometransition might be empty
		homeTransition := node.HomeTransition
		if homeTransition == nil {
			binary.Write(buf, binary.LittleEndian, int32(-1))
			binary.Write(buf, binary.LittleEndian, int32(-1))
			binary.Write(buf, binary.LittleEndian, int32(-1))
		} else {
			actionNode := GetLisNodeIndexedById(homeTransition.ActionNode, listNodeIndex)
			binary.Write(buf, binary.LittleEndian, int32(actionNode.Position))
			binary.Write(buf, binary.LittleEndian, int32(actionNode.Length))
			binary.Write(buf, binary.LittleEndian, int32(homeTransition.OptionIndex))
		}

		binary.Write(buf, binary.LittleEndian, boolToShort(node.ControlSettings.Wheel))
		binary.Write(buf, binary.LittleEndian, boolToShort(node.ControlSettings.Ok))
		binary.Write(buf, binary.LittleEndian, boolToShort(node.ControlSettings.Home))
		binary.Write(buf, binary.LittleEndian, boolToShort(node.ControlSettings.Pause))
		binary.Write(buf, binary.LittleEndian, boolToShort(node.ControlSettings.Autoplay))
		binary.Write(buf, binary.LittleEndian, int16(0))

	}

	return buf.Bytes()

}

func generateBtBinary(niBinary []byte) []byte {
	chunckSize := min(len(niBinary), 64)
	return cipherBlockSpecificKey(niBinary[0:chunckSize])
}
