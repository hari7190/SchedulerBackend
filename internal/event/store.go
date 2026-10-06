package event

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
}

type EvenCreationRequest struct {
	Title string `json:"title"`
}

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(db *sql.DB) *MySQLStore {
	return &MySQLStore{db: db}
}

func (p *MySQLStore) Create(e EvenCreationRequest) error {
	query := `INSERT INTO events (id, title, created_on, modified_on) VALUES (?, ?, ?, ?)`
	_, err := p.db.Exec(query, uuid.NewString(), e.Title, time.Now(), time.Now())

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (p *MySQLStore) GetById(id string) (Event, error) {
	query := `SELECT id, title, created_on FROM events WHERE id = $1`
	row := p.db.QueryRow(query, id)

	var e Event
	err := row.Scan(&e.ID, e.Title, e.CreatedOn)
	if err != nil {
		return Event{}, err
	}
	return e, nil
}
