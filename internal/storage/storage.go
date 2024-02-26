package storage

import (
	"log/slog"
)

type Storage interface {
	Log(msg string, lvl slog.Level)
	SaveURL(urlToSave, alias string) error
	GetURL(alias string) (string, error)
	GetAlias(url string) string
	AliasExist(alias string) bool
}
