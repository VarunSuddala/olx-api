package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VarunSuddala/olx-api/internal/config"
	"github.com/VarunSuddala/olx-api/internal/db"
	"github.com/VarunSuddala/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	mux := http.NewServeMux()
	db, err := db.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("main.db.connect :%v", err)
	}
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"welcome":"olx-api"}`))
	})
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", handlers.Get_listings(db))
	mux.HandleFunc("DELETE /listings/{id}", handlers.Delete_listing(db))
	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	fmt.Println(("database connected"))
	fmt.Println(("server is running"))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed : %v", err)
	}

}
