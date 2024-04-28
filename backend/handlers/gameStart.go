package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"scribl-clone/data"
	"scribl-clone/game"
	"scribl-clone/utils"

	"github.com/go-chi/chi"
)

func StartGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	db := data.GetDb()
	var err error

	// ====== Fetch Required data ======
	g := data.Game{}
	err = db.QueryRowx(`SELECT id, state FROM game WHERE id = $1 LIMIT 1;`, gameId).StructScan(&g)
	if err == sql.ErrNoRows {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if g.State != data.GAME_STATE_WAITING_FOR_PLAYERS {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var playerId string
	err = db.QueryRow(`
		SELECT id 
		FROM player WHERE 
		date_created = (
			SELECT MIN(date_created)
			FROM player
			WHERE game = $1
		)
		LIMIT 1		
	`, gameId).Scan(&playerId)

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	_, err = db.Exec(`UPDATE game SET state = $2, turn = $3 WHERE id = $1;`, g.Id, data.GAME_STATE_SELECTING_WORD, playerId)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Bad Request", http.StatusInternalServerError)
		return
	}
	game.StartRound(gameId, playerId)

	w.Write(utils.STANDARD_SUCCESS_RESPONSE)
}
