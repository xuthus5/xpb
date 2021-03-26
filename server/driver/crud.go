package driver

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	LifeCycleOneDay = iota
	LifeCycleThreeDay
	LifeCycleOneWeek
	LifeCycleOneMonth
	LifeCycleSixMonth
	LifeCycleOneYear
	LifeCycleForever
)

type CodeSegmentRecord struct {
	Title     string   `bson:"title" json:"title"`
	Content   string   `bson:"content" json:"content"`
	Author    string   `bson:"author" json:"author"`
	Lang      string   `bson:"lang" json:"lang"`
	Tags      []string `bson:"tags" json:"tags"`
	Lifecycle uint8    `bson:"lifecycle" json:"lifecycle"`
	ShortKey  string   `bson:"short_key" json:"short_key"`
	CreatedAt int64    `bson:"created_at" json:"created_at"`
	UpdatedAt int64    `bson:"updated_at" json:"updated_at"`
}

func NewMongoClient() *mongo.Collection {
	return GetCollection()
}

func GetCollection() *mongo.Collection {
	return client.Database(confer.Mongo.Dbname).Collection("pastebin")
}
