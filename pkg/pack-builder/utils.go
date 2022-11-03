package studiopackbuilder

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

func copyFile(from string, to string) error {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(to, input, 0777)
	if err != nil {
		return err
	}
	return nil
}

func zipDir(directoryPath string, archivePath string) error {
	fmt.Println(directoryPath)
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		directoryPath + string(filepath.Separator): "",
	})
	if err != nil {
		return err
	}

	// create the output file we'll write to
	out, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// we can use the CompressedArchive type to gzip a tarball
	// (compression is not required; you could use Tar directly)
	format := archiver.CompressedArchive{
		Archival: archiver.Zip{},
	}

	// create the archive
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return err
	}
	return nil
}
