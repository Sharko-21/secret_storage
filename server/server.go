package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"main/logger"
	"main/storage"
)

type Server interface {
	Start() <-chan struct{}
}
type server struct {
	http.Server
	service service
}

type service struct {
	storage storage.Storage
}

func (s *server) SaveSecret(w http.ResponseWriter, r *http.Request) {
	logger.ErrorLog.Println(r.ParseForm())
	secret := r.Form.Get("secret")
	key := r.Form.Get("key")
	id := s.service.storage.Save(secret, key)
	logger.ErrorLog.Println(w.Write([]byte(id)))
}

func (s *server) GetSecret(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	id := r.URL.Query().Get("id")
	key := r.URL.Query().Get("key")
	secret := s.service.storage.Get(id, key)
	logger.ErrorLog.Println(w.Write([]byte(secret)))
}

func (s *server) Start() <-chan struct{} {
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.ErrorLog.Println(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	exit := make(chan struct{}, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		_ = <-sig
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.Shutdown(ctx)
		exit <- struct{}{}
		return
	}()
	return exit
}

func NewServer() Server {
	s := &server{
		Server: http.Server{
			Addr:         ":8082",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		service: service{
			storage: storage.GetStorage(),
		},
	}
	r := mux.NewRouter()
	r.HandleFunc("/", s.SaveSecret).Methods("POST")
	r.HandleFunc("/", s.GetSecret).Methods("GET")
	r.Use(logger.LogMiddleware)
	s.Handler = r
	return s
}
