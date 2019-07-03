package mongo

import (
	"context"
	"fmt"

	"github.com/guilhermebr/botzito/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Intent struct {
	Tag           string `bson:"id"`
	Patterns      []string
	Responses     []string
	ContextSet    string      `bson:"context_set" json:"context_set"`
	ContextFilter string      `bson:"context_filter" json:"context_filter"`
	Score         interface{} `bson:",omitempty" json:",omitempty"`
	//Action string
}

func (s *mongoEngine) collection() *mongo.Collection {
	c := s.db.Session.Database(s.db.Database).Collection("intents")
	_, err := c.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{"patterns", bsonx.String("text")},
				{"responses", bsonx.String("text")},
				{"tag", bsonx.String("text")},
			},
			Options: options.Index().
				SetWeights(bsonx.Doc{
					{"patterns", bsonx.Int32(3)},
					{"responses", bsonx.Int32(1)},
					{"tag", bsonx.Int32(3)},
				}).
				SetName("textIndex").
				SetDefaultLanguage(s.language),
		})

	if err != nil {
		fmt.Printf("error creating index - err: %v", err)
	}
	return c
}

func (s *mongoEngine) UpsertIntent(i Intent) error {
	c := s.collection()
	filter := bson.D{
		{"tag", i.Tag},
	}
	doc, err := mongodb.Atodoc(i)
	if err != nil {
		return err
	}
	query := bson.D{
		{"$set", doc},
	}

	_, err = c.UpdateOne(context.Background(), filter, query, options.Update().SetUpsert(true))
	return err
}

func (s *mongoEngine) AskQuestion(question string) ([]Intent, error) {
	c := s.collection()
	filter := bson.D{{"$text", bson.D{{"$search", question}}}}

	opts := options.FindOptions{
		Projection: bson.D{
			{"score", bson.D{
				{"$meta", "textScore"},
			}},
		},
		Sort: bson.D{
			{"score", bson.D{
				{"$meta", "textScore"},
			}},
		},
	}
	cursor, err := c.Find(context.Background(), filter, &opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var intents []Intent
	for cursor.Next(context.Background()) {
		var i Intent
		err := cursor.Decode(&i)
		if err != nil {
			return nil, err
		}
		intents = append(intents, i)
	}

	return intents, err
}
