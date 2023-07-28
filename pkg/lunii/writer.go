package lunii

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	cp "github.com/otiai10/copy"
	yaml "gopkg.in/yaml.v3"
)

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

func (device *Device) AddStudioPack(studioPack *StudioPack) error {
	start := time.Now()
	CurrentJob = &Job{
		IsComplete: false,
		InitDone:   true,
	}

	// 1. Get path on devide
	rootPath := device.MountPoint
	contentPath := filepath.Join(rootPath, ".content", studioPack.Ref)

	fmt.Println("Generating Binaries...")

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
	fmt.Println("Generating temp data: " + tempPath)

	CurrentJob.BinGenerationDone = true

	fmt.Println("Generating Binaries : Operation took : ", time.Since(start))
	start = time.Now()

	fmt.Println("Preparing asset in " + tempPath)

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

	// unzip all files first
	assetPath := filepath.Join(tempPath, "assets")
	err = os.MkdirAll(assetPath, 0700)
	if err != nil {
		return err
	}

	err = Unzip(studioPack.OriginalPath, assetPath)
	if err != nil {
		return err
	}

	CurrentJob.UnpackDone = true

	reader := os.DirFS(assetPath)

	fmt.Println("Unzipping asset : Operation took : ", time.Since(start))
	start = time.Now()

	// copy images in rf
	fmt.Println("Converting images ...")
	fmt.Println(len(*imageIndex))

	CurrentJob.TotalImages = len(*imageIndex)

	deviceImageDirectory := filepath.Join(tempPath, "rf", "000")
	os.MkdirAll(deviceImageDirectory, 0700)

	wg := sync.WaitGroup{}
	for i, image := range *imageIndex {
		go func() {
			wg.Add(1)
			convertAndWriteImage(reader, deviceImageDirectory, image, i)
			CurrentJob.ImagesDone++
			wg.Done()
		}()
	}
	wg.Wait()

	CurrentJob.ImagesConversionDone = true

	fmt.Println("Converting images : Operation took : ", time.Since(start))
	start = time.Now()

	// copy audios in sf
	fmt.Println("Converting audios ...")

	deviceAudioDirectory := filepath.Join(tempPath, "sf", "000")
	os.MkdirAll(deviceAudioDirectory, 0700)

	CurrentJob.TotalAudios = len(*audioIndex)

	wg = sync.WaitGroup{}

	for i, audio := range *audioIndex {
		go func() {
			wg.Add(1)
			convertAndWriteAudio(reader, deviceAudioDirectory, audio, i)
			CurrentJob.AudiosDone++
			wg.Done()
		}()
	}

	wg.Wait()

	CurrentJob.AudiosConversionDone = true

	fmt.Println("Converting audios : Operation took : ", time.Since(start))
	start = time.Now()

	// adding metadata
	fmt.Println("Writing metadata...")

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

	CurrentJob.MetadataDone = true

	fmt.Println("Writing metadata : Operation took : ", time.Since(start))
	start = time.Now()

	// copy temp to lunii
	fmt.Println("Copying directory to the device...")

	cp.Copy(tempPath, contentPath)

	CurrentJob.CopyingDone = true

	fmt.Println("Copying directory to the device : Operation took : ", time.Since(start))
	start = time.Now()

	fmt.Println("Adding pack to root index...")

	// // update .pi root file with uuid
	err = device.AddPackToIndex(studioPack.Uuid)
	if err != nil {
		return err
	}

	CurrentJob.IndexDone = true

	fmt.Println("Adding pack to root index : Operation took : ", time.Since(start))
	start = time.Now()

	fmt.Println("Cleaning...")
	_ = os.RemoveAll(tempPath)

	fmt.Println("Cleaning : Operation took : ", time.Since(start))

	return nil
}

func convertAndWriteAudio(reader fs.FS, deviceAudioDirectory string, audio Asset, index int) error {
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

func convertAndWriteImage(reader fs.FS, deviceImageDirectory string, image Asset, index int) error {
	file, err := reader.Open("assets/" + image.SourceName)
	if err != nil {
		return err
	}

	bmpFile, err := ImageReadOrDecode(image.SourceName, file)
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
