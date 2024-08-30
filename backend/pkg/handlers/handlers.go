package handlers

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"party-game/pkg/gamelogic"
)

const playerIdCookie string = "player-id"
const questionIdCookie string = "question-id"

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("Cannot parse form.", "error", err)
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}

	playerName := r.FormValue("player-name")
	if playerName == "" {
		slog.Error("Player name empty during create player")
		http.Error(w, "Player name empty. Try again.", http.StatusInternalServerError)
		return
	}

	player, ok := gamelogic.CreatePlayer(playerName)
	if !ok {
		http.Error(w, "Player with this name already exists.", http.StatusInternalServerError)
		return
	}

	// Set the player ID in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  playerIdCookie,
		Value: player.Id,
		Path:  "/",
	})

	io.WriteString(w, "Player created succesfully!")
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("Cannot parse form.", "error", err)
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}

	password := r.FormValue("game-password")
	if password == "" {
		slog.Error("Empty password in create game request", "Form", r.PostForm)
		http.Error(w, "Cannot create game with empty password", http.StatusUnprocessableEntity)
		return
	}

	_, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Player not identified. Make sure you have created one.", http.StatusForbidden)
		return
	}

	game, created := gamelogic.CreateGame(password)
	if !created {
		http.Error(w, "Could not create game. Probably a game with the same password is already running", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/round-question.html"))
	tmpl.Execute(w, game)
	slog.Debug("Serving game template.", "template", tmpl)
}

func JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("Cannot parse form.", "error", err)
		http.Error(w, "Error. Check server logs.", http.StatusInternalServerError)
		return
	}
	password := r.FormValue("game-password")
	if password == "" {
		slog.Error("Empty password in create game request", "Form", r.PostForm)
		http.Error(w, "Cannot create game with empty password", http.StatusUnprocessableEntity)
		return
	}

	playerId, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Player not identified. Make sure you have created one.", http.StatusForbidden)
		return
	}

	game, err := gamelogic.JoinGame(password, playerId.Value)
	if err != nil {
		http.Error(w, "Could not join game. Check server logs", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/round-question.html"))
	tmpl.Execute(w, game)
	slog.Debug("Serving game template.", "template", tmpl)
}

func SumbitReplyGameHandler(w http.ResponseWriter, r *http.Request) {

}

func IsPost(r *http.Request) bool {
	switch r.Method {
	case "POST":
		return true

	default:
		slog.Debug("Request is not POST.", "RemoteAddress", r.RemoteAddr, "Method", r.Method, "URL", r.URL)
	}
	return false
}

func ServeGame(w http.ResponseWriter, r *http.Request, game gamelogic.Game) {
	tmpl := template.Must(template.ParseFiles("templates/round-question.html"))
	tmpl.Execute(w, game)
	slog.Debug("Serving game template.", "template", tmpl)
}
