package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	collection *mongo.Database
}

func newStorage()

type LogEntry struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

// func (sg *Storage) log(msg string, lvl int) *LogEntry {

//  	LogEntry{
// 		Time:  time.Now().Format(time.RFC3339),
// 		Level: logLevelFrom(lvl),
// 		Msg:   msg,
// 	}
// }

func logLevelFrom(lvl int) string {

	switch lvl {
	case 1:
		return "debug"
	case 2:
		return "info"
	case 3:
		return "error"
	}
	return "incorrect lvl"
}
