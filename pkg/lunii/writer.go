package lunii

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
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
	unpkgPath := filepath.Join(os.TempDir(), "packsUnpkg", studioPack.Ref)
	err = os.MkdirAll(unpkgPath, 0700)
	if err != nil {
		return err
	}
	fmt.Println("Unzipping asset in " + unpkgPath)
	defer os.RemoveAll(unpkgPath)

	err = Unzip(studioPack.OriginalPath, unpkgPath)
	if err != nil {
		return fmt.Errorf("Unzip error: %v", err)
	}

	CurrentJob.UnpackDone = true

	fmt.Println("Unzipping asset : Operation took : ", time.Since(start))

	reader := os.DirFS(unpkgPath)

	start = time.Now()

	// copy images in rf
	fmt.Println("Converting images ...")
	fmt.Println(len(*imageIndex))

	CurrentJob.TotalImages = len(*imageIndex)

	deviceImageDirectory := filepath.Join(tempPath, "rf", "000")
	os.MkdirAll(deviceImageDirectory, 0700)

	var imgwg sync.WaitGroup

	for i, _image := range *imageIndex {
		index := i
		image := _image

		go func() {
			imgwg.Add(1)
			convertAndWriteImage(reader, deviceImageDirectory, image, index)
			CurrentJob.ImagesDone++
			imgwg.Done()
		}()
	}
	imgwg.Wait()

	CurrentJob.ImagesConversionDone = true

	fmt.Println("Converting images : Operation took : ", time.Since(start))
	start = time.Now()

	// copy audios in sf
	fmt.Println("Converting audios ...")

	deviceAudioDirectory := filepath.Join(tempPath, "sf", "000")
	os.MkdirAll(deviceAudioDirectory, 0700)

	CurrentJob.TotalAudios = len(*audioIndex)

	var sdwg sync.WaitGroup

	for i, _audio := range *audioIndex {
		fmt.Println("Converting audio : ", _audio.SourceName)
		index := i
		audio := _audio
		go func() {
			sdwg.Add(1)
			err := convertAndWriteAudio(reader, deviceAudioDirectory, audio, index)
			if err != nil {
				fmt.Println("Error converting audio : ", err)
			}
			CurrentJob.AudiosDone++
			sdwg.Done()
		}()
	}

	sdwg.Wait()

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
	outPutFile, err := os.OpenFile(filepath.Join(deviceAudioDirectory, intTo8Chars(index)), os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer outPutFile.Close()

	// if this is an empty file, we just write a blank mp3
	if audio.SourceName == "EMPTY_SOUND" {
		mp3, _ := hex.DecodeString(BLANK_MP3_FILE)
		_, err = outPutFile.Write(mp3)
		if err != nil {
			return err
		}
	} else {
		start := time.Now()

		//check extension
		extension := audio.SourceName[len(audio.SourceName)-3:]

		file, err := reader.Open("assets/" + audio.SourceName)
		if err != nil {
			return err
		}

		defer file.Close()

		if extension == "mp3" {
			err := Mp3ToMp3(file, outPutFile)
			if err != nil {
				return err
			}

		} else if extension == "ogg" {
			// it's an ogg file

			err = OggToMp3(file, outPutFile)
			if err != nil {
				return err
			}

		} else {
			return errors.New("Audio file format not supported")
		}
		fmt.Println("Audio converted in ", time.Since(start))
	}

	// TODO cipher from a file on disk to save memory
	//cipherFileFirstBlockCommonKey(outPutFile)

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
