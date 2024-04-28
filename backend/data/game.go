package data

import (
	"errors"
	"scribl-clone/utils"
	"time"

	"github.com/google/uuid"
)

const (
	GAME_STATE_WAITING_FOR_PLAYERS = 0
	GAME_STATE_END                 = 1
	GAME_STATE_DRAWING             = 2
	GAME_STATE_SELECTING_WORD      = 3
)

type Game struct {
	Id                  string           `json:"id"`
	Word                string           `json:"word"`
	Rounds              int              `json:"rounds"`
	CurrentRound        int              `db:"current_round" json:"currentRound"`
	Turn                utils.NullString `json:"turn"`
	MaxPlayers          int              `db:"max_players" json:"maxPlayers"`
	State               int              `json:"state"`
	LastStateChangeTime time.Time        `db:"last_state_change_time" json:"lastStateChangeTime"`
	DateCreated         time.Time        `db:"date_created" json:"dateCreated"`
}

func GetGame(id string) (*Game, error) {
	if uuid.Validate(id) != nil {
		return nil, errors.New("invalid id")
	}
	db := GetDb()
	data := db.QueryRowx(`SELECT * FROM game WHERE id = $1`, id)
	if data.Err() != nil {
		return nil, data.Err()
	}
	game := Game{}
	err := data.StructScan(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CreateNewGame() (*Game, error) {
	db := GetDb()
	game := Game{
		Id:           uuid.NewString(),
		Word:         "",
		CurrentRound: 1,
		Turn:         utils.CreateNullString(nil),
		MaxPlayers:   10,
		State:        GAME_STATE_WAITING_FOR_PLAYERS,
		DateCreated:  time.Now(),
	}

	result := db.QueryRow(`INSERT INTO game
		(id, word, current_round, turn, max_players, state, date_created)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`,
		game.Id,
		game.Word,
		game.CurrentRound,
		game.Turn,
		game.MaxPlayers,
		game.State,
		game.DateCreated,
	)

	if result.Err() != nil {
		return nil, result.Err()
	}
	return &game, nil
}
