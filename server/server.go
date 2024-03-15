package server

import (
	sqlpkg "VK_Internship_Go/sql"
	"database/sql"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
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
	UpdateDatabaseById(id int, db *sql.DB) error
	GetById(id int, db *sql.DB) error
	GetByName(name string, db *sql.DB) error
}

type Rout struct {
	router *gin.Engine
	db     *sql.DB
	uCache *lru.Cache[int, any] // cache for users
	qCache *lru.Cache[int, any] // cache for quests
}

func (rt *Rout) UserPost(cont *gin.Context) {
	var u sqlpkg.User
	if err := cont.BindJSON(&u); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	if u.Balance < 0 {
		cont.JSON(400, gin.H{"message": "error balance cannot be below zero"})
		return
	}
	if u.Name == "" {
		cont.JSON(400, gin.H{"message": "error name cannot be empty"})
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

func (rt *Rout) ItemReadDatabaseById(cont *gin.Context, i Item, id int) {
	if err := i.GetById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	cont.IndentedJSON(http.StatusOK, i)
}

func (rt *Rout) ItemReadDatabaseByName(cont *gin.Context, i Item, name string) {
	if err := i.GetByName(name, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	cont.IndentedJSON(http.StatusOK, i)
}

func (rt *Rout) QuestPost(cont *gin.Context) {
	var q sqlpkg.Quest
	if err := cont.BindJSON(&q); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	if q.Cost < 0 {
		cont.JSON(400, gin.H{"message": "error cost cannot be below zero"})
		return
	}
	if q.Name == "" {
		cont.JSON(400, gin.H{"message": "error name cannot be empty"})
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
	if e.AppendDatabase(rt.db) != nil {
		cont.JSON(400, gin.H{"message": "error database indexing"})
	} else {
		cont.JSON(200, gin.H{"message": "ok"})
	}
}

func (rt *Rout) UserGetById(cont *gin.Context) {
	var u sqlpkg.User
	id, err := strconv.Atoi(cont.Param("id"))
	/*
		if val, ok := rt.uCache.Get(id); ok {
			fmt.Println("RETURN CACHED")
			cont.IndentedJSON(http.StatusOK, val)
			return
		}

	*/
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	rt.ItemReadDatabaseById(cont, &u, id)
	//rt.uCache.Add(u.Id, u)
}

func (rt *Rout) UserGetByName(cont *gin.Context) {
	var u sqlpkg.User
	name := cont.Param("name")
	rt.ItemReadDatabaseByName(cont, &u, name)
}

func (rt *Rout) QuestGetById(cont *gin.Context) {
	var q sqlpkg.Quest
	id, err := strconv.Atoi(cont.Param("id"))
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	rt.ItemReadDatabaseById(cont, &q, id)
}

func (rt *Rout) QuestGetByName(cont *gin.Context) {
	var q sqlpkg.Quest
	name := cont.Param("name")
	rt.ItemReadDatabaseByName(cont, &q, name)
}

func NewRout(g *gin.Engine, d *sql.DB) *Rout {
	return &Rout{
		router: g,
		db:     d,
	}
}

func (rt *Rout) UserPutById(cont *gin.Context) {
	var u sqlpkg.User
	u.Balance = -1
	u.Name = ""
	if err := cont.BindJSON(&u); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	id, err := strconv.Atoi(cont.Param("id"))
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	if err = u.UpdateDatabaseById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error sql execution"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})

}

func (rt *Rout) UserPutByName(cont *gin.Context) {
	var u sqlpkg.User
	u.Balance = -1
	u.Name = ""
	if err := cont.BindJSON(&u); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	name := cont.Param("name")
	if err := u.UpdateDatabaseByName(name, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error sql execution"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})

}

func (rt *Rout) ItemUpdateDatabaseById(cont *gin.Context, id int, i Item) {
	if err := i.UpdateDatabaseById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error sql execution"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})
}

func (rt *Rout) QuestPutById(cont *gin.Context) {
	var q sqlpkg.Quest
	q.Cost = -1
	q.Name = ""
	if err := cont.BindJSON(&q); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	id, err := strconv.Atoi(cont.Param("id"))
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	if err = q.UpdateDatabaseById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error sql execution"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})
}

func (rt *Rout) QuestPutByName(cont *gin.Context) {
	var q sqlpkg.Quest
	q.Cost = -1
	q.Name = ""
	if err := cont.BindJSON(&q); err != nil {
		cont.JSON(400, gin.H{"message": "error binding json"})
		return
	}
	name := cont.Param("name")
	if err := q.UpdateDatabaseByName(name, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error sql execution"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})
}

func (rt *Rout) UserDeleteById(cont *gin.Context) {
	id, err := strconv.Atoi(cont.Param("id"))
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	if err = sqlpkg.RemoveUserFromDatabaseById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})
}

func (rt *Rout) UserDeleteByName(cont *gin.Context) {
	name := cont.Param("name")
	if err := sqlpkg.RemoveUserFromDatabaseByName(name, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})
}

func (rt *Rout) QuestDeleteById(cont *gin.Context) {
	id, err := strconv.Atoi(cont.Param("id"))
	if err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	if err = sqlpkg.RemoveQuestFromDatabaseById(id, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})
}

func (rt *Rout) QuestDeleteByName(cont *gin.Context) {
	name := cont.Param("name")
	if err := sqlpkg.RemoveQuestFromDatabaseByName(name, rt.db); err != nil {
		cont.JSON(400, gin.H{"message": "error param type"})
		return
	}
	cont.JSON(200, gin.H{"message": "ok"})
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

	rout.router.PUT("/users/id/:id", rout.UserPutById)
	rout.router.PUT("/users/name/:name", rout.UserPutByName)
	rout.router.PUT("/quests/id/:id", rout.QuestPutById)
	rout.router.PUT("/quests/name/:name", rout.QuestPutByName)

	rout.router.DELETE("/users/id/:id", rout.UserDeleteById)
	rout.router.DELETE("/users/name/:name", rout.UserDeleteByName)
	rout.router.DELETE("/quests/id/:id", rout.QuestDeleteById)
	rout.router.DELETE("/quests/name/:name", rout.QuestDeleteByName)

	err := rout.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
