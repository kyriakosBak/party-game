package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"party-game/pkg/gamelogic"
)

const playerIdCookie string = "player-id"
const gameIdCookie string = "game-id"
const roundIdCookie string = "round-id"

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering Home handler")
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering CreatePlayer handler")

	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("Cannot parse form.", "error", err)
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	playerName := r.FormValue("player-name")
	if playerName == "" {
		slog.Error("Player name empty during create player")
		http.Error(w, "Player name empty. Try again.", http.StatusBadRequest)
		return
	}

	player, ok := gamelogic.CreatePlayer(playerName)
	if !ok {
		http.Error(w, "Player with this name already exists.", http.StatusBadRequest)
		return
	}

	// Set the player ID in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  playerIdCookie,
		Value: player.Id,
		Path:  "/",
	})
	w.Write([]byte("Player " + playerName + " created."))
}

func PlayerReadyHandler(w http.ResponseWriter, r *http.Request) {
	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	playerId, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Player not identified. Make sure you have created one.", http.StatusBadRequest)
		return
	}

	gameId, err := r.Cookie(gameIdCookie)
	if err != nil {
		http.Error(w, "Could not find game id cookie.", http.StatusBadRequest)
		return
	}

	gamelogic.PlayerReady(gameId.Value, playerId.Value)
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering CreateGame handler")
	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("Cannot parse form.", "error", err)
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	password := r.FormValue("game-password")
	if password == "" {
		slog.Error("Empty password in create game request", "Form", r.PostForm)
		http.Error(w, "Cannot create game with empty password", http.StatusBadRequest)
		return
	}

	player, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Player not identified. Make sure you have created one.", http.StatusBadRequest)
		return
	}

	game, created := gamelogic.CreateGame(password)
	if !created {
		http.Error(w, "Could not create game. Probably a game with the same password is already running", http.StatusInternalServerError)
		return
	}

	_, err = gamelogic.JoinGame(password, player.Value)
	if err != nil {
		http.Error(w, "Could not join game.", http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/round-question")
	http.SetCookie(w, &http.Cookie{
		Name:  gameIdCookie,
		Value: game.Id,
		Path:  "/",
	})
	slog.Debug("Redirecting to /round-question")
	w.Write(nil)
	return
}

func JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Entering JoinGame handler")
	if !IsPost(r) {
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("Cannot parse form.", "error", err)
		http.Error(w, "Error. Check server logs.", http.StatusBadRequest)
		return
	}

	password := r.FormValue("game-password")
	if password == "" {
		slog.Error("Empty password in create game request", "Form", r.PostForm)
		http.Error(w, "Cannot create game with empty password", http.StatusBadRequest)
		return
	}

	playerId, err := r.Cookie(playerIdCookie)
	if err != nil {
		http.Error(w, "Player not identified. Make sure you have created one.", http.StatusBadRequest)
		return
	}

	game, err := gamelogic.JoinGame(password, playerId.Value)
	if err != nil {
		http.Error(w, "Could not join game. Check server logs", http.StatusBadRequest)
		return
	}

	w.Header().Set("HX-Redirect", "/round-question")
	http.SetCookie(w, &http.Cookie{
		Name:  gameIdCookie,
		Value: game.Id,
		Path:  "/",
	})

	slog.Debug("Redirecting to /round-question")
	w.Write(nil)
	return
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
