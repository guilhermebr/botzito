# Botzito

# Install from Source

```
$ git clone git@github.com/guilhermebr/botzito
$ cd botzito
$ make compile
```

# Running

```
$ ./cmd/botzito
```

# Creating a Bot

```
$ curl -i 127.0.0.1:5000/bot -d '{
	"name": "teste-bot",
	"language": "english",
	"engine_type": "mongo",
	"engine_data": {
		"endpoint": "mongodb://localhost:27017",
		"database": "teste_bot"
	}
}'
```

# Creating a new Intent

```
$ curl -i 127.0.0.1:5000/bot/teste-bot -d '{
	"command": "learn",
  "data": {
  	"tag": "hello",
   	"patterns": [
    	"Hy",
      "Hello",
      "Whats up?"
    ],
    "responses": [
    	"Hy! how are you?",
      "Hello :)"
    ]
  }
}'
```

# Asking to bot

```
$ curl -i 127.0.0.1:5000/bot/teste-bot -d '{
	"command": "ask",
  "data": "Hello my friend!!"
}'
```

# Features

- [x] Engine: Mongo
- [ ] Engine/Mongo: Entity recognizer
- [ ] Engine: Tensorflow
- [ ] Engine: Dialogflow
- [ ] Platforms: Slack
- [ ] Platforms: Telegram
