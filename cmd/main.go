package main

import (
	"os"

	"github.com/guilhermebr/botzito/api"
	"github.com/guilhermebr/botzito/storage"
	"github.com/guilhermebr/botzito/storage/mongodb"
	"github.com/sirupsen/logrus"
)

func main() {

	// start log
	log := logrus.StandardLogger()
	log.Infoln("Starting account api...")

	// Start Storages
	endpoint := os.Getenv("MONGODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "botzito"
	}
	log.Infof("connecting to mongodb at %s database %s", endpoint, dbName)
	db, err := mongodb.New(endpoint, dbName)
	if err != nil {
		log.Fatal(err)
	}

	storage := &storage.Storage{
		Bots: mongodb.NewBotStorage(db),
	}

	// Start Platforms

	// Start API
	if err := api.Start(log, storage); err != nil {
		log.WithError(err).Fatalln("Couldn't start api!")
	}
}
