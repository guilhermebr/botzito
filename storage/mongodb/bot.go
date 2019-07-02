package mongodb

import (
	"context"
	"fmt"

	"github.com/guilhermebr/botzito/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var _ types.BotStorage = &botStorage{}

const (
	botCollectionName = "bots"
)

type botStorage struct {
	db  *DB
	ctx context.Context
}

type bot struct {
	Name       string `bson:"name"`
	Language   string
	EngineType string                 `bson:"engine_type" json:"engine_type"`
	EngineData map[string]interface{} `bson:"engine_data" json:"engine_data"`
}

func NewBotStorage(db *DB) *botStorage {
	ctx := context.Background()
	return &botStorage{
		db:  db,
		ctx: ctx,
	}
}

func (s *botStorage) collection() *mongo.Collection {
	c := s.db.Session.Database(s.db.Database).Collection(botCollectionName)
	_, err := c.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{"name", bsonx.Int32(1)},
			},
			Options: options.Index().
				SetUnique(true),
		})

	if err != nil {
		fmt.Printf("error creating index - err: %v", err)
	}
	return c
}

func (s *botStorage) Create(form *types.Bot) error {
	c := s.collection()

	b := bot{
		Name:       form.Name,
		Language:   form.Language,
		EngineType: form.EngineType,
		EngineData: form.EngineData,
	}
	_, err := c.InsertOne(s.ctx, b)
	return err
}

func (s *botStorage) GetById(id string) (*types.Bot, error) {
	c := s.collection()
	var b bot
	filter := bson.M{"name": id}
	err := c.FindOne(s.ctx, filter).Decode(&b)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = types.DataNotFound
		}
		return nil, err
	}
	form := types.Bot{
		Name:       b.Name,
		Language:   b.Language,
		EngineType: b.EngineType,
		EngineData: b.EngineData,
	}
	return &form, err
}

func (s *botStorage) ListAll(limit int64, skip int64) ([]*types.Bot, error) {
	filter := bson.M{}
	opts := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}

	return s.listByFilter(filter, opts)
}

func (s *botStorage) listByFilter(filter bson.M, opts options.FindOptions) ([]*types.Bot, error) {
	c := s.collection()
	cursor, err := c.Find(s.ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	var bots []*types.Bot
	for cursor.Next(s.ctx) {
		b := bot{}
		err := cursor.Decode(&b)
		if err != nil {
			return nil, err
		}
		bots = append(bots, &types.Bot{
			Name:       b.Name,
			Language:   b.Language,
			EngineType: b.EngineType,
			EngineData: b.EngineData,
		})
	}

	return bots, nil
}
