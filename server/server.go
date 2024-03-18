package server

import (
	"VK_Internship_Go/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (rt *Rout) UserPost(cont *gin.Context) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	result := make(chan *Resp)
	go func() {
		var u db.User
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
			return
		}
		if u.Balance < 0 {
			result <- NewResponse(http.StatusBadRequest, "error balance cannot be negative", nil)
			return
		}
		if u.Name == "" {
			result <- NewResponse(http.StatusBadRequest, "error name cannot be empty", nil)
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
		result <- NewResponse(http.StatusConflict, "error database indexing", nil)
		return err
	}
	result <- NewResponse(http.StatusCreated, "ok", i)
	return nil
}

func (rt *Rout) ItemReadDatabaseById(i Item, id int, result chan<- *Resp) error {
	if err := i.GetById(id, rt.db); err != nil {
		result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
		return err
	}
	result <- NewResponse(http.StatusOK, "ok", i)
	return nil
}

func (rt *Rout) ItemReadDatabaseByName(i Item, name string, result chan<- *Resp) error {
	if err := i.GetByName(name, rt.db); err != nil {
		result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
		return err
	}
	result <- NewResponse(http.StatusOK, "ok", i)
	return nil
}

func (rt *Rout) QuestPost(cont *gin.Context) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	result := make(chan *Resp)
	go func() {
		var q db.Quest
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
			return
		}
		if q.Cost < 0 {
			result <- NewResponse(http.StatusBadRequest, "error cost cannot be negative", nil)
			return
		}
		if q.Name == "" {
			result <- NewResponse(http.StatusBadRequest, "error name cannot be empty", nil)
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
	rt.mu.Lock()
	defer rt.mu.Unlock()
	result := make(chan *Resp)
	go func() {
		var e db.Event
		if err := cont.BindJSON(&e); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
			return
		}
		if e.AppendDatabase(rt.db) != nil {
			result <- NewResponse(http.StatusBadRequest, "error database indexing", nil)
			return
		} else {
			result <- NewResponse(http.StatusCreated, "ok", nil)
		}
		u := db.User{
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
		var q db.Quest
		id, err := strconv.Atoi(cont.Param("id"))
		if qc, ok := rt.qCacheId.Get(id); ok {
			result <- NewResponse(http.StatusOK, "ok", qc)
			return
		}
		if err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
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
		var u db.User
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		if uc, ok := rt.uCacheId.Get(id); ok {
			result <- NewResponse(http.StatusOK, "ok", uc)
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
		var u db.User
		name := cont.Param("name")
		if uc, ok := rt.uCacheName.Get(name); ok {
			result <- NewResponse(http.StatusOK, "ok", uc)
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
		var q db.Quest
		name := cont.Param("name")
		if qc, ok := rt.qCacheName.Get(name); ok {
			result <- NewResponse(http.StatusOK, "ok", qc)
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

func (rt *Rout) UserPutById(cont *gin.Context) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	result := make(chan *Resp)
	go func() {
		var u db.User
		u.Balance = -1
		u.Name = ""
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
			return
		}
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type json", nil)
			return
		}
		if err = u.UpdateDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error sql execution", nil)
			return
		}
		rt.RemoveFromUCacheById(id)
		result <- NewResponse(http.StatusOK, "ok", nil)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserPutByName(cont *gin.Context) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	result := make(chan *Resp)
	go func() {
		var u db.User
		u.Balance = -1
		u.Name = ""
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
			return
		}
		name := cont.Param("name")
		if err := u.UpdateDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error sql execution", nil)
			return
		}
		rt.RemoveFromUCacheByName(name)
		result <- NewResponse(http.StatusOK, "ok", &u)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})

}

