package server

import (
	sqlpkg "VK_Internship_Go/sql"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

// NEED TO USE MANNERS PACKAGE FOR SHUT DOWN SERVER
// http.Request.Method
// net/url
// "github.com/braintree/manners"

// localhost:8080/user
type User struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

//localhost:8080/quest
type Quest struct { // мб сделать из этого интерфейс
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type Item interface {
	AppendDatabase(db *sql.DB)
}

type Rout struct {
	router *gin.Engine
	db     *sql.DB
}

func (u *User) AppendDatabase(db *sql.DB) {
	response, err := db.Exec("SELECT name FROM users WHERE name = $1", u.Name)
	if err != nil {
		return
	}
	count, err := response.RowsAffected()
	if err != nil {
		return
	}
	if count > 0 {
		log.Println("User already exist")
		return
	}
	response, err = db.Exec("INSERT INTO users (name, balance) VALUES ($1, $2)", u.Name, u.Balance)
	if err != nil {
		return
	}
	fmt.Println(response)
}

func (q *Quest) AppendDatabase(db *sql.DB) {
	response, err := db.Exec("INSERT INTO quests (name, cost) VALUES ($1, $2)", q.Name, q.Cost)
	if err != nil {
		return
	}
	fmt.Println(response)
}

func (rt *Rout) UserPost(cont *gin.Context) {
	var u User
	if err := cont.BindJSON(&u); err != nil {
		return
	}
	u.AppendDatabase(rt.db)
}

func (rt *Rout) QuestPost(cont *gin.Context) {
	var u User
	if err := cont.BindJSON(&u); err != nil {
		return
	}

	u.AppendDatabase(rt.db)
}

func NewRout(g *gin.Engine, d *sql.DB) *Rout {
	return &Rout{
		router: g,
		db:     d,
	}
}

func HttpServer() {
	rout := NewRout(gin.Default(), sqlpkg.NewConn())

	rout.router.POST("/users", rout.UserPost)
	rout.router.POST("/quests", rout.QuestPost)

	err := rout.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
