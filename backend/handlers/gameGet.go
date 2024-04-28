package handlers

import (
	"encoding/json"
	"net/http"
	"scribl-clone/data"

	"github.com/go-chi/chi"
)

func GetGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "gameId")

	game, err := data.GetGame(id)

	if game == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	payload, err := json.Marshal(game)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Write(payload)
}
