package server

import (
	sqlpkg "VK_Internship_Go/sql"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}

	rt.ItemAppendDatabase(cont, &u)

}

func (rt *Rout) ItemAppendDatabase(cont *gin.Context, i Item) {
	if i.AppendDatabase(rt.db) != nil {
		cont.JSON(400, gin.H{"message": "error database indexing"})
	} else {
		cont.JSON(200, gin.H{"message": "ok"})
	}
}

func (rt *Rout) QuestPost(cont *gin.Context) {
	var q sqlpkg.Quest
	if err := cont.BindJSON(&q); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}

	rt.ItemAppendDatabase(cont, &q)

}
func (rt *Rout) EventPost(cont *gin.Context) {
	var e sqlpkg.Event
	if err := cont.BindJSON(&e); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}

	rt.ItemAppendDatabase(cont, &e)
}

func (rt *Rout) UserGetById(cont *gin.Context) {
	var u sqlpkg.User

	id, err := strconv.Atoi(cont.Param("id"))
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}

	if err = u.GetById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	cont.IndentedJSON(http.StatusOK, u)
}

func (rt *Rout) UserGetByName(cont *gin.Context) {
	var u sqlpkg.User

	name := cont.Param("name")
	if err := u.GetByName(name, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	cont.IndentedJSON(http.StatusOK, u)
}

func (rt *Rout) QuestGetById(cont *gin.Context) {
	var q sqlpkg.Quest
	id, err := strconv.Atoi(cont.Param("id"))
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}

	if err = q.GetById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	cont.IndentedJSON(http.StatusOK, q)
}

func (rt *Rout) QuestGetByName(cont *gin.Context) {
	var q sqlpkg.Quest
	name := cont.Param("name")

	if err := q.GetByName(name, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	cont.IndentedJSON(http.StatusOK, q)
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
	rout.router.GET("/users/id/:id", rout.UserGetById)
	rout.router.GET("/users/name/:name", rout.UserGetByName)
	rout.router.GET("/quests/id/:id", rout.QuestGetById)
	rout.router.GET("/quests/name/:name", rout.QuestGetByName)
	//rout.router.GET("/event", rout.EventGet)
	err := rout.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
