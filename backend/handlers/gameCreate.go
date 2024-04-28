package handlers

import (
	"encoding/json"
	"net/http"
	"scribl-clone/data"
	"scribl-clone/utils"
)

const ID_LENGTH = 6

func CreateGame(w http.ResponseWriter, r *http.Request) {
	game, err := data.CreateNewGame()
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	payload, err := json.Marshal(game)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	w.Write(payload)
}
