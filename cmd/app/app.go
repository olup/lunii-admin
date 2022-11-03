package main

import (
	"context"
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"github.com/olup/lunii-admin/pkg/lunii"
	studiopackbuilder "github.com/olup/lunii-admin/pkg/pack-builder"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	log "github.com/sirupsen/logrus"
)

var version string
var machineId, _ = machineid.ID()
var nrApp *newrelic.Application

// App struct
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

var NR_LICENCE = ""

func (a *App) startup(ctx context.Context) {
	var err error
	log.SetFormatter(&log.JSONFormatter{})
	defer HandlePanic()

	nrApp, err = newrelic.NewApplication(
		newrelic.ConfigAppName("lunii-admin"),
		newrelic.ConfigLicense(NR_LICENCE),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigAppLogEnabled(true),
	)

	if err != nil {
		log.Error(err)
		log.Error("NR init: %s", err)
	}
	logHook := NewHook(nrApp, DefaultLevels)
	log.AddHook(logHook)

	log.Info("NR app loaded")

	a.ctx = ctx

	log.Info("App started")
}

func (a *App) GetDeviceInfos() *lunii.Device {
	defer HandlePanic()
	log.Info("Get device info")

	device, err := lunii.GetDevice()
	if err != nil {
		log.Error(err)
		return nil
	}
	return device
}

func (a *App) ListPacks() []lunii.Metadata {
	defer HandlePanic()
	log.Info("List pack")

	device, err := lunii.GetDevice()
	if err != nil {
		log.Error(err)
		return nil
	}
	metadatas, err := device.GetPacks()
	if err != nil {
		log.Error(err)
		return nil
	}
	return metadatas
}

func (a *App) RemovePack(uuid uuid.UUID) (bool, error) {
	defer HandlePanic()
	log.Info("Remove pack")

	device, err := lunii.GetDevice()
	if err != nil {
		log.Error(err)
		return false, err
	}

	err = device.RemovePackFromIndex(uuid)
	if err != nil {
		log.Error(err)
		return false, err
	}

	err = os.RemoveAll(filepath.Join(device.MountPoint, ".content", lunii.GetRefFromUUid(uuid)))
	if err != nil {
		log.Error(err)
		return false, err
	}

	return true, nil
}

func (a *App) CreatePack(directoryPath string, destinationPath string) (string, error) {
	defer HandlePanic()
	log.Info("Create pack")

	_, err := studiopackbuilder.CreateStudioPack(directoryPath, destinationPath)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return "", nil
}

func (a *App) OpenDirectory(title string) string {
	defer HandlePanic()
	log.Info("Open directory")

	path, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	return path
}

func (a *App) SaveFile(title string, defaultDirectory string, defaultFileName string) string {
	defer HandlePanic()
	log.Info("Save file")

	log.Info("Select save path - options : ", defaultDirectory, defaultFileName)
	path, _ := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            title,
		DefaultDirectory: defaultDirectory,
		DefaultFilename:  defaultFileName,
	})
	return path
}

func (a *App) OpenFile(title string) string {
	defer HandlePanic()
	log.Info("Open file")

	path, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	return path
}

// global var, so that we can't start two installing jobs at the same time
var isInstalling = false

func (a *App) InstallPack(packPath string) (string, error) {
	if isInstalling == true {
		log.Error("Installation already in process")
	}

	defer HandlePanic()
	log.Info("Install pack")

	device, err := lunii.GetDevice()
	studioPack, err := lunii.ReadStudioPack(packPath)

	if err != nil {
		log.Error(err)
		return "", err
	}

	channel := make(chan string, 100)

	isInstalling = true
	go func() {
		err := device.AddStudioPack(studioPack, &channel)
		isInstalling = false

		if err != nil {
			log.Error(err)
			channel <- "ERROR"
			return
		}
	}()

	for message := range channel {
		log.Info("Update: " + message)
		runtime.EventsEmit(a.ctx, "INSTALL_EVENT", message)
		if message == "DONE" {
			break
		}
		if message == "ERROR" {
			return "", errors.New("An error happened while installing on device")
		}
	}

	return "", nil
}

