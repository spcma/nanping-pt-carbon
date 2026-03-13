package shared

import (
	"app/internal/shared/event"
	"app/internal/shared/persistence"
	"gorm.io/gorm"
)

var GlobalEventBus event.EventBus

var GlobalDB *gorm.DB

func Init(config persistence.Config) (*gorm.DB, error) {
	db, err := persistence.NewGORM(config)
	if err != nil {
		return nil, err
	}

	GlobalDB = db

	GlobalEventBus = event.NewInMemoryEventBus()

	return db, nil
}
