package gamelogic

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var games map[string]Game = make(map[string]Game)
var players map[string]Player = make(map[string]Player)

func CreateGame(password string) (Game, bool) {
	for _, v := range games {
		if v.Password == password && !v.IsComplete {
			return Game{}, false
		}
	}
	game := Game{
		Id:         uuid.New().String(),
		Players:    []Player{},
		Rounds:     []Round{},
		Password:   password,
		Started:    false,
		IsComplete: false,
	}

	// TODO: Remove dummy values
	dummyPlayer := Player{uuid.New().String(), "dummy player", 0}
	round := Round{uuid.New().String(), "q1", []Answer{
		{"asfad", "dummy answer", dummyPlayer, 0},
	}}
	game.Rounds = append(game.Rounds, round)

	games[game.Id] = game
	slog.Info("Created game", "game", game)
	return game, true
}

func JoinGame(password string, playerId string) (Game, error) {
	game := Game{} //gameExists := false
	for _, g := range games {
		if g.Password == password {
			game = g
			break
		}
	}

	if game.Id == "" {
		slog.Info("Game does not exist", "requested-password", password)
		return Game{}, errors.New("Game does not exist.")
	}

	player, playerExists := players[playerId]
	if !playerExists {
		slog.Info("Player does not exist", "requested-player", player)
		return Game{}, errors.New("Player does not exist.")
	}

	game.AddPlayerToGame(player)

	return game, nil
}

func GetLatestRound(gameId string) Round {
	game := games[gameId]
	slog.Debug("game found", "game", game)
	currentRound := game.Rounds[len(game.Rounds)-1]
	return currentRound
}

func (g *Game) AddPlayerToGame(player Player) {
	g.Players = append(g.Players, player)
	slog.Info("Player added to game", "player", player, "game", g)
}

func CreatePlayer(playerName string) (Player, bool) {
	playerId := uuid.New().String()
	player := Player{playerId, playerName, 0}
	players[playerId] = player
	slog.Info("Created player.", "player", player)
	return player, true
}

type Game struct {
	Id         string
	Players    []Player
	Rounds     []Round
	Password   string
	Started    bool
	IsComplete bool
}

type Player struct {
	Id     string
	Name   string
	Points int
}

type Round struct {
	Id       string
	Question string
	Answers  []Answer
}

type Answer struct {
	Id    string
	Text  string
	Owner Player
	Votes int
}
