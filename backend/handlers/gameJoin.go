package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"scribl-clone/data"
	"scribl-clone/game"
	"scribl-clone/player"

	"github.com/go-chi/chi"
)

type bodySchema struct {
	Name string `json:"name"`
}

type returnSchema struct {
	Player *data.Player `json:"player"`
	Token  string       `json:"token"`
}

func JoinGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	body := bodySchema{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// ====== Fetch Required data ======
	g, err := data.GetGame(gameId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if g == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// TODO: add in limit to the number of players that can join

	// ======= Create Player ========
	playerId, err := data.CreatePlayer(body.Name, gameId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	p, err := data.GetPlayer(playerId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	game.AddPlayer(gameId, p)

	payload, err := json.Marshal(returnSchema{
		Token: player.GenerateToken(player.PlayerClaim{
			PlayerId: playerId,
			GameId:   gameId,
		}),
		Player: p,
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Write(payload)
}
