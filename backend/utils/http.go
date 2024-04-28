package utils

import (
	"errors"
	"log/slog"
	"net/http"
	"scribl-clone/player"
)

var STANDARD_SUCCESS_RESPONSE = []byte(`{"message": "success"}`)

func HandleError(w http.ResponseWriter, err error) {
	slog.Debug(err.Error())
	if errors.Is(err, ErrResourceNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if errors.Is(err, ErrInvalidArguments) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if errors.Is(err, player.ErrUnauthorized) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	slog.Error(err.Error())
	http.Error(w, http.StatusText(500), 500)
}
