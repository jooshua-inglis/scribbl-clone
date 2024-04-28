package main

import (
	"log/slog"
	"net/http"
	"os"
	"scribl-clone/handlers"
	"scribl-clone/sockets"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func configureLogger() {
	var programLevel = new(slog.LevelVar)
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RealIP)

	configureLogger()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "ws://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/game/new", handlers.CreateGame)
	r.Post("/game/{gameId}/start", handlers.StartGame)
	r.Post("/game/{gameId}/join", handlers.JoinGame)
	r.Post("/game/{gameId}/select_word", handlers.SelectWord)
	r.Post("/game/{gameId}/guess", handlers.MakeGuess)
	r.Get("/game/{gameId}/dummy_event", handlers.DummyEvent)
	r.Get("/game/{gameId}/players", handlers.GetPlayers)
	r.Get("/game/{gameId}", handlers.GetGame)

	r.Get("/player/{playerId}", handlers.GetPlayer)
	r.Patch("/player/{playerId}", handlers.PatchPlayer)

	r.Get("/game_connection/{userId}", sockets.InitGameConnection)

	if err := http.ListenAndServe(":4000", r); err != nil {
		slog.Error(err.Error())
		return
	}
}
