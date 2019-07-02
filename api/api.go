package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/guilhermebr/botzito/storage"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type Service struct {
	log     *logrus.Logger
	storage *storage.Storage
	Cfg     *config
}

type config struct {
	Port      string
	SecretKey string
}

func (s *Service) LoadConfig() {
	cfg := config{}
	cfg.SecretKey = os.Getenv("SECRET_KEY")
	//	if len(cfg.SecretKey) == 0 {
	//		log.Fatal("SECRET_KEY env var is required")
	//	}
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "5000"
	}
	s.Cfg = &cfg
}

func Start(log *logrus.Logger, storage *storage.Storage) error {
	service := Service{
		log:     log,
		storage: storage,
	}
	service.LoadConfig()

	//Router
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", service.healthcheck).Methods("GET")
	//	r.HandleFunc("/login", service.login).Methods("POST")
	r.HandleFunc("/bot", service.createBot).Methods("POST")
	r.HandleFunc("/bot", service.listBot).Methods("GET")
	r.HandleFunc("/bot/{id}", service.botCommand).Methods("POST")
	//r.HandleFunc("/bot", middleware.ValidateTokenIfExists(service.Cfg.SecretKey, service.bot)).Methods("POST")

	//Negroni
	n := negroni.Classic()
	n.UseHandler(r)
	//n.UseHandler(middleware.Cors(r))

	service.log.Infoln("Listen at 0.0.0.0:" + service.Cfg.Port)
	http.ListenAndServe(":"+service.Cfg.Port, n)

	return nil
}
