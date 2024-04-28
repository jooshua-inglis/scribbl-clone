package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"scribl-clone/data"

	"github.com/go-chi/chi"
)

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	playerId := chi.URLParam(r, "playerId")

	player, err := data.GetPlayer(playerId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if player == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	payload, err := json.Marshal(player)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Write(payload)
}
