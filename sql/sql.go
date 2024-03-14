package sql

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

// localhost:8080/user
type User struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

// localhost:8080/quest
type Quest struct { // мб сделать из этого интерфейс
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

// localhost:8080/event
type Event struct {
	User_id  int `json:"user_id"`
	Quest_id int `json:"quest_id"`
}

func (q *Quest) AppendDatabase(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO quests (name, cost) VALUES ($1, $2)", q.Name, q.Cost)
	if err != nil {
		return errors.New("error")
	}
	return nil
}

func (u *User) AppendDatabase(db *sql.DB) error {
	response, err := db.Exec("SELECT name FROM users WHERE name = $1", u.Name)

	if err != nil {
		return errors.New("error")
	}

	count, err := response.RowsAffected()

	if err != nil {
		return errors.New("error")
	}
	if count > 0 {
		return errors.New("user already exist")
	}

	response, err = db.Exec("INSERT INTO users (name, balance) VALUES ($1, $2)", u.Name, u.Balance)

	if err != nil {
		return errors.New("error")
	}

	return nil
}

func (e *Event) AppendDatabase(db *sql.DB) error {
	response, err := db.Exec("SELECT user_id, quest_id FROM user_quest WHERE user_id = $1 AND quest_id = $2", e.User_id, e.Quest_id)

	if err != nil {
		return errors.New("error")
	}

	count, err := response.RowsAffected()
	if err != nil {
		return errors.New("error")
	}
	if count > 0 {
		return errors.New("event already completed")
	}

	response, err = db.Exec("INSERT INTO user_quest (user_id, quest_id) VALUES ($1, $2)", e.User_id, e.Quest_id)

	if err != nil {
		return errors.New("error")
	}

	response, err = db.Exec("UPDATE users SET balance = balance + (SELECT cost FROM quests WHERE id = $1) WHERE id = $2", e.Quest_id, e.User_id)

	if err != nil {
		return errors.New("error")
	}
	return nil
}

func NewConn() *sql.DB { // need to write up some config file read
	connStr := "user=server password=server dbname=userdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
