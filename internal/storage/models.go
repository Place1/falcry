package storage

import (
	"time"

	"gorm.io/datatypes"
)

// Event is a model that stores the data from the original
// Falco event in the database.
// All fields are "create + read only" because they should
// be immutable once written to the database.
type Event struct {
	ID       uint           `gorm:"primaryKey"`
	Output   string         `gorm:"<-:create"`
	Priority string         `gorm:"<-:create"`
	Rule     string         `gorm:"<-:create"`
	Time     time.Time      `gorm:"<-:create"`
	Raw      datatypes.JSON `gorm:"<-:create"`
}
