package mongo_storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"os"
	"time"
	"url-shortener/internal/storage"
)

type MongoStorage struct {
	db *mongo.Database
}

func NewStorage(storagePath string, storageName string) (storage.Storage, error) {
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
		return nil, err
	}
	mongo := &MongoStorage{
		db: client.Database(storageName),
	}
	err = mongo.createCollectionLog()
	if err != nil {
		fmt.Printf("logs collection is not created: %s", err.Error())
		return nil, err
	}
	err = mongo.createCollectionUrl()
	if err != nil {
		fmt.Printf("urls collection is not created: %s", err.Error())
	}
	return mongo, nil
}

func (ms *MongoStorage) createCollectionLog() error {
	command := bson.D{{Key: "create", Value: "logs"}}
	var res bson.D
	return ms.db.RunCommand(context.TODO(), command).Decode(&res)
}

func (ms *MongoStorage) createCollectionUrl() error {
	command := bson.D{{Key: "create", Value: "urls"}}
	var res bson.D
	err := ms.db.RunCommand(context.TODO(), command).Decode(&res)
	if err != nil {
		return err
	}

	index := mongo.IndexModel{
		Keys:    bson.M{"alias": 1},
		Options: options.Index().SetUnique(true),
	}
	ms.db.Collection("urls").Indexes().CreateOne(context.Background(), index)
	return err
}

type LogEntry struct {
	Time string `json:"time"`
	Lvl  string `json:"lvl"`
	Msg  string `json:"msg"`
}

func (ms *MongoStorage) Log(msg string, lvl slog.Level) {
	l := LogEntry{
		Time: time.Now().Format(time.RFC3339),
		Lvl:  LogLevelFrom(lvl),
		Msg:  msg,
	}
	_, err := ms.db.Collection("logs").InsertOne(context.Background(), l)
	if err != nil {
		fmt.Errorf("Err: ujas log")
	}
}

func (ms *MongoStorage) GetURL(alias string) (string, error) {
	coll := ms.db.Collection("urls")
	filter := bson.D{{Key: "alias", Value: alias}}
	var urlEntry UrlEntry
	err := coll.FindOne(context.TODO(), filter).Decode(&urlEntry)
	if err != nil {
		return "", err
	}
	return urlEntry.Url, nil
}

func (ms *MongoStorage) SaveURL(urlToSave, alias string) error {
	u := UrlEntry{
		Url:   urlToSave,
		Alias: alias,
	}
	_, err := ms.db.Collection("urls").InsertOne(context.Background(), u)
	return err
}

func (ms *MongoStorage) AliasExist(alias string) bool {
	coll := ms.db.Collection("urls")
	filter := bson.D{{Key: "alias", Value: alias}}
	err := coll.FindOne(context.Background(), filter)
	if err != nil {
		return false
	}
	return true
}

func (ms *MongoStorage) GetAlias(url string) string {
	coll := ms.db.Collection("urls")
	filter := bson.D{{Key: "url", Value: url}}
	var urlEntry UrlEntry
	err := coll.FindOne(context.Background(), filter).Decode(urlEntry)
	if err != nil {
		return ""
	}
	return urlEntry.Alias
}

type UrlEntry struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

func LogLevelFrom(lvl slog.Level) string {
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
