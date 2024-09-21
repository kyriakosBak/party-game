package handlers

import "net/http"

func AddHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", HomePageHandler)
	mux.HandleFunc("/create-player", CreatePlayerHandler)
	mux.HandleFunc("/create-game", CreateGameHandler)
	mux.HandleFunc("/join-game", JoinGameHandler)
	mux.HandleFunc("/player-ready", PlayerReadyHandler)
	mux.HandleFunc("/round-question", RoundQuestionHandler)
	mux.HandleFunc("/submit-answer", SubmitAnswerHandler)
	mux.HandleFunc("/round-choice", RoundChoiceHandler)
	mux.HandleFunc("/submit-choice", SubmitChoiceHandler)
	mux.HandleFunc("/round-results", RoundResultsHandler)
	mux.HandleFunc("/new-round-ready", NewRoundReady)
}
