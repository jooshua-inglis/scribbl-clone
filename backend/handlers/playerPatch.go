package handlers

import (
	"encoding/json"
	"net/http"
	"scribl-clone/data"
	"scribl-clone/game"
	"scribl-clone/player"
	"scribl-clone/utils"
)

func PatchPlayer(w http.ResponseWriter, r *http.Request) {
	claim, err := player.AuthorizeRequest(r)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	playerId := claim.PlayerId

	updateSet := make(map[string]any)
	if err = json.NewDecoder(r.Body).Decode(&updateSet); err != nil {
		utils.HandleError(w, err)
		return
	}

	updatedPlayer, err := data.UpdatePlayer(playerId, updateSet)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	game.UpdatePlayer(claim.GameId, playerId, updateSet)

	payload, err := json.Marshal(updatedPlayer)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	w.Write(payload)
}
