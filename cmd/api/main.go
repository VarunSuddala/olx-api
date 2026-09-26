package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/VarunSuddala/olx-api/internal/config"
	"github.com/VarunSuddala/olx-api/internal/db"
	"github.com/VarunSuddala/olx-api/internal/handlers"
	"github.com/VarunSuddala/olx-api/internal/middleware"
)

func main() {
	cfg := config.MustLoad()
	mux := http.NewServeMux()
	db, err := db.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("main.db.connect :%v", err)
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	lh := handlers.NewListingHandler(db, logger)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"welcome":"olx-api"}`))
	})
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.Get)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)
	mux.HandleFunc("POST /listings", lh.Create)
	handler := middleware.RequestId(mux)
	srv := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      handler,
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
