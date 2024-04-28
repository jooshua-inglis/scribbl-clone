package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"scribl-clone/data"
	"scribl-clone/game"
	"scribl-clone/player"
	"scribl-clone/utils"
	"slices"
)

func MakeGuess(w http.ResponseWriter, r *http.Request) {
	db := data.GetDb()
	var err error
	userClaim, err := player.AuthorizeRequest(r)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	body := struct {
		Guess string `json:"guess"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	gameId := userClaim.GameId
	guesserId := userClaim.PlayerId

	// Find game
	g := data.Game{}
	if err := db.Get(&g, `SELECT id, word, state, turn FROM game WHERE id = $1`, gameId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Could not find game", "gameId", gameId)
			http.Error(w, http.StatusText(404), 404)
		} else {
			utils.HandleError(w, err)
		}
		return
	}

	if g.State != data.GAME_STATE_DRAWING {
		slog.Debug("In wrong state")
		http.Error(w, "In wrong state", 400)
		return
	}

	if g.Turn.String == guesserId {
		slog.Debug("Requester is the drawer")
		http.Error(w, "Your the drawer! You can't guess", 400)
		return
	}

	if g.Word != body.Guess {
		slog.Debug("Incorrect guess")
		w.Write(utils.STANDARD_SUCCESS_RESPONSE)
		return
	}

	// Get players of the game
	players := []data.Player{}
	if err := db.Select(&players, `SELECT id, guessed_correct FROM player WHERE game = $1`, gameId); err != nil {
		utils.HandleError(w, err)
		return
	}

	gotoNextRound, err := handleCorrectGuess(gameId, guesserId, g.Turn.String, players)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	if gotoNextRound {
		if err := game.GotoNextTurn(gameId, g.Turn.String); err != nil {
			utils.HandleError(w, err)
			return
		}
	}

	w.Write(utils.STANDARD_SUCCESS_RESPONSE)
}

func handleCorrectGuess(gameId string, playerId string, drawer string, players []data.Player) (bool, error) {
	db := data.GetDb()

	guessingPlayerIndex := slices.IndexFunc(players, func(p data.Player) bool { return p.Id == playerId })
	if guessingPlayerIndex == -1 {
		return false, fmt.Errorf("cannot find player with id %s", playerId)
	}

	if players[guessingPlayerIndex].GuessedCorrect {
		return false, nil
	}

	scoreIncrease := 10
	players[guessingPlayerIndex].GuessedCorrect = true

	var drawerScore int
	err := db.Get(
		&drawerScore,
		`UPDATE player SET 
			guessed_correct = true, score = score + $2 
			WHERE id = $1 
		RETURNING score`,
		playerId,
		scoreIncrease)

	if err != nil {
		return false, err
	}
	// TODO
	game.UpdatePlayer(gameId, playerId, map[string]any{"GuessedCorrect": true})
	game.ScoreUpdate(gameId, map[string]int{playerId: drawerScore})

	return shouldGotoNextRound(players, drawer), nil
}

func shouldGotoNextRound(players []data.Player, drawer string) bool {
	for i := range players {
		if !players[i].GuessedCorrect && players[i].Id != drawer {
			return false
		}
	}
	return true
}
