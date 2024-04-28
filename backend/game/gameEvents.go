package game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"scribl-clone/data"
	"scribl-clone/eventListener"
	"scribl-clone/utils"
	"slices"
	"time"
)

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}
type Line struct {
	Points []Point `json:"points"`
	Size   int     `json:"size"`
	Rgb    [3]int  `json:"rgb"`
}

type GameEvent struct {
	EventType    int
	EventPayload any
}

type PlayerUpdatePayload struct {
	PlayerId string
	Updates  map[string]any
}

type GameUpdatePayload map[string]any
type ScoreUpdatePayload map[string]int

type DrawingEventPayload struct {
	Line  Line `json:"line"`
	Index int  `json:"index"`
}

const (
	GAME_EVENT_GUESS_OCCURRED = 0
	GAME_EVENT_SCORE_UPDATE   = 2
	GAME_EVENT_GAME_UPDATE    = 3
	GAME_EVENT_PLAYER_UPDATE  = 4
	GAME_EVENT_PLAYER_JOIN    = 5
	GAME_EVENT_DRAWING        = 6

	ROUND_END_REASON_TIMEOUT = "TIMER_RAN_OUT"
)

const (
	POINTS_DRAWER_TIMEOUT        = 10
	POINTS_DRAWER_CORRECT_GUESS  = 30
	POINTS_GUESSER_CORRECT_GUESS = 30
)

func GetGameChannelName(gameId string) string {
	return fmt.Sprintf("game/%s", gameId)
}

func StartRound(gameId string, drawer string) error {
	return UpdateGame(gameId, map[string]any{
		"state":               data.GAME_STATE_SELECTING_WORD,
		"turn":                drawer,
		"lastStateChangeTime": time.Now().UTC().Format(time.RFC3339),
	})
}

func SelectWord(gameId string, word string) error {
	return UpdateGame(gameId, map[string]any{
		"state":               data.GAME_STATE_DRAWING,
		"lastStateChangeTime": time.Now().UTC().Format(time.RFC3339),
	})
}

func EndTurn(gameId string, reason string) error {
	return UpdateGame(gameId, map[string]any{
		"state":               data.GAME_STATE_SELECTING_WORD,
		"lastStateChangeTime": time.Now().UTC().Format(time.RFC3339),
	})
}

func EndGame(gameId string) error {
	return UpdateGame(gameId, map[string]any{
		"state":               data.GAME_STATE_SELECTING_WORD,
		"lastStateChangeTime": time.Now().UTC().Format(time.RFC3339),
	})
}

func ScoreUpdate(gameId string, scores map[string]int) error {
	return publishEvent(gameId, GameEvent{
		EventType:    GAME_EVENT_SCORE_UPDATE,
		EventPayload: scores,
	})
}

func UpdateGame(gameId string, gameUpdates map[string]any) error {
	return publishEvent(gameId, GameEvent{
		EventType:    GAME_EVENT_GAME_UPDATE,
		EventPayload: gameUpdates,
	})
}

func UpdatePlayer(gameId string, playerId string, playerUpdates map[string]any) error {
	return publishEvent(gameId, GameEvent{
		EventType: GAME_EVENT_PLAYER_UPDATE,
		EventPayload: PlayerUpdatePayload{
			PlayerId: playerId,
			Updates:  playerUpdates,
		},
	})
}

func AddPlayer(gameId string, player *data.Player) error {
	return publishEvent(gameId, GameEvent{
		EventType:    GAME_EVENT_PLAYER_JOIN,
		EventPayload: player,
	})
}

func UpsertLine(gameId string, line Line, lineIndex int) error {
	slog.Info("hello", "there", lineIndex)
	return publishEvent(gameId, GameEvent{
		EventType: GAME_EVENT_DRAWING,
		EventPayload: DrawingEventPayload{
			Line:  line,
			Index: lineIndex,
		},
	})
}

func publishEvent(gameId string, event GameEvent) error {
	pubSub := eventListener.GetPubSub()
	data, err := json.Marshal(event)

	if err != nil {
		return err
	}

	err = pubSub.Publish(
		context.Background(),
		GetGameChannelName(gameId),
		string(data),
	).Err()

	if err != nil {
		slog.Error(err.Error())
	}

	return err
}

func getNextDrawer(gameId string, currentDrawer string) (playerId string, startNextRound bool, err error) {
	db := data.GetDb()

	players := []data.Player{}
	err = db.Select(&players, `SELECT id, date_created FROM player WHERE game = $1`, gameId)

	if err != nil {
		return "", false, err
	}
	slices.SortFunc(players, func(a, b data.Player) int {
		return a.DateCreated.Compare(b.DateCreated)
	})

	for i := range players {
		if players[i].Id == currentDrawer {
			if i == len(players)-1 {
				// No more players left,
				// either go to next round or end game
				return players[0].Id, true, nil
			}
			return players[i+1].Id, false, nil
		}
	}
	return "", false, errors.New("bad game state, current drawer is not in this game")

}

/*
To go to the next round, the following need to happen

- A new drawer is selected
- If all players have played this round go onto the next round
- If all rounds have ended, end the game. This means the state, round and turn will be updated.
*/
func GotoNextTurn(gameId string, currentDrawerId string) error {
	db := data.GetDb()
	nextDrawerId, startNextRound, err := getNextDrawer(gameId, currentDrawerId)
	if err != nil {
		return err
	}

	g := data.Game{}
	if err := db.Get(&g, `SELECT id, rounds, current_round FROM game WHERE id = $1`, gameId); err != nil {
		return err
	}

	newState := utils.If(g.Rounds == g.CurrentRound && startNextRound, data.GAME_STATE_WAITING_FOR_PLAYERS, data.GAME_STATE_SELECTING_WORD)
	var nextRound int
	if newState == data.GAME_STATE_WAITING_FOR_PLAYERS {
		nextRound = 1
	} else if startNextRound {
		nextRound = g.CurrentRound + 1
	} else {
		nextRound = g.CurrentRound
	}

	gameUpdates := map[string]any{
		"turn":         nextDrawerId,
		"state":        newState,
		"currentRound": nextRound,
	}
	UpdateGame(gameId, gameUpdates)
	if _, err = db.Exec(`UPDATE game SET turn = $2, state = $3, current_round = $4 WHERE id = $1;`, g.Id, nextDrawerId, newState, nextRound); err != nil {
		return err
	}

	_, err = db.Query(`UPDATE player SET guessed_correct = false WHERE game = $1`, g.Id)
	if err != nil {
		return err
	}

	players := []data.Player{}
	db.Select(&players, `SELECT id FROM player WHERE game = $1`, g.Id)

	for i := range players {
		log.Println(players[i].Id)
		UpdatePlayer(gameId, players[i].Id, map[string]any{"GuessedCorrect": false})
	}
	return nil
}
