package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/google/uuid"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"github.com/olup/lunii-cli/pkg/lunii"
	studiopackbuilder "github.com/olup/lunii-cli/pkg/pack-builder"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDeviceInfos() *lunii.Device {
	device, err := lunii.GetDevice()
	if err != nil {
		return nil
	}
	return device
}

func (a *App) ListPacks() []lunii.Metadata {
	device, err := lunii.GetDevice()
	if err != nil {
		return nil
	}
	metadatas, err := device.GetPacks()
	if err != nil {
		return nil
	}
	return metadatas
}

func (a *App) RemovePack(uuid uuid.UUID) (bool, error) {
	device, err := lunii.GetDevice()
	if err != nil {
		return false, err
	}

	err = device.RemovePackFromIndex(uuid)
	if err != nil {
		return false, err
	}

	err = os.RemoveAll(filepath.Join(device.MountPoint, ".content", lunii.GetRefFromUUid(uuid)))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *App) CreatePack(directoryPath string, destinationPath string) (string, error) {
	_, err := studiopackbuilder.CreateStudioPack(directoryPath, destinationPath)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (a *App) OpenDirectory(title string) string {
	path, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	return path
}

func (a *App) SaveFile(title string, defaultDirectory string, defaultFileName string) string {
	fmt.Println("Select save path - options : ", defaultDirectory, defaultFileName)
	path, _ := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            title,
		DefaultDirectory: defaultDirectory,
		DefaultFilename:  defaultFileName,
	})
	return path
}

func (a *App) OpenFile(title string) string {
	path, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	return path
}

func (a *App) InstallPack(packPath string) (string, error) {

	device, err := lunii.GetDevice()
	studioPack, err := lunii.ReadStudioPack(packPath)
	if err != nil {
		return "", err
	}

	err = device.AddStudioPack(studioPack)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (a *App) ChangePackOrder(uuid uuid.UUID, index int) (string, error) {

	fmt.Println("Moving ", uuid, index)

	device, err := lunii.GetDevice()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	err = device.ChangePackOrder(uuid, index)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return "", nil
}

func (a *App) SyncLuniiStoreMetadata(uuids []uuid.UUID) (string, error) {
	device, err := lunii.GetDevice()
	if err != nil {
		return "", err
	}

	db, err := lunii.GetLuniiMetadataDb()
	if err != nil {
		return "", err
	}

	for _, thisUuid := range uuids {
		device.SyncMetadataFromDb(thisUuid, db)
	}

	return "", nil
}

func (a *App) SyncStudioMetadata(uuids []uuid.UUID, dbPath string) (string, error) {
	device, err := lunii.GetDevice()
	if err != nil {
		return "", err
	}

	db, err := lunii.GetStudioMetadataDb(dbPath)
	if err != nil {
		return "", err
	}

	for _, thisUuid := range uuids {
		device.SyncMetadataFromDb(thisUuid, db)
	}

	return "", nil
}

func (a *App) CheckForUpdate() (bool, string, string) {
	latest, found, err := selfupdate.DetectLatest("olup/lunii-admin")
	v := semver.MustParse("0.0.3")
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return false, "", ""
	}

	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return false, "", ""
	}
	return true, latest.Version.String(), latest.ReleaseNotes
}
