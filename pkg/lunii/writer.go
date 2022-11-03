package lunii

import (
	"archive/zip"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	cp "github.com/otiai10/copy"
	yaml "gopkg.in/yaml.v3"
)

var wg sync.WaitGroup

func sendUpdate(updateChan *chan string, msg string) {
	// if no channel ?
	if updateChan == nil {
		return
	}

	// Non blocking send
	select {
	case *updateChan <- msg:
		return
	default:
		return
	}
}

func (device *Device) AddStudioPack(studioPack *StudioPack, updateChan *chan string) error {
	start := time.Now()
	sendUpdate(updateChan, "STARTING")
	defer sendUpdate(updateChan, "DONE")

	// 1. Get path on devide
	rootPath := device.MountPoint
	contentPath := filepath.Join(rootPath, ".content", studioPack.Ref)

	fmt.Println("Generating Binaries...")
	sendUpdate(updateChan, "GENERATING_BINS")

	stageNodeIndex := &studioPack.StageNodes
	// generate list node index
	listNodeIndex := GetListNodeIndex(&studioPack.ListNodes)

	// create image index - with exporter & lookup
	imageIndex := GetImageAssetListFromPack(studioPack)
	// create sound index - with exporter & lookup
	audioIndex := GetAudioAssetListFromPack(studioPack)

	// prepare ni (index of stage nodes)
	niBinary := GenerateNiBinary(studioPack, stageNodeIndex, listNodeIndex, imageIndex, audioIndex)

	// prepare li (index of action nodes)
	liBinary := generateLiBinary(listNodeIndex, stageNodeIndex)
	liBinaryCiphered := cipherFirstBlockCommonKey(liBinary)

	// prepare ri (index of image assets)
	riBinary := GenerateBinaryFromAssetIndex(imageIndex)
	riBinaryCiphered := cipherFirstBlockCommonKey(riBinary)

	// prepare si (index of sound assets)
	siBinary := GenerateBinaryFromAssetIndex(audioIndex)
	siBinaryCiphered := cipherFirstBlockCommonKey(siBinary)

	// create boot file
	btBinary := generateBtBinary(riBinaryCiphered)

	// prepare pack
	tempPath := filepath.Join(os.TempDir(), "packs", studioPack.Ref)
	err := os.MkdirAll(tempPath, 0700)
	if err != nil {
		return err
	}

	fmt.Println("Generating Binaries : Operation took : ", time.Since(start))
	start = time.Now()

	fmt.Println("Preparing asset in " + tempPath)
	sendUpdate(updateChan, "PREPARING_ASSETS")

	err = os.WriteFile(filepath.Join(tempPath, "ni"), niBinary, 0777)
	err = os.WriteFile(filepath.Join(tempPath, "li"), liBinaryCiphered, 0777)
	err = os.WriteFile(filepath.Join(tempPath, "ri"), riBinaryCiphered, 0777)
	err = os.WriteFile(filepath.Join(tempPath, "si"), siBinaryCiphered, 0777)
	err = os.WriteFile(filepath.Join(tempPath, "bt"), btBinary, 0777)
	err = os.MkdirAll(filepath.Join(tempPath, "sf"), 0700)
	err = os.MkdirAll(filepath.Join(tempPath, "rf"), 0700)
	if err != nil {
		return err
	}

	// prepare zip reader
	reader, err := zip.OpenReader(studioPack.OriginalPath)
	if err != nil {
		return errors.New("Zip package could not be opened")
	}
	defer reader.Close()

	fmt.Println("Preparing asset : Operation took : ", time.Since(start))
	start = time.Now()

	// copy images in rf
	fmt.Println("Converting images ...")
	sendUpdate(updateChan, "CONVERTING_IMAGES")

	deviceImageDirectory := filepath.Join(tempPath, "rf", "000")
	os.MkdirAll(deviceImageDirectory, 0700)

	for i, image := range *imageIndex {
		wg.Add(1)
		go convertAndWriteImage(*reader, deviceImageDirectory, image, i)
	}

	wg.Wait()

	fmt.Println("Converting images : Operation took : ", time.Since(start))
	start = time.Now()

	// copy audios in sf
	fmt.Println("Converting audios ...")
	sendUpdate(updateChan, "CONVERTING_AUDIOS")

	deviceAudioDirectory := filepath.Join(tempPath, "sf", "000")
	os.MkdirAll(deviceAudioDirectory, 0700)

	for i, audio := range *audioIndex {
		wg.Add(1)
		go convertAndWriteAudio(*reader, deviceAudioDirectory, audio, i)
	}

	wg.Wait()

	fmt.Println("Converting audios : Operation took : ", time.Since(start))
	start = time.Now()

	// adding metadata
	fmt.Println("Writing metadata...")
	sendUpdate(updateChan, "WRITING_METADATA")

	md := Metadata{
		Uuid:        studioPack.Uuid,
		Ref:         GetRefFromUUid(studioPack.Uuid),
		Title:       studioPack.Title,
		Description: studioPack.Description,
		PackType:    "custom",
	}

	yaml, err := yaml.Marshal(&md)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(tempPath, "md"), yaml, 0777)
	if err != nil {
		return err
	}
	fmt.Println("Writing metadata : Operation took : ", time.Since(start))
	start = time.Now()

	// copy temp to lunii
	fmt.Println("Copying directory to the device...")
	sendUpdate(updateChan, "COPYING")

	cp.Copy(tempPath, contentPath)

	fmt.Println("Copying directory to the device : Operation took : ", time.Since(start))
	start = time.Now()

	fmt.Println("Adding pack to root index...")
	sendUpdate(updateChan, "UPDATING_INDEX")

	// // update .pi root file with uuid
	err = device.AddPackToIndex(studioPack.Uuid)
	if err != nil {
		return err
	}

	fmt.Println("Adding pack to root index : Operation took : ", time.Since(start))
	start = time.Now()

	fmt.Println("Cleaning...")
	_ = os.RemoveAll(tempPath)

	fmt.Println("Cleaning : Operation took : ", time.Since(start))

	return nil
}

