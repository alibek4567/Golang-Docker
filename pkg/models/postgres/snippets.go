package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"time"

	"se02.com/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES ($1, $2, $3, $4) RETURNING id"

	var id int
	days, err := strconv.Atoi(expires)
	if err != nil {
		return 0, err
	}
	result := m.DB.QueryRow(context.Background(), stmt, title, content, time.Now(), time.Now().AddDate(0, 0, days))
	err = result.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > now() AND id = $1`

	row := m.DB.QueryRow(context.Background(), stmt, id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil

}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > $1 ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(context.Background(), stmt, time.Now())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {

		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