func (a *App) ChangePackOrder(uuid uuid.UUID, index int) (string, error) {
	defer HandlePanic()
	log.Info("Change pack order")

	log.Info("Moving ", uuid, index)

	device, err := lunii.GetDevice()
	if err != nil {
		log.Error(err)
		log.Error(err)
		return "", err
	}

	err = device.ChangePackOrder(uuid, index)
	if err != nil {
		log.Error(err)
		log.Error(err)
		return "", err
	}

	return "", nil
}

func (a *App) SyncLuniiStoreMetadata(uuids []uuid.UUID) (string, error) {
	defer HandlePanic()
	log.Info("Sync store metadata")

	device, err := lunii.GetDevice()
	if err != nil {
		log.Error(err)
		return "", err
	}

	db, err := lunii.GetLuniiMetadataDb()
	if err != nil {
		log.Error(err)
		return "", err
	}

	for _, thisUuid := range uuids {
		device.SyncMetadataFromDb(thisUuid, db)
	}

	return "", nil
}

func (a *App) SyncStudioMetadata(uuids []uuid.UUID, dbPath string) (string, error) {
	defer HandlePanic()
	log.Info("Sync studio metadata")

	device, err := lunii.GetDevice()
	if err != nil {
		log.Error(err)
		return "", err
	}

	db, err := lunii.GetStudioMetadataDb(dbPath)
	if err != nil {
		log.Error(err)
		return "", err
	}

	for _, thisUuid := range uuids {
		device.SyncMetadataFromDb(thisUuid, db)
	}

	return "", nil
}

type CheckUpdateResponse struct {
	CanUpdate     bool   `json:"canUpdate"`
	LatestVersion string `json:"latestVersion"`
	ReleaseNotes  string `json:"releaseNotes"`
}

var lastCheck time.Time
var lastResponse *CheckUpdateResponse

func (a *App) CheckForUpdate() (*CheckUpdateResponse, error) {
	defer HandlePanic()
	if lastCheck.Add(time.Hour*1).Before(time.Now()) && lastResponse != nil {
		return lastResponse, nil
	}

	log.Info("Check updates")

	latest, found, err := selfupdate.DetectLatest("olup/lunii-admin")

	trimmedVersion := strings.TrimPrefix(version, "v")
	trimmedVersion = strings.TrimSuffix(trimmedVersion, "\n")

	if trimmedVersion == "next" {
		return &CheckUpdateResponse{
			CanUpdate:     false,
			LatestVersion: "",
			ReleaseNotes:  "",
		}, nil
	}

	v := semver.MustParse(trimmedVersion)

	if err != nil {
		log.Error(err)
		log.Error("Error occurred while detecting version:", err)
		return nil, err
	}

	if !found || latest.Version.LTE(v) {
		log.Info("Current version is the latest")
		return &CheckUpdateResponse{
			CanUpdate:     false,
			LatestVersion: "",
			ReleaseNotes:  "",
		}, nil
	}

	lastCheck = time.Now()
	lastResponse = &CheckUpdateResponse{
		CanUpdate:     true,
		LatestVersion: latest.Version.String(),
		ReleaseNotes:  latest.ReleaseNotes,
	}

	return lastResponse, nil
}

type Infos struct {
	Version   string `json:"version"`
	MachineId string `json:"machineId"`
	Os        string `json:"os"`
	Arch      string `json:"arch"`
}

func (a *App) GetInfos() (*Infos, error) {
	defer HandlePanic()
	trx := nrApp.StartTransaction("lunii_get_infos")
	defer trx.End()
	log.Info("Get infos")

	return &Infos{
		Version:   version,
		MachineId: machineId,
		Os:        goruntime.GOOS,
		Arch:      goruntime.GOARCH,
	}, nil
}
