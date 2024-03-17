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

type EventResp struct {
	code int
	msg  string

	//event *sqlpkg.Event
	event *sqlpkg.EventResponse
}

type Rout struct {
	router     *gin.Engine
	db         *sql.DB
	uCacheId   *lru.Cache[int, *sqlpkg.User]
	qCacheId   *lru.Cache[int, *sqlpkg.Quest]
	uCacheName *lru.Cache[string, *sqlpkg.User]
	qCacheName *lru.Cache[string, *sqlpkg.Quest]
}

func NewEventResponse(code int, msg string, event *sqlpkg.EventResponse) *EventResp {
	return &EventResp{
		code:  code,
		msg:   msg,
		event: event,
	}
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
		u := sqlpkg.User{
			Id:      e.UserId,
			Name:    "",
			Balance: e.UserBalance,
		}
		rt.UpdateUCachesById(u, u.Id)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestGetById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var q sqlpkg.Quest
		id, err := strconv.Atoi(cont.Param("id"))
		if qc, ok := rt.qCacheId.Get(id); ok {
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
		rt.UpdateUCachesById(u, id)
		result <- NewResponse(200, "ok", nil)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func AssignmentUsersForUpdateCache(dst *sqlpkg.User, src sqlpkg.User) {
	if src.Balance >= 0 {
		dst.Balance = src.Balance
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
}

func (rt *Rout) UpdateUCachesByName(u sqlpkg.User, name string) {
	uc, ok := rt.uCacheName.Get(name)
	if ok {
		AssignmentUsersForUpdateCache(uc, u)
	}
	rt.uCacheName.Remove(name)
	rt.uCacheName.Add(name, uc)
	rt.uCacheId.Remove(u.Id)
	rt.uCacheId.Add(u.Id, uc)
}

func (rt *Rout) UpdateUCachesById(u sqlpkg.User, id int) {
	uc, ok := rt.uCacheId.Get(id)
	if ok {
		AssignmentUsersForUpdateCache(uc, u)
	}
	rt.uCacheId.Remove(id)
	rt.uCacheId.Add(id, uc)
	rt.uCacheName.Remove(uc.Name)
	rt.uCacheName.Add(uc.Name, uc)
}

func AssignmentQuestsForUpdateCache(dst *sqlpkg.Quest, src sqlpkg.Quest) {
	if src.Cost >= 0 {
		dst.Cost = src.Cost
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
}

func (rt *Rout) UpdateQCachesByName(q sqlpkg.Quest, name string) {
	qc, ok := rt.qCacheName.Get(name)
	if ok {
		AssignmentQuestsForUpdateCache(qc, q)
	}
	rt.qCacheName.Remove(name)
	rt.qCacheName.Add(name, qc)
	rt.qCacheId.Remove(q.Id)
	rt.qCacheId.Add(q.Id, qc)
}

func (rt *Rout) UpdateQCachesById(q sqlpkg.Quest, id int) {
	qc, ok := rt.qCacheId.Get(id)
	if ok {
		AssignmentQuestsForUpdateCache(qc, q)
	}
	rt.qCacheId.Remove(id)
	rt.qCacheId.Add(id, qc)
	rt.qCacheName.Remove(qc.Name)
	rt.qCacheName.Add(qc.Name, qc)
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
		rt.UpdateUCachesByName(u, name)
		result <- NewResponse(200, "ok", &u)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})

}

func (rt *Rout) QuestPutById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var q sqlpkg.Quest
		q.Cost = -1
		q.Name = ""
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(400, "error binding json", nil)
			return
		}
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if err = q.UpdateDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(400, "error sql execution", nil)
			return
		}
		rt.UpdateQCachesById(q, id)
		result <- NewResponse(200, "ok", &q)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestPutByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var q sqlpkg.Quest
		q.Cost = -1
		q.Name = ""
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(400, "error binding json", nil)
			return
		}
		name := cont.Param("name")
		if err := q.UpdateDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(400, "error sql execution", nil)
			return
		}
		rt.UpdateQCachesByName(q, name)
		result <- NewResponse(200, "ok", nil)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserDeleteById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if err = sqlpkg.RemoveUserFromDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if u, ok := rt.uCacheId.Get(id); ok {
			rt.uCacheId.Remove(id)
			rt.uCacheName.Remove(u.Name)
		}
		result <- NewResponse(200, "ok", nil)
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserDeleteByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		name := cont.Param("name")
		if err := sqlpkg.RemoveUserFromDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if u, ok := rt.uCacheName.Get(name); ok {
			rt.uCacheId.Remove(u.Id)
			rt.uCacheName.Remove(name)
		}
		result <- NewResponse(200, "ok", nil)
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestDeleteById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if err = sqlpkg.RemoveQuestFromDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if q, ok := rt.qCacheId.Get(id); ok {
			rt.qCacheId.Remove(id)
			rt.qCacheName.Remove(q.Name)
		}
		result <- NewResponse(200, "ok", nil)
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestDeleteByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		name := cont.Param("name")
		if err := sqlpkg.RemoveQuestFromDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(400, "error param type", nil)
			return
		}
		if q, ok := rt.qCacheName.Get(name); ok {
			rt.qCacheId.Remove(q.Id)
			rt.qCacheName.Remove(name)
		}
		result <- NewResponse(200, "ok", nil)
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) EventsGetByUserId(cont *gin.Context) {
	result := make(chan *EventResp)
	go func() {
		var event sqlpkg.Event
		userId, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewEventResponse(400, "error param type", nil)
			return
		}
		event.UserId = userId
		res, err := event.GetByUser(rt.db)
		if err != nil {
			result <- NewEventResponse(400, "error user_id", nil)
			return
		}
		result <- NewEventResponse(200, "ok", res)
	}()
	resp := <-result
	if resp.event != nil {
		cont.IndentedJSON(http.StatusOK, resp.event)
	} else {
		cont.JSON(resp.code, gin.H{"message": resp.msg})
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
	rout.router.GET("/events/user_id/:id", rout.EventsGetByUserId)
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
