package lunii

import (
	"github.com/google/uuid"
)

type Story struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Uuid        uuid.UUID `json:"uuid"`
	PackType    string
}

type Db struct {
	stories []Story
}

func (db *Db) GetStoryById(uuid uuid.UUID) *Story {
	for _, story := range db.stories {
		if story.Uuid == uuid {
			return &story
		}
	}
	return nil
}
