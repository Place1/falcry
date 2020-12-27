package storage

import (
	"net/url"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MustOpen(rawurl string) *gorm.DB {
	u, err := url.Parse(rawurl)
	if err != nil {
		logrus.Panic(err)
	}

	db, err := Open(u)
	if err != nil {
		logrus.Panic(err)
	}

	return db
}

func Open(u *url.URL) (*gorm.DB, error) {
	logrus.Infof("storing data in sql backend: %s", u.Scheme)

	db, err := gorm.Open(dialector(u), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Event{})

	return db, nil
}
