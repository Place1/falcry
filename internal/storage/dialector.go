package storage

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func dialector(u *url.URL) gorm.Dialector {
	switch u.Scheme {
	case "postgresql":
		u.Scheme = "postgres"
		fallthrough
	case "postgres":
		return postgres.Open(pgconn(u))
	case "sqlite3":
		return sqlite.Open(sqlite3conn(u))
	default:
		logrus.Panicf("unknown sql storage backend %s", u.Scheme)
		return nil // unreachable
	}
}

func pgconn(u *url.URL) string {
	password, _ := u.User.Password()
	decodedQuery, err := url.QueryUnescape(u.RawQuery)
	if err != nil {
		logrus.Warnf("failed to unescape connection string query parameters - they will be ignored")
		decodedQuery = ""
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		u.Hostname(),
		u.Port(),
		u.User.Username(),
		password,
		strings.TrimLeft(u.Path, "/"),
		decodedQuery,
	)
}

func sqlite3conn(u *url.URL) string {
	return filepath.Join(u.Host, u.Path)
}
