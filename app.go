package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"github.com/olup/lunii-cli/pkg/lunii"
	studiopackbuilder "github.com/olup/lunii-cli/pkg/pack-builder"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/getsentry/sentry-go"
)

//go:generate bash scripts/get-version.sh
//go:embed version.txt
var version string
var machineId, _ = machineid.ID()

// App struct
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Sentry config
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://1c80a20fce80413db1008a4d4dc4a06d@o1341821.ingest.sentry.io/6615295",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
		Environment:      "production",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("infos", map[string]interface{}{
			"version":   version,
			"machineId": machineId,
		})
	})

	sentry.CaptureMessage("initial-event")
}

func (a *App) GetDeviceInfos() *lunii.Device {
	defer sentry.Recover()

	device, err := lunii.GetDevice()
	if err != nil {
		return nil
	}
	return device
}

func (a *App) ListPacks() []lunii.Metadata {
	defer sentry.Recover()

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
	defer sentry.Recover()

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
	defer sentry.Recover()

	_, err := studiopackbuilder.CreateStudioPack(directoryPath, destinationPath)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (a *App) OpenDirectory(title string) string {
	defer sentry.Recover()

	path, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	return path
}

func (a *App) SaveFile(title string, defaultDirectory string, defaultFileName string) string {
	defer sentry.Recover()

	fmt.Println("Select save path - options : ", defaultDirectory, defaultFileName)
	path, _ := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            title,
		DefaultDirectory: defaultDirectory,
		DefaultFilename:  defaultFileName,
	})
	return path
}

func (a *App) OpenFile(title string) string {
	defer sentry.Recover()

	path, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	return path
}

func (a *App) InstallPack(packPath string) (string, error) {
	defer sentry.Recover()

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
	defer sentry.Recover()

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
	defer sentry.Recover()

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
	defer sentry.Recover()

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

type CheckUpdateResponse struct {
	CanUpdate     bool   `json:"canUpdate"`
	LatestVersion string `json:"latestVersion"`
	ReleaseNotes  string `json:"releaseNotes"`
}

func (a *App) CheckForUpdate() (*CheckUpdateResponse, error) {
	defer sentry.Recover()

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
		log.Println("Error occurred while detecting version:", err)
		return nil, err
	}

	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return &CheckUpdateResponse{
			CanUpdate:     false,
			LatestVersion: "",
			ReleaseNotes:  "",
		}, nil
	}

	return &CheckUpdateResponse{
		CanUpdate:     true,
		LatestVersion: latest.Version.String(),
		ReleaseNotes:  latest.ReleaseNotes,
	}, nil
}

type Infos struct {
	Version   string `json:"version"`
	MachineId string `json:"machineId"`
	Os        string `json:"os"`
	Arch      string `json:"arch"`
}

func (a *App) GetInfos() (*Infos, error) {
	defer sentry.Recover()

	return &Infos{
		Version:   version,
		MachineId: machineId,
		Os:        goruntime.GOOS,
		Arch:      goruntime.GOARCH,
	}, nil
}