func (rt *Rout) QuestPutById(cont *gin.Context) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	result := make(chan *Resp)
	go func() {
		var q db.Quest
		q.Cost = -1
		q.Name = ""
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
			return
		}
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		if err = q.UpdateDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusConflict, "error sql execution", nil)
			return
		}
		rt.UpdateQCachesById(q, id)
		result <- NewResponse(http.StatusOK, "ok", &q)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestPutByName(cont *gin.Context) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	result := make(chan *Resp)
	go func() {
		var q db.Quest
		q.Cost = -1
		q.Name = ""
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error binding json", nil)
			return
		}
		name := cont.Param("name")
		if err := q.UpdateDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusConflict, "error sql execution", nil)
			return
		}
		rt.UpdateQCachesByName(q, name)
		result <- NewResponse(http.StatusOK, "ok", nil)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserDeleteById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		if err = db.RemoveUserFromDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusConflict, "error param type", nil)
			return
		}
		if u, ok := rt.uCacheId.Get(id); ok {
			rt.uCacheId.Remove(id)
			rt.uCacheName.Remove(u.Name)
		}
		result <- NewResponse(http.StatusCreated, "ok", nil)
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserDeleteByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		name := cont.Param("name")
		if err := db.RemoveUserFromDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		if u, ok := rt.uCacheName.Get(name); ok {
			rt.uCacheId.Remove(u.Id)
			rt.uCacheName.Remove(name)
		}
		result <- NewResponse(http.StatusCreated, "ok", nil)
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestDeleteById(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		id, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		if err = db.RemoveQuestFromDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		if q, ok := rt.qCacheId.Get(id); ok {
			rt.qCacheId.Remove(id)
			rt.qCacheName.Remove(q.Name)
		}
		result <- NewResponse(http.StatusCreated, "ok", nil)
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestDeleteByName(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		name := cont.Param("name")
		if err := db.RemoveQuestFromDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		if q, ok := rt.qCacheName.Get(name); ok {
			rt.qCacheId.Remove(q.Id)
			rt.qCacheName.Remove(name)
		}
	}()
	resp := <-result
	cont.IndentedJSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) EventsGetByUserId(cont *gin.Context) {
	result := make(chan *EventResp)
	go func() {
		var event db.Event
		userId, err := strconv.Atoi(cont.Param("id"))
		if err != nil {
			result <- NewEventResponse(http.StatusBadRequest, "error param type", nil)
			return
		}
		event.UserId = userId
		res, err := event.GetByUser(rt.db)
		if err != nil {
			result <- NewEventResponse(http.StatusConflict, "error user_id", nil)
			return
		}
		result <- NewEventResponse(http.StatusOK, "ok", res)
	}()
	resp := <-result
	if resp.event != nil {
		cont.IndentedJSON(http.StatusOK, resp.event)
	} else {
		cont.JSON(resp.code, gin.H{"message": resp.msg})
	}
}

func HttpServer() {
	rout := NewRout(gin.Default(), db.NewConn())

	rout.router.POST("/users", rout.UserPost)
	rout.router.POST("/quests", rout.QuestPost)
	rout.router.POST("/event", rout.EventPost)

	rout.router.GET("/users/id/:id", rout.UserGetById)
	rout.router.GET("/users/name/:name", rout.UserGetByName)
	rout.router.GET("/quests/id/:id", rout.QuestGetById)
	rout.router.GET("/quests/name/:name", rout.QuestGetByName)
	rout.router.GET("/events/user_id/:id", rout.EventsGetByUserId)

	rout.router.PUT("/users/id/:id", rout.UserPutById)
	rout.router.PUT("/users/name/:name", rout.UserPutByName)
	rout.router.PUT("/quests/id/:id", rout.QuestPutById)
	rout.router.PUT("/quests/name/:name", rout.QuestPutByName)

	rout.router.DELETE("/users/id/:id", rout.UserDeleteById)
	rout.router.DELETE("/users/name/:name", rout.UserDeleteByName)
	rout.router.DELETE("/quests/id/:id", rout.QuestDeleteById)
	rout.router.DELETE("/quests/name/:name", rout.QuestDeleteByName)

	err := rout.router.Run("0.0.0.0:8080")
	if err != nil {
		log.Println(err)
	}
}
