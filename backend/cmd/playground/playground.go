package main

import (
	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	"io"
	"log/slog"
	"net/http"
	"os"
	"party-game/pkg/handlers"
)

func main() {
	var logger *slog.Logger
	gelfWriter, err := gelf.NewWriter("localhost:12201")
	if err != nil {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		file, err := os.Create("playground.log")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		logger = slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), opts))
	} else {
		gelfWriter.CompressionType = gelf.CompressNone
		logger = slog.New(sloggraylog.Option{Level: slog.LevelDebug, Writer: gelfWriter}.NewGraylogHandler())
	}
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	handlers.AddHandlers(mux)
	loggedMux := logRequest(mux)

	slog.Info("Server is starting on port 8888...")
	if err := http.ListenAndServe(":8888", loggedMux); err != nil {
		slog.Error("error", err)
	}
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
