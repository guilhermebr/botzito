package types

var (
	LearnCommand commands = "learn"
)

type commands string

type Bot struct {
	Name       string
	EngineType string
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
