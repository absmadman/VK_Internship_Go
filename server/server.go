package server

import (
	sqlpkg "VK_Internship_Go/sql"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

// NEED TO USE MANNERS PACKAGE FOR SHUT DOWN SERVER
// http.Request.Method
// net/url
// "github.com/braintree/manners"

type Item interface {
	AppendDatabase(db *sql.DB) error
}

type Rout struct {
	router *gin.Engine
	db     *sql.DB
}

func (rt *Rout) UserPost(cont *gin.Context) {
	var u sqlpkg.User
	if err := cont.BindJSON(&u); err != nil {
		cont.JSON(400, "")
		return
	}

	if u.AppendDatabase(rt.db) != nil {
		cont.JSON(400, "")
	} else {
		cont.JSON(200, "")
	}
}

func (rt *Rout) QuestPost(cont *gin.Context) {
	var q sqlpkg.Quest
	if err := cont.BindJSON(&q); err != nil {
		cont.JSON(400, "")
		return
	}

	if q.AppendDatabase(rt.db) != nil {
		cont.JSON(400, "")
	} else {
		cont.JSON(200, "")
	}
}

func (rt *Rout) EventPost(cont *gin.Context) {
	var e sqlpkg.Event
	if err := cont.BindJSON(&e); err != nil {
		cont.JSON(200, "")
		return
	}

	if e.AppendDatabase(rt.db) != nil {
		cont.JSON(400, "")
	} else {
		cont.JSON(200, "")
	}
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
	rout.router.POST("/event", rout.EventPost)

	err := rout.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
