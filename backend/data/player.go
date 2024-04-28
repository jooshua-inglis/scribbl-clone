package data

import (
	"log"
	"log/slog"
	"maps"
	"scribl-clone/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Id             string    `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Score          int       `db:"score" json:"score"`
	Game           string    `db:"game" json:"game"`
	DateCreated    time.Time `db:"date_created" json:"dateCreated"`
	GuessedCorrect bool      `db:"guessed_correct" json:"guessedCorrect"`
	ActiveState    string    `db:"active_state" json:"activeState"`
}

func CreatePlayer(name string, game string) (string, error) {
	db := GetDb()
	player := Player{
		Id:          uuid.NewString(),
		Name:        name,
		Score:       0,
		Game:        game,
		DateCreated: time.Now(),
		ActiveState: "creating",
	}
	slog.Info(name)

	_, err := db.NamedQuery(`INSERT INTO player
		(id, name, score, game, date_created, active_state)
		VALUES
		(:id, :name, :score, :game, :date_created, :active_state)
		`,
		player,
	)
	return player.Id, err
}

func GetPlayer(id string) (*Player, error) {
	db := GetDb()

	query := db.QueryRowx(`SELECT * FROM player WHERE id=$1`, id)
	if query.Err() != nil {
		return nil, query.Err()
	}

	player := Player{}
	err := query.StructScan(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func toSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

// TODO: SCRIBBL-1
func UpdatePlayer(id string, updateSet map[string]any) (*Player, error) {
	setQuery := dynamicUpdateSet([]string{"Name", "ActiveState"}, updateSet)
	if setQuery == "" {
		return nil, utils.ErrInvalidArguments
	}
	query := `UPDATE player SET ` + setQuery + " WHERE id = :Id RETURNING *"

	inputData := map[string]any{"Id": id}
	maps.Copy(inputData, updateSet)

	db := GetDb()
	p := Player{}
	row, err := db.NamedQuery(query, inputData)
	log.Println(query)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, utils.ErrResourceNotFound
	}

	if err = row.StructScan(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