func convertAndWriteAudio(reader zip.ReadCloser, deviceAudioDirectory string, audio Asset, index int) error {
	defer wg.Done()
	var mp3 []byte
	var err error

	// if this is an empty file, we just write a blank mp3
	if audio.SourceName == "EMPTY_SOUND" {
		mp3, _ = hex.DecodeString(BLANK_MP3_FILE)

	} else {
		start := time.Now()

		//check extension
		extension := audio.SourceName[len(audio.SourceName)-3:]

		// otherwise, let's convert the real file
		file, err := reader.Open("assets/" + audio.SourceName)
		if err != nil {
			return err
		}

		defer file.Close()

		if extension == "mp3" {
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			mp3, err = Mp3ToMp3(fileBytes)
			if err != nil {
				return err
			}
		} else if extension == "ogg" {
			// it's an ogg file

			mp3, err = OggToMp3(file)
			if err != nil {
				return err
			}

		} else {
			return errors.New("Audio file format not supported")
		}
		fmt.Println("Audio converted in ", time.Since(start))

	}

	cypheredFile := cipherFirstBlockCommonKey(mp3)

	err = os.WriteFile(filepath.Join(deviceAudioDirectory, intTo8Chars(index)), cypheredFile, 0777)
	if err != nil {
		return err
	}

	return nil
}

func convertAndWriteImage(reader zip.ReadCloser, deviceImageDirectory string, image Asset, index int) error {
	defer wg.Done()

	file, err := reader.Open("assets/" + image.SourceName)
	if err != nil {
		return err
	}
	bmpFile, err := ImageToBmp4(file)
	if err != nil {
		return err
	}
	cypheredBmp := cipherFirstBlockCommonKey(bmpFile)
	err = os.WriteFile(filepath.Join(deviceImageDirectory, intTo8Chars(index)), cypheredBmp, 0777)
	if err != nil {
		return err
	}
	return nil
}
