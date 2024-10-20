package gamelogic

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

var games map[string]Game = make(map[string]Game)
var players map[string]Player = make(map[string]Player)

func CreateGame(password string, playerId string) (Game, bool) {
	for _, v := range games {
		if v.Password == password && !v.IsComplete {
			return Game{}, false
		}
	}

	player := GetPlayer(playerId)

	game := Game{
		Id:         uuid.New().String(),
		Players:    []Player{player},
		Rounds:     []Round{},
		Password:   password,
		Started:    false,
		IsComplete: false,
		Score:      make(map[string]int),
	}

	games[game.Id] = game
	CreateNewRound(game.Id)
	slog.Info("Created game", "game", game)
	return game, true
}

func JoinGame(password string, playerId string) (Game, error) {
	var game *Game
	for _, g := range games {
		if g.Password == password {
			game = &g
			break
		}
	}

	if game == nil {
		slog.Error("Game does not exist", "requested-password", password)
		return Game{}, errors.New("Game does not exist.")
	}

	player, playerExists := players[playerId]
	if !playerExists {
		slog.Error("Player does not exist", "requested-player", player)
		return Game{}, errors.New("Player does not exist.")
	}

	AddPlayerToGame(game.Id, player)

	return *game, nil
}

func CreateNewRound(gameId string) {
	game := games[gameId]
	round := Round{}
	round.Id = uuid.New().String()
	round.Question = GetRandomQuestion(game.GetNextPlayerName())
	round.Answers = []Answer{}
	game.Rounds = append(game.Rounds, round)
	games[gameId] = game
	slog.Debug("Created new round", "game", game)
}

func AllPlayerAnswered(gameId string, roundId string) bool {
	game := games[gameId]

	for _, r := range game.Rounds {
		slog.Debug("Checking all players answered", "round", r, "players", game.Players)
		if r.Id == roundId {
			if len(r.Answers) == len(game.Players) {
				return true
			}
		}
	}
	return false
}

func AllPlayersSelectedChoice(gameId string, roundId string) bool {
	game := games[gameId]

	for _, r := range game.Rounds {
		slog.Debug("Checking all players selected choice", "round", r, "players", game.Players)
		if r.Id == roundId {
			if r.ChoiceCount == len(game.Players) {
				return true
			}
		}
	}
	return false
}

func AllPlayersReady(gameId string) bool {
	game := games[gameId]
	slog.Debug("Checking all players ready")
	for _, p := range game.Players {
		if !p.PlayerReady {
			return false
		}
	}
	return true
}

var playerReadyLock sync.Mutex

func PlayerReady(gameId string, playerId string) {
	game := games[gameId]
	playerReadyLock.Lock()
	allPlayersReady := true
	for i := range game.Players {
		p := &game.Players[i]
		if p.Id == playerId {
			p.PlayerReady = true
			slog.Info("Player is ready", "player", p)
		}
		if !p.PlayerReady {
			allPlayersReady = false
		}
	}

	// if all players are ready start a  new round
	if allPlayersReady {
		slog.Debug("All players ready", "players", game.Players)
		CreateNewRound(gameId)
	} else {
		slog.Debug("Not all players ready", "players", game.Players)
	}
	playerReadyLock.Unlock()
}

func GetLatestRound(gameId string) (Round, error) {
	game := games[gameId]
	if len(game.Rounds) == 0 {
		logMessage := "Error when trying to get latest round. Game has no rounds yet."
		slog.Error(logMessage)
		return Round{}, errors.New(logMessage)
	}
	currentRound := game.Rounds[len(game.Rounds)-1]
	slog.Debug("Getting latest round", "game", game, "roundId", currentRound.Id)
	return currentRound, nil
}

func AddAnswer(gameId string, playerId string, roundId string, answerText string) error {
	game, ok := games[gameId]
	if !ok {
		return errors.New("Game " + gameId + " does not exist")
	}

	player, ok := players[playerId]
	if !ok {
		return errors.New("Player " + playerId + " does not exist")
	}

	for i := range game.Rounds {
		r := &game.Rounds[i]
		if r.Id != roundId {
			continue
		}
		answer := Answer{
			Id:     uuid.New().String(),
			Text:   answerText,
			Owner:  player,
			Voters: []Player{}}

		// If player answer exists, overwrite it
		updatedAnswer := false
		for j, a := range r.Answers {
			if a.Owner.Id == playerId {
				r.Answers[j] = answer
				updatedAnswer = true
			}
		}

		if !updatedAnswer {
			game.Rounds[i].Answers = append(game.Rounds[i].Answers, answer)
		}

		slog.Debug("Adding answer", "game", game, "player", player, "roundId", r.Id, "answer", answer)
		return nil
	}
	return errors.New("Could not add answer")
}

func AddChoice(gameId string, playerId string, roundId string, choiceId string) error {
	game, ok := games[gameId]
	if !ok {
		return errors.New("Game " + gameId + " does not exist")
	}

	player, ok := players[playerId]
	if !ok {
		return errors.New("Player " + playerId + " does not exist")
	}

	for i, r := range game.Rounds {
		if r.Id == roundId {
			for j := range r.Answers {
				a := &r.Answers[j]
				if a.Id == choiceId {
					a.Voters = append(a.Voters, player)
					game.Score[a.Owner.Id] += 1
					game.Rounds[i].ChoiceCount++
					slog.Debug("Added choice", "game", game, "player", player, "roundId", r.Id, "answer", a)
					slog.Info("Score update", "score", game.Score)
					// Setting player ready in order to be able to check when starting next round
					player.PlayerReady = false
					players[playerId] = player
					return nil
				}
			}
		}
	}
	return errors.New("Could not add choice")

}

func GetScore(gameId string) map[string]int {
	game := games[gameId]
	return game.Score
}

func GetPlayer(playerId string) Player {
	return players[playerId]
}

func AddPlayerToGame(gameId string, player Player) {
	// Adding a player copy so the variables are not carried over to different games
	playerCopy := player
	game := games[gameId]
	game.Players = append(game.Players, playerCopy)
	games[gameId] = game
	slog.Info("Player added to game", "player", playerCopy, "game", game)
}

func CreatePlayer(playerName string) (Player, bool) {
	playerId := uuid.New().String()
	player := Player{playerId, playerName, false}
	players[playerId] = player
	slog.Info("Created player.", "player", player)
	return player, true
}

type Game struct {
	Id              string
	Players         []Player
	Rounds          []Round
	Password        string
	Started         bool
	IsComplete      bool
	Score           map[string]int // map[playerId]points
	NextPlayerIndex int
}

func (g *Game) GetNextPlayerName() string {
	slog.Debug("getting next player name", "index", g.NextPlayerIndex, "players", len(g.Players),
		"modulo", g.NextPlayerIndex%len(g.Players))
	index := g.NextPlayerIndex % len(g.Players)
	g.NextPlayerIndex++
	return g.Players[index].Name
}

type Player struct {
	Id          string
	Name        string
	PlayerReady bool
}

type Round struct {
	Id          string
	Question    string
	Answers     []Answer
	ChoiceCount int
}

type Answer struct {
	Id     string
	Text   string
	Owner  Player
	Voters []Player
}
