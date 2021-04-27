package words

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetRandomWord(db *sql.DB) string {
	row := db.QueryRow(`SELECT uid, content FROM woord LIMIT 1`)
	if row.Err() != nil {
		return row.Err().Error()
	}

	var uid, content string

	err := row.Scan(&uid, &content)
	if err != nil {
		return err.Error()
	}

	return content
}
