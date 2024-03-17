package sql

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

// localhost:8080/user
type User struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

// localhost:8080/quest
type Quest struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

// localhost:8080/event
type Event struct {
	UserId  int `json:"user_id"`
	QuestId int `json:"quest_id"`
	//Users       []int   `json:"users_ids"`
	//Quests      []int   `json:"quest_ids"`
	UserBalance float64 `json:"user_balance"`
}

type EventResponse struct {
	UserId  int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Quests  []int   `json:"quests"`
}

func (u *User) UpdateDatabaseById(id int, db *sql.DB) error {
	if u.Balance > 0 {
		if _, err := db.Exec("UPDATE users SET balance = $1 WHERE id = $2", u.Balance, id); err != nil {
			return errors.New("error sql request")
		}
	}
	if u.Name != "" {
		if _, err := db.Exec("UPDATE users SET name = $1 WHERE id = $2", u.Name, id); err != nil {
			return errors.New("error sql request")
		}
	}

	return nil
}

func (u *User) UpdateDatabaseByName(name string, db *sql.DB) error {
	if u.Balance > 0 {
		responce, err := db.Exec("UPDATE users SET balance = $1 WHERE name = $2", u.Balance, name)

		if err != nil {
			return errors.New("error sql request")
		}
		num, err := responce.RowsAffected()
		if num == 0 {
			return errors.New("error user not exist")
		}
		if err != nil {
			return errors.New("error sql request")
		}
	}
	if u.Name != "" {
		responce, err := db.Exec("UPDATE users SET name = $1 WHERE name = $2", u.Name, name)

		if err != nil {
			return errors.New("error sql request")
		}
		num, err := responce.RowsAffected()
		if num == 0 {
			return errors.New("error user not exist")
		}
		if err != nil {
			return errors.New("error sql request")
		}
	}

	return nil
}

func (q *Quest) UpdateDatabaseById(id int, db *sql.DB) error {
	if q.Cost > 0 {
		if _, err := db.Exec("UPDATE quests SET cost = $1 WHERE id = $2", q.Cost, id); err != nil {
			return errors.New("error sql request")
		}
	}
	if q.Name != "" {
		if _, err := db.Exec("UPDATE quests SET name = $1 WHERE id = $2", q.Name, id); err != nil {
			return errors.New("error sql request")
		}
	}

	return nil
}

func (q *Quest) UpdateDatabaseByName(name string, db *sql.DB) error {
	if q.Cost > 0 {
		if _, err := db.Exec("UPDATE quests SET cost = $1 WHERE name = $2", q.Cost, name); err != nil {
			return errors.New("error sql request")
		}
	}
	if q.Name != "" {
		if _, err := db.Exec("UPDATE quests SET name = $1 WHERE name = $2", q.Name, name); err != nil {
			return errors.New("error sql request")
		}
	}

	return nil
}

func (q *Quest) AppendDatabase(db *sql.DB) error {
	response, err := db.Exec("SELECT name FROM quests WHERE name = $1", q.Name)

	if err != nil {
		return errors.New("error")
	}

	count, err := response.RowsAffected()

	if err != nil {
		return errors.New("error")
	}
	if count > 0 {
		return errors.New("quest already exist")
	}

	_, err = db.Exec("INSERT INTO quests (name, cost) VALUES ($1, $2)", q.Name, q.Cost)
	if err != nil {
		return errors.New("error")
	}

	err = db.QueryRow("SELECT id FROM quests WHERE name = $1", q.Name).Scan(&q.Id)

	if err != nil {
		return errors.New("error query")
	}
	return nil
}

func (u *User) AppendDatabaseById(db *sql.DB) error {
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

	err = db.QueryRow("SELECT id FROM users WHERE name = $1", u.Name).Scan(&u.Id)

	if err != nil {
		return errors.New("error query")
	}
	return nil
}

