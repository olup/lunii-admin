package lunii

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMp3(t *testing.T) {
	reader, err := zip.OpenReader("/Users/louptopalian/lunii/yuna.zip")
	if err != nil {
		log.Fatal(err)
	}
	file, err := reader.Open("assets/" + "15e2004b-31b2-4af7-9550-976b2b802760.mp3")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	result, err := Mp3ToMp3(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(result))
	os.WriteFile("../../reference/test.mp3", result, 0777)
}
