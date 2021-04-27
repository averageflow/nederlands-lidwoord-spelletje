package words

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type GameWord struct {
	Content  string
	UID      string
	Lidwoord string
}

func GetRandomWord(db *sql.DB) (*GameWord, error) {
	row := db.QueryRow(`
		SELECT woord.uid, woord.content, lidwoord.content
		FROM woord
				 INNER JOIN woord_lidwoord on woord.id = woord_lidwoord.woord_id
				 INNER JOIN lidwoord on lidwoord.id = woord_lidwoord.lidwoord_id
		ORDER BY RANDOM()
		LIMIT 1;
	`)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var gameWord GameWord

	err := row.Scan(&gameWord.UID, &gameWord.Content, &gameWord.Lidwoord)
	if err != nil {
		return nil, err
	}

	return &gameWord, nil
}

func InsertNewWord(db *sql.DB, woord, lidwoord, plural string) error {
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

	insertWordPluralStmt, err := tx.Prepare(`
		INSERT INTO woord_plural(singular_id, plural_id) VALUES (
		  (SELECT id FROM woord WHERE content = ? LIMIT 1), (SELECT id FROM woord WHERE content = ? LIMIT 1))
		`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	defer insertWordPluralStmt.Close()

	// Insert singular word
	_, err = insertWordStmt.Exec(uuid.NewString(), woord, time.Now().Unix())
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert plural word
	_, err = insertWordStmt.Exec(uuid.NewString(), plural, time.Now().Unix())
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert lidwoord for singular
	_, err = insertWordRelationLidwoordStmt.Exec(woord, lidwoord)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert lidwoord (always de) for plural
	_, err = insertWordRelationLidwoordStmt.Exec(plural, "de")
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert plural relationship
	_, err = insertWordPluralStmt.Exec(woord, plural)
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
