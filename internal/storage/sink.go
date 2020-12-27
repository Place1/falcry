package storage

import (
	"context"

	"github.com/place1/falcry/internal/webhook"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StorageSink struct {
	db          *gorm.DB
	broadcaster *webhook.Broadcaster
}

func NewStorageSink(db *gorm.DB, broadcaster *webhook.Broadcaster) *StorageSink {
	return &StorageSink{db, broadcaster}
}

func (s *StorageSink) ListenAndSave() {
	for event := range s.broadcaster.Channel(context.Background()) {
		model := &Event{
			Output:   event.Output,
			Rule:     event.Rule,
			Priority: event.Priority,
			Time:     event.Time,
			Raw:      datatypes.JSON(event.Raw),
		}
		if err := s.db.Save(model).Error; err != nil {
			logrus.Error(err)
		}
	}
}
