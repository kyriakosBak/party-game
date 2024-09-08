package main

import (
	"io"
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
	file, err := os.Create("playground.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), opts))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	handlers.AddHandlers(mux)
	loggedMux := logRequest(mux)

	slog.Info("Server is starting on port 8888...")
	if err := http.ListenAndServe(":8888", loggedMux); err != nil {
		slog.Error("error", err)
	}

	_, _ = gamelogic.CreateGame("asdf")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()
		}
		slog.Info("Received http request", "RemoteAddress", r.RemoteAddr, "Method", r.Method, "URL", r.URL, "Body", r.Body, "PostForm", r.PostForm)
		handler.ServeHTTP(w, r)
	})
}
