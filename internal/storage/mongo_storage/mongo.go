package mongo_storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	db *mongo.Database
}

func NewStorage(storagePath string, storageName string) MongoStorage {
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
	mongo := MongoStorage{
		db: client.Database(storageName),
	}
	err = mongo.createCollectionLog()
	if err != nil {
		fmt.Printf("logs collection is not created: %s", err.Error())
		os.Exit(1)
	}
	mongo.createCollectionUrl()
	return mongo
}

func (sg *MongoStorage) createCollectionLog() error {
	command := bson.D{{Key: "create", Value: "logs"}}
	var res bson.D
	return sg.db.RunCommand(context.TODO(), command).Decode(&res)
}

func (sg *MongoStorage) createCollectionUrl() error {
	command := bson.D{{Key: "create", Value: "url"}}
	var res bson.D
	err := sg.db.RunCommand(context.TODO(), command).Decode(&res)

	index := mongo.IndexModel{
		Keys:    bson.M{"alias": 1},
		Options: options.Index().SetUnique(true),
	}
	sg.db.Collection("url").Indexes().CreateOne(context.Background(), index)
	return err
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
type UrlEntry struct {
	Alias string `json:"alias"`
	Url   string `json:"url"`
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
