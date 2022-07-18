package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

func (a *App) CreatePack() (string, error) {
	directoryPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose directory",
	})
	if err != nil || directoryPath == "" {
		return "", err
	}

	outputPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Where to save the pack ?",
		DefaultFilename: "pack.zip",
	})
	if err != nil || outputPath == "" {
		return "", err
	}

	_, err = studiopackbuilder.CreateStudioPack(directoryPath, outputPath)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (a *App) InstallPack() (string, error) {
	packPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose Zip Pack",
	})
	if err != nil || packPath == "" {
		return "", err
	}

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

	device, err := lunii.GetDevice()
	if err != nil {
		return "", err
	}

	err = device.ChangePackOrder(uuid, index)
	if err != nil {
		return "", err
	}

	return "", nil
}