func (u *User) GetById(id int, db *sql.DB) error {
	err := db.QueryRow("SELECT id, name, balance FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Balance)
	if err != nil {
		return errors.New("error query")
	}
	return nil
}

func (u *User) GetByName(name string, db *sql.DB) error {
	err := db.QueryRow("SELECT id, name, balance FROM users WHERE name = $1", name).Scan(&u.Id, &u.Name, &u.Balance)
	if err != nil {
		return errors.New("error query")
	}
	return nil
}

func (q *Quest) GetById(id int, db *sql.DB) error {
	err := db.QueryRow("SELECT id, name, cost FROM quests WHERE id = $1", id).Scan(&q.Id, &q.Name, &q.Cost)
	if err != nil {
		return errors.New("error query")
	}
	return nil
}

func (q *Quest) GetByName(name string, db *sql.DB) error {
	err := db.QueryRow("SELECT id, name, cost FROM quests WHERE name = $1", name).Scan(&q.Id, &q.Name, &q.Cost)
	if err != nil {
		return errors.New("error query")
	}
	return nil
}

func (e *Event) AppendDatabase(db *sql.DB) error {
	response, _ := db.Exec("SELECT id FROM users WHERE id = $1", e.UserId)
	if num, _ := response.RowsAffected(); num == 0 {
		return errors.New("user is not exist")
	}

	response, _ = db.Exec("SELECT id FROM quests WHERE id = $1", e.QuestId)
	if num, _ := response.RowsAffected(); num == 0 {
		return errors.New("quest is not exist")
	}

	response, err := db.Exec("SELECT user_id, quest_id FROM user_quest WHERE user_id = $1 AND quest_id = $2", e.UserId, e.QuestId)

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

	response, err = db.Exec("INSERT INTO user_quest (user_id, quest_id) VALUES ($1, $2)", e.UserId, e.QuestId)

	if err != nil {
		return errors.New("error")
	}

	response, err = db.Exec("UPDATE users SET balance = balance + (SELECT cost FROM quests WHERE id = $1) WHERE id = $2", e.QuestId, e.UserId)

	if err != nil {
		return errors.New("error")
	}

	_ = db.QueryRow("SELECT balance FROM users WHERE id = $1", e.UserId).Scan(&e.UserBalance)

	return nil
}

func (e *Event) GetByUser(db *sql.DB) (*EventResponse, error) {
	var eventResp EventResponse
	rows, err := db.Query("SELECT quest_id FROM user_quest WHERE user_id = $1", e.UserId)
	if err != nil {
		return nil, errors.New("error sql request")
	}
	err = db.QueryRow("SELECT id, name, balance FROM users WHERE id = $1", e.UserId).Scan(&eventResp.UserId, &eventResp.Name, &eventResp.Balance)
	if err != nil {
		return nil, errors.New("error query")
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, errors.New("error parsing rows")
		}
		eventResp.Quests = append(eventResp.Quests, id)
	}
	return &eventResp, nil
}

func RemoveUserFromDatabaseById(id int, db *sql.DB) error {
	res, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return errors.New("error removing from database")
	}
	if val, _ := res.RowsAffected(); val == 0 {
		return errors.New("error user is not exist")
	}
	return nil
}

func RemoveUserFromDatabaseByName(name string, db *sql.DB) error {
	res, err := db.Exec("DELETE FROM users WHERE name = $1", name)
	if err != nil {
		return errors.New("error removing from database")
	}
	if val, _ := res.RowsAffected(); val == 0 {
		return errors.New("error user is not exist")
	}
	return nil
}

func RemoveQuestFromDatabaseById(id int, db *sql.DB) error {
	res, err := db.Exec("DELETE FROM quests WHERE id = $1", id)
	if err != nil {
		return errors.New("error removing from database")
	}
	if val, _ := res.RowsAffected(); val == 0 {
		return errors.New("error user is not exist")
	}
	return nil
}

func RemoveQuestFromDatabaseByName(name string, db *sql.DB) error {
	res, err := db.Exec("DELETE FROM quests WHERE name = $1", name)
	if err != nil {
		return errors.New("error removing from database")
	}
	if val, _ := res.RowsAffected(); val == 0 {
		return errors.New("error user is not exist")
	}
	return nil
}

func NewConn() *sql.DB { // need to write up some config file read
	connStr := "user=server password=server dbname=api_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
