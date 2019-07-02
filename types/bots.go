package types

var (
	LearnCommand commands = "learn"
)

type commands string

type Bot struct {
	Name       string
	Language   string
	EngineType string                 `json:"engine_type"`
	EngineData map[string]interface{} `json:"engine_data"`
}

type BotCommand struct {
	Command commands
	Data    interface{}
}

type BotStorage interface {
	Create(*Bot) error
	GetById(string) (*Bot, error)
	ListAll(limit int64, skip int64) ([]*Bot, error)
}
