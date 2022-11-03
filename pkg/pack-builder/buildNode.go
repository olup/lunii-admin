// TODO : returns
package studiopackbuilder

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/uuid"
	"github.com/olup/lunii-admin/pkg/lunii"
)

type NodeContext struct {
	uuid  string
	index int
}

func getTitleNode(ctx *NodeContext, directoryPath string, tempOutputAssetPath string) ([]lunii.StageNode, []lunii.ListNode, error) {
	fsys := os.DirFS(directoryPath)

	// Create title node
	nodeUuid := uuid.New()
	titleNode := lunii.StageNode{
		Uuid: nodeUuid,
		Name: filepath.Base(directoryPath) + "_title_node",
		Type: "node",
	}

	// audio
	matches, err := doublestar.Glob(fsys, "title.{mp3,ogg,wav}")
	if err != nil || len(matches) == 0 {
		return nil, nil, errors.New("No title audio file found")
	}
	audioPath := filepath.Join(directoryPath, matches[0])

	audioFileName := uuid.NewString() + filepath.Ext(audioPath)
	err = copyFile(audioPath, filepath.Join(tempOutputAssetPath, audioFileName))
	if err != nil {
		return nil, nil, err
	}

	titleNode.Audio = audioFileName

	// cover
	matches, err = doublestar.Glob(fsys, "cover.{jpg,jpeg,png}")
	if err != nil || len(matches) == 0 {
		return nil, nil, errors.New("No cover image file found")
	}
	imagePath := filepath.Join(directoryPath, matches[0])

	imageFileName := uuid.NewString() + filepath.Ext(imagePath)
	copyFile(imagePath, filepath.Join(tempOutputAssetPath, imageFileName))
	if err != nil {
		return nil, nil, err
	}

	titleNode.Image = imageFileName

	// control settings for title node

	titleNode.ControlSettings = &lunii.ControlSettings{
		Wheel:    true,
		Ok:       true,
		Home:     true,
		Pause:    false,
		Autoplay: false,
	}

	// Is there a story node or more title nodes ?
	matches, err = doublestar.Glob(fsys, "story.{mp3,ogg,wav}")
	if err == nil && len(matches) != 0 {
		// We have a story node
		storyAudioPath := filepath.Join(directoryPath, matches[0])

		// copy audio
		audioFileName := uuid.NewString() + ".mp3"
		err = copyFile(storyAudioPath, filepath.Join(tempOutputAssetPath, audioFileName))
		if err != nil {
			return nil, nil, err
		}

		// create node
		storyNode := lunii.StageNode{
			Uuid:  uuid.New(),
			Name:  filepath.Base(directoryPath) + "_story_node",
			Audio: audioFileName,
			Type:  "node",
		}

		// if there is a context, attach home and ok transition to it
		if ctx != nil {
			storyNode.HomeTransition = &lunii.Transition{ActionNode: ctx.uuid, OptionIndex: ctx.index}
			storyNode.OkTransition = &lunii.Transition{ActionNode: ctx.uuid, OptionIndex: ctx.index}
		}

		// create a list node
		listNode := lunii.ListNode{
			Id:      uuid.NewString(),
			Name:    filepath.Base(directoryPath) + "_story_list_node",
			Options: []uuid.UUID{storyNode.Uuid},
		}

		// add list node to title node
		titleNode.OkTransition = &lunii.Transition{
			ActionNode:  listNode.Id,
			OptionIndex: 0,
		}

		titleNode.HomeTransition = nil

		// set story node control settings
		storyNode.ControlSettings = &lunii.ControlSettings{
			Wheel:    false,
			Ok:       false,
			Home:     true,
			Pause:    true,
			Autoplay: true,
		}

		// return nodes and lists
		return []lunii.StageNode{titleNode, storyNode}, []lunii.ListNode{listNode}, nil
	} else {
		// There is no story node - it is a title node
		stageNodes, listNodes, err := listNodesFromDirectory(directoryPath, tempOutputAssetPath)
		if err != nil {
			return nil, nil, err
		}

		titleNode.OkTransition = &lunii.Transition{
			OptionIndex: 0,
			ActionNode:  listNodes[0].Id,
		}

		stageNodes = append([]lunii.StageNode{titleNode}, stageNodes...)
		return stageNodes, listNodes, nil
	}
}

func listNodesFromDirectory(directoryPath string, tempOutputPath string) ([]lunii.StageNode, []lunii.ListNode, error) {
	var stageNodes []lunii.StageNode
	var listNodes []lunii.ListNode
	thisListNode := lunii.ListNode{
		Id:   uuid.NewString(),
		Name: filepath.Base(directoryPath) + "_title_list_node",
	}

	// read each files in directory
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, nil, err
	}

	// for each directory
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		ctx := NodeContext{uuid: thisListNode.Id, index: len(thisListNode.Options)}
		thisStageNodes, thisListNodes, err := getTitleNode(&ctx, filepath.Join(directoryPath, file.Name()), tempOutputPath)
		if err != nil {
			return nil, nil, err
		}
		stageNodes = append(stageNodes, thisStageNodes...)
		listNodes = append(listNodes, thisListNodes...)

		// link top stage node in this list
		thisListNode.Options = append(thisListNode.Options, thisStageNodes[0].Uuid)
	}

	finalListNodes := append([]lunii.ListNode{thisListNode}, listNodes...)

	return stageNodes, finalListNodes, nil
}
