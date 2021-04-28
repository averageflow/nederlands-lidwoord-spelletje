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

type GameWordWithPlural struct {
	GameWord
	Plural    string
	PluralUID string
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

func GetRandomWordWithPlural(db *sql.DB) (*GameWordWithPlural, error) {
	row := db.QueryRow(`
		SELECT (SELECT woord.uid FROM woord WHERE id = woord_plural.singular_id LIMIT 1)     AS singular_uid,
			   (SELECT woord.content FROM woord WHERE id = woord_plural.singular_id LIMIT 1) AS singular_content,
			   (SELECT lidwoord.content
				FROM lidwoord
				WHERE id = (SELECT lidwoord_id
							FROM woord_lidwoord
							WHERE woord_id = woord_plural.singular_id
							LIMIT 1)
				LIMIT 1)                                                                     AS lidwoord_content,
			   (SELECT woord.uid FROM woord WHERE id = woord_plural.plural_id LIMIT 1)       AS plural_uid,
			   (SELECT woord.content FROM woord WHERE id = woord_plural.plural_id LIMIT 1)   AS plural_content
		FROM woord_plural
		ORDER BY RANDOM()
		LIMIT 1;
	`)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var gameWord GameWordWithPlural

	err := row.Scan(&gameWord.UID, &gameWord.Content, &gameWord.Lidwoord, &gameWord.PluralUID, &gameWord.Plural)
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

	// Insert singular word
	_, err = insertWordStmt.Exec(uuid.NewString(), woord, time.Now().Unix())
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

	if plural != "" {
		insertWordPluralStmt, err := tx.Prepare(`
		INSERT INTO woord_plural(singular_id, plural_id) VALUES (
		  (SELECT id FROM woord WHERE content = ? LIMIT 1), (SELECT id FROM woord WHERE content = ? LIMIT 1))
		`)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		defer insertWordPluralStmt.Close()

		// Insert plural word
		_, err = insertWordStmt.Exec(uuid.NewString(), plural, time.Now().Unix())
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
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
