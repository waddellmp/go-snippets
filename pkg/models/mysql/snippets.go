package mysql

import (
	"database/sql"

	"github.com/waddellmp/go-snippets/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for the placeholder parameters. This
	// method returns a sql.Result object, which contains some basic
	// information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// Construct a snippet struct to return
	s := &models.Snippet{}

	// Query snippets table with id & map to Snippet struct
	// Handle no records found and error
	err := m.DB.QueryRow(
		`SELECT id, title, content, created, expires FROM snippets
	 WHERE expires > UTC_TIMESTAMP() AND id = ?`,
		id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// If everything went OK then return the snippet object
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	rows, err := m.DB.Query(`
	SELECT id, title, content, created, expires 
	FROM snippets
	WHERE expires > UTC_TIMESTAMP() 
	ORDER BY DESC LIMIT 10
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	// Iterate through the rows in the resultset. This
	// prepares the first (and then subsequent) row to be acted on by the
	// rows.Scan() method.
	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// Handle if there were issues reading the rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
