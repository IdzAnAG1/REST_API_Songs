package main

import (
	_ "REST_API_Songs/docs"
	"REST_API_Songs/internal/config"
	"REST_API_Songs/internal/db"
	"REST_API_Songs/internal/handlers"
	"REST_API_Songs/internal/repository"
	"REST_API_Songs/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	Swag "github.com/swaggo/http-swagger"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}
	database, err := db.NewDataBase(cfg.DataBaseURL)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()
	repo := repository.NewRepository(database)
	svc := service.NewService(repo, cfg.ExternalApiURL)
	h := handlers.NewHandler(svc)
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/list", http.StatusMovedPermanently)
	}).Methods("GET")
	r.HandleFunc("/list", h.GetSongs).Methods("GET")
	r.HandleFunc("/song/{id}/text", h.GetSongText).Methods("GET")
	r.HandleFunc("/song/{id}", h.DeleteSong).Methods("DELETE")
	r.HandleFunc("/song/{id}", h.UpdateSong).Methods("PUT")
	r.HandleFunc("/song", h.CreateSong).Methods("POST")
	r.PathPrefix("/swagger/").Handler(Swag.WrapHandler)
	logrus.Infof("Service work on http://localhost:%s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
