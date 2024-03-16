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

type Item interface {
	AppendDatabase(db *sql.DB) error
	UpdateDatabaseById(id int, db *sql.DB) error
	GetById(id int, db *sql.DB) error
	GetByName(name string, db *sql.DB) error
}

type Resp struct {
	code int
	msg  string
	item Item
}

type Rout struct {
	router     *gin.Engine
	db         *sql.DB
	uCacheId   *lru.Cache[int, *sqlpkg.User]
	qCacheId   *lru.Cache[int, *sqlpkg.Quest]
	uCacheName *lru.Cache[string, *sqlpkg.User]
	qCacheName *lru.Cache[string, *sqlpkg.Quest]
}

func NewResponse(code int, msg string, item Item) *Resp {
	return &Resp{
		code: code,
		msg:  msg,
		item: item,
	}
}

func (rt *Rout) UserPost(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var u sqlpkg.User
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(400, "error binding json", nil)
			return
		}
		if u.Balance < 0 {
			result <- NewResponse(400, "error balance cannot be negative", nil)
			return
		}
		if u.Name == "" {
			result <- NewResponse(400, "error name cannot be empty", nil)
			return
		}
		if err := rt.ItemAppendDatabase(&u, result); err == nil {
			rt.uCacheId.Add(u.Id, &u)
			rt.uCacheName.Add(u.Name, &u)
		}
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) ItemAppendDatabase(i Item, result chan<- *Resp) error {
	if err := i.AppendDatabase(rt.db); err != nil {
		result <- NewResponse(400, "error database indexing", nil)
		return err
	}
	result <- NewResponse(200, "ok", i)
	return nil
}

func (rt *Rout) ItemReadDatabaseById(i Item, id int, result chan<- *Resp) error {
	if err := i.GetById(id, rt.db); err != nil {
		result <- NewResponse(400, "error binding json", nil)
		return err
	}
	result <- NewResponse(200, "ok", i)
	return nil
}

func (rt *Rout) ItemReadDatabaseByName(i Item, name string, result chan<- *Resp) error {
	if err := i.GetByName(name, rt.db); err != nil {
		result <- NewResponse(400, "error binding json", nil)
		return err
	}
	result <- NewResponse(200, "ok", i)
	return nil
}

func (rt *Rout) QuestPost(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var q sqlpkg.Quest
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(400, "error binding json", nil)
			return
		}
		if q.Cost < 0 {
			result <- NewResponse(400, "error cost cannot be negative", nil)
			return
		}
		if q.Name == "" {
			result <- NewResponse(400, "error name cannot be empty", nil)
			return
		}
		if err := rt.ItemAppendDatabase(&q, result); err == nil {
			rt.qCacheId.Add(q.Id, &q)
			rt.qCacheName.Add(q.Name, &q)
		}
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})

}
func (rt *Rout) EventPost(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var e sqlpkg.Event
		if err := cont.BindJSON(&e); err != nil {
			result <- NewResponse(400, "error binding json", nil)
			return
		}
		if e.AppendDatabase(rt.db) != nil {
			result <- NewResponse(400, "error database indexing", nil)
		} else {
			result <- NewResponse(200, "ok", nil)
		}
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserGetById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var u sqlpkg.User
		id, err := strconv.Atoi(cont.Param("id"))
		if uc, ok := rt.uCacheId.Get(id); ok {
			result <- NewResponse(200, "ok", uc)
			return
		}
		if err != nil {
			result <- NewResponse(400, "error param indexing", nil)
			return
		}

		if err = rt.ItemReadDatabaseById(&u, id, result); err == nil {
			rt.uCacheId.Add(id, &u)
			rt.uCacheName.Add(u.Name, &u)
		}
	}()
	resp := <-result
	if resp.item == nil {
		cont.JSON(resp.code, gin.H{"message": resp.msg})
	} else {
		cont.IndentedJSON(http.StatusOK, resp.item)
	}
}

func (rt *Rout) UserGetByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var u sqlpkg.User
		name := cont.Param("name")
		if uc, ok := rt.uCacheName.Get(name); ok {
			result <- NewResponse(200, "ok", uc)
			return
		}
		if err := rt.ItemReadDatabaseByName(&u, name, result); err == nil {
			rt.uCacheName.Add(name, &u)
			rt.uCacheId.Add(u.Id, &u)
		}

	}()
	resp := <-result
	if resp.item == nil {
		cont.JSON(resp.code, gin.H{"message": resp.msg})
	} else {
		cont.IndentedJSON(http.StatusOK, resp.item)
	}
}

func (rt *Rout) QuestGetById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var q sqlpkg.Quest
		id, err := strconv.Atoi(cont.Param("id"))
		if qc, ok := rt.uCacheId.Get(id); ok {
			result <- NewResponse(200, "ok", qc)
			return
		}
		if err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if err = rt.ItemReadDatabaseById(&q, id, result); err == nil {
			rt.qCacheId.Add(id, &q)
			rt.qCacheName.Add(q.Name, &q)
		}
	}()
	resp := <-result
	if resp.item == nil {
		cont.JSON(resp.code, gin.H{"message": resp.msg})
	} else {
		cont.IndentedJSON(http.StatusOK, resp.item)
	}
}

func (rt *Rout) QuestGetByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var q sqlpkg.Quest
		name := cont.Param("name")
		if qc, ok := rt.qCacheName.Get(name); ok {
			result <- NewResponse(200, "ok", qc)
			return
		}
		if err := rt.ItemReadDatabaseByName(&q, name, result); err == nil {
			rt.qCacheId.Add(q.Id, &q)
			rt.qCacheName.Add(name, &q)
		}
	}()
	resp := <-result
	if resp.item == nil {
		cont.JSON(resp.code, gin.H{"message": resp.msg})
	} else {
		cont.IndentedJSON(http.StatusOK, resp.item)
	}
}

func NewRout(g *gin.Engine, d *sql.DB) *Rout {
	ucId, _ := lru.New[int, *sqlpkg.User](128)
	qcId, _ := lru.New[int, *sqlpkg.Quest](128)
	ucName, _ := lru.New[string, *sqlpkg.User](128)
	qcName, _ := lru.New[string, *sqlpkg.Quest](128)
	return &Rout{
		router:     g,
		db:         d,
		uCacheId:   ucId,
		qCacheId:   qcId,
		uCacheName: ucName,
		qCacheName: qcName,
	}
}

func (rt *Rout) UserPutById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var u sqlpkg.User
		u.Balance = -1
		u.Name = ""
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(400, "error binding json", nil)
			return
		}
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(400, "error param type json", nil)
			return
		}
		if err = u.UpdateDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(400, "error sql execution", nil)
			return
		}
		result <- NewResponse(200, "ok", nil)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})

}

func (rt *Rout) UserPutByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var u sqlpkg.User
		u.Balance = -1
		u.Name = ""
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(400, "error binding json", nil)
			return
		}
		name := cont.Param("name")
		if err := u.UpdateDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(400, "error sql execution", nil)
			return
		}
		result <- NewResponse(200, "ok", &u)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})

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
