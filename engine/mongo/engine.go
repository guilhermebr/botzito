package mongo

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/guilhermebr/botzito/storage/mongodb"
)

type mongoEngine struct {
	name     string
	language string
	db       *mongodb.DB
}

func NewMongoEngine(name string, language string, data map[string]interface{}) (*mongoEngine, error) {
	endpoint, ok := data["endpoint"].(string)
	if !ok {
		return nil, errors.New("missing mongo endpoint")
	}
	database, ok := data["database"].(string)
	if !ok {
		database = name
	}
	db, err := mongodb.New(endpoint, database)
	if err != nil {
		return nil, err
	}
	return &mongoEngine{name: name, language: language, db: db}, nil
}

func (b *mongoEngine) Learn(data []byte) string {
	var i Intent
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err.Error()
	}

	fmt.Printf("intent: %#v\n", i)

	err = b.UpsertIntent(i)
	if err != nil {
		return err.Error()
	}

	return "success"
}

func (b *mongoEngine) Ask(question string) string {
	fmt.Println(question)
	resps, err := b.AskQuestion(question)
	if err != nil {
		return err.Error()
	}

	fmt.Printf("resps: %v\n", resps)

	return resps[0].Responses[0]
}
