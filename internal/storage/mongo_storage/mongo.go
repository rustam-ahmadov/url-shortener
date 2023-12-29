package mongo_storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	db *mongo.Database
}

func NewStorage(storagePath string) MongoStorage {
	clientOptions := options.Client().ApplyURI(storagePath)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	mongo:= MongoStorage{
		db: client.Database("url-shortener"),
	}
	return mongo
}

func (sg MongoStorage) Log(msg string, lvl slog.Level) {
	l := LogEntry{
		Time: time.Now().Format(time.RFC3339),
		Lvl:  logLevelFrom(lvl),
		Msg:  msg,
	}
	sg.db.Collection("logs").InsertOne(context.Background(), l)
}

type LogEntry struct {
	Time string `json:"time"`
	Lvl  string `json:"lvl"`
	Msg  string `json:"msg"`
}

func logLevelFrom(lvl slog.Level) string {
	switch lvl {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelInfo:
		return "info"
	case slog.LevelError:
		return "error"
	}
	return "incorrect lvl"
}
