package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"scribl-clone/data"

	"github.com/go-chi/chi"
)

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	db := data.GetDb()

	players := []data.Player{}
	if err := db.Select(&players, `SELECT * FROM player WHERE game = $1`, gameId); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error(err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}

	payload, err := json.Marshal(players)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Write(payload)
}
