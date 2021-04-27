package words

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func InsertNewWord(db *sql.DB, woord, lidwoord string) error {
	tx, err := db.Begin()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	insertWordStmt, err := tx.Prepare(`INSERT INTO woord(uid, content, created_at) VALUES (?, ?, ?)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	defer insertWordStmt.Close()

	insertWordRelationLidwoordStmt, err := tx.Prepare(`
		INSERT INTO woord_lidwoord(woord_id, lidwoord_id) VALUES (
		  (SELECT id FROM woord WHERE content = ? LIMIT 1), (SELECT id FROM lidwoord WHERE content = ? LIMIT 1))
		`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	defer insertWordRelationLidwoordStmt.Close()

	_, err = insertWordStmt.Exec(uuid.NewString(), woord, time.Now().Unix())
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = insertWordRelationLidwoordStmt.Exec(woord, lidwoord)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
