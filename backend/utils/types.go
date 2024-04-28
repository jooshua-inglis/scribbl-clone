package utils

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	value, err := s.Value()
	if err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	var value any
	err := json.Unmarshal(data, value)
	if err != nil {
		return nil
	}
	s.String, s.Valid = value.(string)
	return nil
}

func CreateNullString(data any) NullString {
	s := NullString{}
	s.String, s.Valid = data.(string)
	return s
}
