package lunii

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buger/jsonparser"
	"github.com/google/uuid"
)

func getStudioDbPath() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, ".studio", "db", "unofficial.json")
}

func GetStudioMetadataDb(dbPath string) (*Db, error) {
	db := Db{
		stories: []Story{},
	}

	// if path is empty get default db path
	if dbPath == "" {
		dbPath = getStudioDbPath()
	}

	fmt.Println("Reading Studio db from ", dbPath)

	// read DB from STUdio local DB
	dbBytes, err := os.ReadFile(dbPath)

	// if no file, return an error
	if err != nil {
		fmt.Println("The db file can't be read")
		return nil, err
	}

	// parse db's json TODO
	jsonparser.ObjectEach(dbBytes, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		storyUuid, err := jsonparser.GetString(value, "uuid")
		if err != nil {
			return err
		}

		title, _ := jsonparser.GetString(value, "title")
		description, _ := jsonparser.GetString(value, "description")

		db.stories = append(db.stories, Story{
			Uuid:  uuid.MustParse(storyUuid),
			Title: title, Description: description,
			PackType: "custom",
		})
		return nil
	})

	return &db, nil
}
