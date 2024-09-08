package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"party-game/pkg/gamelogic"
	"time"
)

var timeoutTime time.Duration = 60 * time.Second

func RoundQuestionHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering RoundQuestion handler")
	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}

	round, err := gamelogic.GetLatestRound(gameId.Value)
	if err != nil {
		http.Error(w, "Could not get latest round", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  roundIdCookie,
		Value: round.Id,
		Path:  "/",
	})

	tmpl := template.Must(template.ParseFiles("templates/round-question.html"))
	tmpl.Execute(w, round)
	slog.Debug("Serving round question template", "round", round)
}

func SubmitAnswerHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering SubmitAnswer handler")
	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		http.Error(w, "Could not find game id cookie.", http.StatusBadRequest)
		return
	}

	playerId, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Could not find player id cookie.", http.StatusBadRequest)
		return
	}

	roundId, err := r.Cookie(roundIdCookie)
	if err != nil {
		http.Error(w, "Could not find round id cookie.", http.StatusBadRequest)
		return
	}

	answer := r.PostFormValue("player-answer")

	err = gamelogic.AddAnswer(gameId.Value, playerId.Value, roundId.Value, answer)
	if err != nil {
		http.Error(w, "Could not add answer. Check server logs", http.StatusInternalServerError)
		slog.Error("Could not add answer", "error", err)
		return
	}

	controlTime := time.Now()
	for {
		time.Sleep(time.Second)
		if gamelogic.AllPlayerAnswered(gameId.Value, roundId.Value) {
			break
		}
		if time.Since(controlTime) > timeoutTime {
			break
		}
	}

	w.Header().Set("HX-Redirect", "/round-choice")
	w.Write(nil)
	slog.Debug("Redirect to /round-choice")
}

type RoundChoiceData struct {
	Question string
	Choices  []gamelogic.Answer
}

func RoundChoiceHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering RoundChoice handler")
	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	playerId, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Could not find player id cookie.", http.StatusBadRequest)
		return
	}

	round, err := gamelogic.GetLatestRound(gameId.Value)
	if err != nil {
		http.Error(w, "Could not get latest round", http.StatusInternalServerError)
		return
	}
	answersCopy := []gamelogic.Answer{}
	for _, a := range round.Answers {
		if a.Owner.Id != playerId.Value {
			answersCopy = append(answersCopy, a)
		}
	}
	responseData := RoundChoiceData{round.Question, answersCopy}

	tmpl := template.Must(template.ParseFiles("templates/round-choices.html"))
	tmpl.Execute(w, responseData)
	slog.Debug("Serving round choice template", "responseData", responseData)
}

func SubmitChoiceHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering SubmitChoice handler")
	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		http.Error(w, "Could not find game id cookie.", http.StatusBadRequest)
		return
	}

	playerId, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Could not find player id cookie.", http.StatusBadRequest)
		return
	}

	roundId, err := r.Cookie(roundIdCookie)
	if err != nil {
		http.Error(w, "Could not find round id cookie.", http.StatusBadRequest)
		return
	}

	choiceId := r.PostFormValue("player-choice-id")
	if choiceId == "" {
		http.Error(w, "Could not find player choice id", http.StatusBadRequest)
		return
	}

	err = gamelogic.AddChoice(gameId.Value, playerId.Value, roundId.Value, choiceId)
	if err != nil {
		slog.Error("Could not add choice " + choiceId)
	}

	controlTime := time.Now()
	for {
		time.Sleep(time.Second)
		if gamelogic.AllPlayerSelectedChoice(gameId.Value, roundId.Value) {
			break
		}
		if time.Since(controlTime) > timeoutTime {
			break
		}
	}

	w.Header().Set("HX-Redirect", "/round-results")
	w.Write(nil)
	slog.Debug("Redirect to /round-results")
}

type RoundResultsData struct {
	Score []ScoreData
}

type ScoreData struct {
	PlayerName string
	Points     int
}

func RoundResultsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering RoundResults handler")

	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		http.Error(w, "Could not find game id cookie.", http.StatusBadRequest)
		return
	}

	score := gamelogic.GetScore(gameId.Value)
	scoreData := []ScoreData{}

	for k, v := range score {
		playerName := gamelogic.GetPlayer(k).Name
		scoreData = append(scoreData, ScoreData{playerName, v})
	}

	responseData := RoundResultsData{scoreData}

	tmpl := template.Must(template.ParseFiles("templates/round-results.html"))
	tmpl.Execute(w, responseData)
	slog.Debug("Serving round results template", "responseData", responseData)
}

func NewRoundReady(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering NewRoundReady handler")

	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		http.Error(w, "Could not find game id cookie.", http.StatusBadRequest)
		return
	}

	gamelogic.AddNewRound(gameId.Value)
}
