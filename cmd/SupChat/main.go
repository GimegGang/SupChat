package main

import (
	"SupChat/internal/Config"
	"SupChat/internal/Handlers/AdminHandler"
	"SupChat/internal/Handlers/UserHandler"
	"SupChat/internal/Storage/SQLite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	config := Config.MustLoad("./config.yaml")
	log.Info("Load config file success", "config", config)
	db, err := SQLite.New(config.StoragePath)
	if err != nil {
		log.Error("Create db error", "error", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)

	router.Get("/", UserHandler.MainHandler(log, db))
	router.Post("/sendUserMessage", UserHandler.SendMessageHandler(log, db))
	router.Get("/admin", AdminHandler.MainHandler(log, db))
	router.Get("/admin/{chatId}", AdminHandler.ChatHandler(log, db))
	router.Post("/admin/{chatId}/send", AdminHandler.SendAdminMessageHandler(log, db))
	router.Get("/admin/{chatId}/closeTicket", AdminHandler.CloseTicketHandler(log, db))
	router.Get("/getMessages/{ticketID}", UserHandler.GetMessageHandler(log, db))

	srv := http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	fs := http.FileServer(http.Dir("static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Listen: %s\n", err)
	}
}
