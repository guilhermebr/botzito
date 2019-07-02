package core

import (
	"encoding/json"
	"fmt"

	"github.com/guilhermebr/botzito/engine/mongo"
	"github.com/guilhermebr/botzito/types"
)

type BotCore struct {
	Name       string
	Language   string
	EngineType string
	EngineData map[string]interface{}
}

type BotEngine interface {
	Learn([]byte) string
	Ask(string) string
}

func LoadBot(botData *types.Bot) (*BotCore, error) {
	return &BotCore{
		Name:       botData.Name,
		Language:   botData.Language,
		EngineType: botData.EngineType,
		EngineData: botData.EngineData,
	}, nil
}

func (b *BotCore) Run(cmd types.BotCommand) (string, error) {
	var resp string
	var err error

	var engine BotEngine

	//load engine
	if b.EngineType == "mongo" {
		engine, err = mongo.NewMongoEngine(b.Name, b.Language, b.EngineData)
		if err != nil {
			return "", err
		}
	}

	switch cmd.Command {
	case types.LearnCommand:
		fmt.Println("Learn")
		fmt.Println(cmd.Data)
		b, err := json.Marshal(cmd.Data)
		if err != nil {
			return "", err
		}
		resp = engine.Learn(b)

	default:
		fmt.Println(cmd)
		resp = engine.Ask(cmd.Data.(string))
	}
	return resp, nil
}
