package sockets

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"scribl-clone/data"
	"scribl-clone/game"
	"scribl-clone/utils"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// TODO: remove this
	CheckOrigin: func(r *http.Request) bool { return true },
}

func InitGameConnection(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	db := data.GetDb()

	p := data.Player{}
	err := db.Get(&p, `SELECT game FROM player WHERE id = $1`, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		utils.HandleError(w, err)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	ws.WriteJSON(struct {
		Kind    string `json:"kind"`
		Message string `json:"message"`
	}{
		Kind:    "init-connection",
		Message: "successfully connected",
	})

	channelName := game.GetGameChannelName(p.Game)
	conn, err := CreateConnection(channelName, userId, ws)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	conn.SetCallback(func(s string) {
		gameEvent := game.GameEvent{}
		err := json.Unmarshal([]byte(s), &gameEvent)
		if err != nil {
			slog.Error("Error decoding websocket connection", "error", err)
		}
		if gameEvent.EventType == game.GAME_EVENT_DRAWING {
			drawingEvent := struct {
				game.GameEvent
				EventPayload game.DrawingEventPayload
			}{}
			err = json.Unmarshal([]byte(s), &drawingEvent)
			if err != nil {
				slog.Error("Error decoding drawing event", "error", err)
			}
			game.UpsertLine(p.Game, drawingEvent.EventPayload.Line, drawingEvent.EventPayload.Index)
		}
	})
}
