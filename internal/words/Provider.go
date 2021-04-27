package words

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type GameWord struct {
	Content string
	UID     string
}

func GetRandomWord(db *sql.DB) (*GameWord, error) {
	row := db.QueryRow(`SELECT uid, content FROM woord LIMIT 1`)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var uid, content string

	err := row.Scan(&uid, &content)
	if err != nil {
		return nil, err
	}

	return &GameWord{UID: uid, Content: content}, nil
}
