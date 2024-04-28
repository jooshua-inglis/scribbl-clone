package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"scribl-clone/data"
	"scribl-clone/eventListener"
	"scribl-clone/game"
	"scribl-clone/player"
	"scribl-clone/utils"
	"time"

	"github.com/google/uuid"
)

func SelectWord(w http.ResponseWriter, r *http.Request) {
	var err error
	userClaim, err := player.AuthorizeRequest(r)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	body := struct {
		Word string `json:"word"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	gameId := userClaim.GameId
	drawerId := userClaim.PlayerId

	db := data.GetDb()

	g := data.Game{}
	if err := db.Get(&g, `SELECT id, turn, state FROM game WHERE id = $1`, gameId); err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if !g.Turn.Valid || g.Turn.String != drawerId {
		http.Error(w, "It's not your turn", 400)
		return
	}

	if g.State != data.GAME_STATE_SELECTING_WORD {
		http.Error(w, "It's not time to select a word yet", 400)
		return
	}

	if _, err = db.Exec(`UPDATE game SET word = $2, state = $3 WHERE id = $1;`, g.Id, body.Word, data.GAME_STATE_DRAWING); err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	game.SelectWord(gameId, body.Word)

	subscriptionId := uuid.NewString()
	timeoutHandler := time.AfterFunc(30*time.Second, func() {
		eventListener.Unsubscribe(game.GetGameChannelName(gameId), subscriptionId)
		if err = db.Get(&g, `SELECT state FROM game WHERE id = $1;`, g.Id); err != nil {
			slog.Error(err.Error())
			return
		}
		if g.State != data.GAME_STATE_DRAWING {
			return
		}
		if _, err = db.Exec(`UPDATE game SET state = $2 WHERE id = $1;`, g.Id, data.GAME_STATE_SELECTING_WORD); err != nil {
			slog.Error(err.Error())
			return
		}

		var drawerScore int
		err := db.Get(&drawerScore, `
			UPDATE player SET score = score + $1
				WHERE id = $2
				RETURNING score
			`,
			game.POINTS_DRAWER_TIMEOUT,
			drawerId,
		)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		game.ScoreUpdate(gameId, map[string]int{drawerId: drawerScore})
		if err := game.GotoNextTurn(gameId, drawerId); err != nil {
			slog.Error(err.Error())
		}
	})
	eventListener.Subscribe(game.GetGameChannelName(gameId), subscriptionId, func(d string) {
		gameEvent := game.GameEvent{}
		json.Unmarshal([]byte(d), &gameEvent)
		if gameEvent.EventType == game.GAME_EVENT_GAME_UPDATE {
			gameUpdateEvent := struct {
				game.GameEvent
				DrawingEventPayload game.GameUpdatePayload
			}{}
			if payload, ok := gameUpdateEvent.EventPayload.(map[string]any); ok {
				if state, ok := payload["state"]; ok && state != data.GAME_STATE_DRAWING {
					// Change of state has occurred, stop the timer.
					slog.Debug("change of state has occurred, stop the timer.")
					timeoutHandler.Stop()
				}
			}
		}
	})

	w.Write([]byte(`{"message": "success"}`))
}
