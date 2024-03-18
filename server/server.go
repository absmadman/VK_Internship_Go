package server

import (
	"VK_Internship_Go/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (rt *Rout) UserPost(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var u db.User
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(http.StatusBadRequest, bindingError, nil)
			return
		}
		if u.Balance < 0 {
			result <- NewResponse(http.StatusBadRequest, balanceError, nil)
			return
		}
		if u.Name == "" {
			result <- NewResponse(http.StatusBadRequest, nameError, nil)
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
		result <- NewResponse(http.StatusConflict, dbError, nil)
		return err
	}
	result <- NewResponse(http.StatusCreated, "ok", i)
	return nil
}

func (rt *Rout) ItemReadDatabaseById(i Item, id int, result chan<- *Resp) error {
	if err := i.GetById(id, rt.db); err != nil {
		result <- NewResponse(http.StatusBadRequest, bindingError, nil)
		return err
	}
	result <- NewResponse(http.StatusOK, "ok", i)
	return nil
}

func (rt *Rout) ItemReadDatabaseByName(i Item, name string, result chan<- *Resp) error {
	if err := i.GetByName(name, rt.db); err != nil {
		result <- NewResponse(http.StatusBadRequest, bindingError, nil)
		return err
	}
	result <- NewResponse(http.StatusOK, "ok", i)
	return nil
}

