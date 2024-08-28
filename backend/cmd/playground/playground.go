package main

import (
	"log/slog"
	"net/http"
	"os"
	"party-game/pkg/gamelogic"
	"party-game/pkg/handlers"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomePageHandler)
	mux.HandleFunc("/create-player", handlers.CreatePlayerHandler)
	mux.HandleFunc("/create-game", handlers.CreateGameHandler)
	mux.HandleFunc("/join-game", handlers.JoinGameHandler)
	loggedMux := logRequest(mux)

	slog.Info("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		slog.Error("error", err)
	}

	_, _ = gamelogic.CreateGame("asdf")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Received http request", "RemoteAddress", r.RemoteAddr, "Method", r.Method, "URL", r.URL, "Body", r.Body, "PostForm", r.PostForm)
		handler.ServeHTTP(w, r)
	})
}
