package driver

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	LifeCycleNil = iota
	LifeCycleOneDay
	LifeCycleOneWeek
	LifeCycleOneMonth
	LifeCycleOneYear
	LifeCycleForever
)

type CodeSegmentRecord struct {
	Title     string   `bson:"title" json:"title"`
	Content   string   `bson:"content" json:"content"`
	Author    string   `bson:"author" json:"author"`
	Lang      string   `bson:"lang" json:"lang"`
	Password  string   `bson:"password" json:"-"`
	ShortKey  string   `bson:"short_key" json:"sk"`
	Tags      []string `bson:"tags" json:"tags"`
	Lifecycle uint8    `bson:"lifecycle" json:"lifecycle"`
	CreatedAt int64    `bson:"created_at" json:"created_at"`
	UpdatedAt int64    `bson:"updated_at" json:"updated_at"`
	ExpiredAt int64    `json:"expired_at"`
	Editable  bool     `json:"editable"`
}

func NewMongoClient() *mongo.Collection {
	return GetCollection()
}

func GetCollection() *mongo.Collection {
	return client.Database(confer.Mongo.Dbname).Collection("pastebin")
}
