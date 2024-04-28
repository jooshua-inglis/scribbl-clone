package handlers

import (
	"net/http"
	"scribl-clone/eventListener"
	"scribl-clone/game"

	"github.com/go-chi/chi"
)

func DummyEvent(w http.ResponseWriter, r *http.Request) {
	rdb := eventListener.GetPubSub()
	ctx := r.Context()
	gameId := chi.URLParam(r, "gameId")

	rdb.Publish(ctx, game.GetGameChannelName(gameId), "Hello there")
}