func (rt *Rout) QuestPost(cont *gin.Context) {
	result := make(chan *Resp)
	go func() {
		var q db.Quest
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(http.StatusBadRequest, bindingError, nil)
			return
		}
		if q.Cost < 0 {
			result <- NewResponse(http.StatusBadRequest, costError, nil)
			return
		}
		if q.Name == "" {
			result <- NewResponse(http.StatusBadRequest, nameError, nil)
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
		var e db.Event
		if err := cont.BindJSON(&e); err != nil {
			result <- NewResponse(http.StatusBadRequest, bindingError, nil)
			return
		}
		if e.AppendDatabase(rt.db) != nil {
			result <- NewResponse(http.StatusBadRequest, dbError, nil)
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

func (rt *Rout) QuestGetById(cont *gin.Context, id int) {
	result := make(chan *Resp)
	go func() {
		var q db.Quest
		if qc, ok := rt.qCacheId.Get(id); ok {
			result <- NewResponse(http.StatusOK, "ok", qc)
			return
		}
		if err := rt.ItemReadDatabaseById(&q, id, result); err == nil {
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

func (rt *Rout) UserGetById(cont *gin.Context, id int) {
	result := make(chan *Resp)
	go func() {
		var u db.User
		if uc, ok := rt.uCacheId.Get(id); ok {
			result <- NewResponse(http.StatusOK, "ok", uc)
			return
		}
		if err := rt.ItemReadDatabaseById(&u, id, result); err == nil {
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

func (rt *Rout) UserGetByName(cont *gin.Context, name string) {
	result := make(chan *Resp)
	go func() {
		var u db.User
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

func (rt *Rout) QuestGetByName(cont *gin.Context, name string) {
	result := make(chan *Resp)
	go func() {
		var q db.Quest
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

func (rt *Rout) UserPutById(cont *gin.Context, id int) {
	result := make(chan *Resp)
	go func() {
		var u db.User
		u.Balance = -1
		u.Name = ""
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(http.StatusBadRequest, bindingError, nil)
			return
		}
		if err := u.UpdateDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, dbError, nil)
			return
		}
		rt.RemoveFromUCacheById(id)
		result <- NewResponse(http.StatusOK, "ok", nil)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserPutByName(cont *gin.Context, name string) {
	result := make(chan *Resp)
	go func() {
		var u db.User
		u.Balance = -1
		u.Name = ""
		if err := cont.BindJSON(&u); err != nil {
			result <- NewResponse(http.StatusBadRequest, bindingError, nil)
			return
		}
		if err := u.UpdateDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, dbError, nil)
			return
		}
		rt.RemoveFromUCacheByName(name)
		result <- NewResponse(http.StatusOK, "ok", &u)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})

}

func (rt *Rout) QuestPutById(cont *gin.Context, id int) {
	result := make(chan *Resp)
	go func() {
		var q db.Quest
		q.Cost = -1
		q.Name = ""
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(http.StatusBadRequest, bindingError, nil)
			return
		}
		if err := q.UpdateDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusConflict, dbError, nil)
			return
		}
		rt.UpdateQCachesById(q, id)
		result <- NewResponse(http.StatusOK, "ok", &q)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) QuestPutByName(cont *gin.Context, name string) {
	result := make(chan *Resp)
	go func() {
		var q db.Quest
		q.Cost = -1
		q.Name = ""
		if err := cont.BindJSON(&q); err != nil {
			result <- NewResponse(http.StatusBadRequest, bindingError, nil)
			return
		}
		if err := q.UpdateDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusConflict, dbError, nil)
			return
		}
		rt.UpdateQCachesByName(q, name)
		result <- NewResponse(http.StatusOK, "ok", nil)
	}()
	resp := <-result
	cont.JSON(resp.code, gin.H{"message": resp.msg})
}

func (rt *Rout) UserDeleteById(cont *gin.Context, id int) {
	result := make(chan *Resp)
	go func() {
		if err := db.RemoveUserFromDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusConflict, paramError, nil)
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

func (rt *Rout) UserDeleteByName(cont *gin.Context, name string) {
	result := make(chan *Resp)
	go func() {
		if err := db.RemoveUserFromDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, paramError, nil)
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

func (rt *Rout) QuestDeleteById(cont *gin.Context, id int) {
	result := make(chan *Resp)
	go func() {
		if err := db.RemoveQuestFromDatabaseById(id, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, paramError, nil)
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

func (rt *Rout) QuestDeleteByName(cont *gin.Context, name string) {
	result := make(chan *Resp)
	go func() {
		if err := db.RemoveQuestFromDatabaseByName(name, rt.db); err != nil {
			result <- NewResponse(http.StatusBadRequest, paramError, nil)
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

func (rt *Rout) EventsGetByUserId(cont *gin.Context, userId int) {
	result := make(chan *EventResp)
	go func() {
		var event db.Event
		event.UserId = userId
		res, err := event.GetByUser(rt.db)
		if err != nil {
			result <- NewEventResponse(http.StatusConflict, userIdError, nil)
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

func (rt *Rout) Wrapper(cont *gin.Context, f1 func(cont *gin.Context, id int), f2 func(cont *gin.Context, name string)) {
	if param, ok := cont.GetQuery("id"); ok {
		id, err := strconv.Atoi(param)
		if err != nil {
			cont.JSON(http.StatusBadRequest, gin.H{"message": paramError})
			return
		}
		f1(cont, id)
	} else if param, ok = cont.GetQuery("name"); ok {
		f2(cont, param)
	} else {
		cont.JSON(http.StatusBadRequest, gin.H{"message": paramError})
	}
}

func (rt *Rout) GetUser(cont *gin.Context) {
	rt.Wrapper(cont, rt.UserGetById, rt.UserGetByName)
}

func (rt *Rout) GetQuest(cont *gin.Context) {
	rt.Wrapper(cont, rt.QuestGetById, rt.QuestGetByName)
}

func (rt *Rout) PutUser(cont *gin.Context) {
	rt.Wrapper(cont, rt.UserPutById, rt.UserPutByName)
}

func (rt *Rout) PutQuest(cont *gin.Context) {
	rt.Wrapper(cont, rt.QuestPutById, rt.QuestPutByName)
}

func (rt *Rout) GetEvent(cont *gin.Context) {
	if param, ok := cont.GetQuery("user_id"); ok {
		id, err := strconv.Atoi(param)
		if err != nil {
			cont.JSON(http.StatusBadRequest, gin.H{"message": paramError})
			return
		}
		rt.EventsGetByUserId(cont, id)
	} else {
		cont.JSON(http.StatusBadRequest, gin.H{"message": paramError})
	}
}

func (rt *Rout) DeleteUser(cont *gin.Context) {
	rt.Wrapper(cont, rt.UserDeleteById, rt.UserDeleteByName)
}

func (rt *Rout) DeleteQuest(cont *gin.Context) {
	rt.Wrapper(cont, rt.QuestDeleteById, rt.QuestDeleteByName)
}

func HttpServer() {
	rout := NewRout(gin.Default(), db.NewConn())
	rout.router.POST("/users", rout.UserPost)
	rout.router.POST("/quests", rout.QuestPost)
	rout.router.POST("/event", rout.EventPost)
	rout.router.GET("/users", rout.GetUser)
	rout.router.GET("/quests", rout.GetQuest)
	rout.router.GET("/events", rout.GetEvent)
	rout.router.PUT("/quests", rout.PutQuest)
	rout.router.PUT("/users", rout.PutUser)
	rout.router.DELETE("/users", rout.DeleteUser)
	rout.router.DELETE("/quests", rout.DeleteQuest)
	err := rout.router.Run("0.0.0.0:8080")
	if err != nil {
		log.Println(err)
	}
}
