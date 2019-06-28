package core

import (
	"fmt"

	"github.com/guilhermebr/botzito/engine/mongo"
	"github.com/guilhermebr/botzito/types"
)

type BotCore struct {
	Name       string
	EngineType string
	EngineData map[string]interface{}
}

type BotEngine interface {
	Learn(map[string]interface{}) string
	Ask(string) string
}

func LoadBot(botData *types.Bot) (*BotCore, error) {
	return &BotCore{
		Name:       botData.Name,
		EngineType: botData.EngineType,
		EngineData: botData.EngineData,
	}, nil
}

func (b *BotCore) Run(cmd types.BotCommand) (string, error) {
	var resp string
	var err error

	var engine BotEngine

	//load engine
	if b.EngineType == "default" {
		engine, err = mongo.NewMongoEngine(b.Enginedata)
		if err != nil {
			return "", err
		}
	}

	switch cmd.Command {
	case types.LearnCommand:
		fmt.Println("Learn")
		fmt.Println(cmd.Data)
		resp = engine.Learn(cmd.Data.(map[string]interafce{}))

	default:
		fmt.Println(cmd)
		resp = engine.Ask(cmd.Data.(string))
	}
	return resp, nil
}
