package storage

import "log/slog"

type Storage interface {
	Log(msg string, lvl slog.Level)
}
