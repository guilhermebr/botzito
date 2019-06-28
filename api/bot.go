package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guilhermebr/botzito/core"
	"github.com/guilhermebr/botzito/types"
	"github.com/sirupsen/logrus"
)

// title: create bot
// path: /bot
// method: POST
// responses:
//		201: created
//		400: bad request
//		500: server error
func (s *Service) createBot(w http.ResponseWriter, r *http.Request) {
	log := s.log.WithFields(logrus.Fields{
		"handler": "createBot",
	})
	log.Infoln("called")

	var bot types.Bot

	if err := json.NewDecoder(r.Body).Decode(&bot); err != nil {
		log.Error(err)
		ErrInvalidJson.Send(w)
		return
	}

	if bot.Name == "" {
		log.Error("missing name")
		respErr := ErrMissingData
		respErr.Message = "missing bot name"
		respErr.Send(w)
		return
	}

	err := s.storage.Bots.Create(&bot)
	if err != nil {
		log.Errorf("error creating bot on storage: %v", err)
		ErrInternalServer.Send(w)
		return
	}

	Success(bot, http.StatusCreated).Send(w)
}

// title: list bots
// path: /bot
// method: GET
// responses:
//		200: ok
//		400: bad request
//		500: server error
func (s *Service) listBot(w http.ResponseWriter, r *http.Request) {
	log := s.log.WithFields(logrus.Fields{
		"handler": "listBot",
	})
	log.Infoln("called")

	//TODO(guilhermebr) get limit and skip in querystring

	bots, err := s.storage.Bots.ListAll(0, 0)
	if err != nil {
		log.Errorf("error listing bots on storage: %v", err)
		ErrInternalServer.Send(w)
		return
	}

	Success(bots, http.StatusOK).Send(w)
}

// title: bot command
// path: /bot/{id}/cmd
// method: POST
// responses:
//		200: ok
//		400: bad request
//		500: server error
func (s *Service) botCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	botID := vars["id"]

	log := s.log.WithFields(logrus.Fields{
		"handler": "botCommand",
		"bot_id":  botID,
	})
	log.Infoln("called")

	// Get Bot
	bot, err := s.storage.Bots.GetById(botID)
	if err != nil {
		log.Errorf("error getting bot on storage: %v", err)
		ErrInternalServer.Send(w)
		return
	}

	var command types.BotCommand

	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		log.Error(err)
		ErrInvalidJson.Send(w)
		return
	}

	botEngine, err := core.LoadBot(bot)
	if err != nil {
		log.Errorf("error loading bot engine: %v", err)
		ErrInternalServer.Send(w)
		return
	}

	resp, err := botEngine.Run(command)
	if err != nil {
		log.Errorf("error executing bot command: %v", err)
		ErrInternalServer.Send(w)
		return
	}

	Success(resp, http.StatusOK).Send(w)
}
