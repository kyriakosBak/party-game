package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"party-game/pkg/gamelogic"
	// "party-game/pkg/gamelogic"
)

func RoundQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// get gameId
	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		slog.Error("Error while fetching game during round", "error", err)
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}

	round := gamelogic.GetLatestRound(gameId.Value)

	tmpl := template.Must(template.ParseFiles("templates/round-question.html"))
	tmpl.Execute(w, round)
	slog.Debug("Serving round question template", "round", round)
}
