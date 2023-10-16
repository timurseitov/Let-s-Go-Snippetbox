package postgreSQL

import (
	"GOLEARN/pkg/models"
	"database/sql"
	"errors"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 DAY')
		RETURNING id
		`

	var id int64

	err := m.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp AND id = $1
	`

	row := m.DB.QueryRow(stmt, id)

	snippet := &models.Snippet{}

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return snippet, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		snippet := &models.Snippet{}

		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
