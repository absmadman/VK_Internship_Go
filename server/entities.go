package server

import (
	"VK_Internship_Go/db"
	"database/sql"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
	"log"
	"os"
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
	code  int
	msg   string
	event *db.EventResponse
}

type Rout struct {
	router     *gin.Engine
	db         *sql.DB
	uCacheId   *lru.Cache[int, *db.User]
	qCacheId   *lru.Cache[int, *db.Quest]
	uCacheName *lru.Cache[string, *db.User]
	qCacheName *lru.Cache[string, *db.Quest]
}

func NewEventResponse(code int, msg string, event *db.EventResponse) *EventResp {
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

func NewRout(g *gin.Engine, d *sql.DB) *Rout {
	cacheSize, err := strconv.Atoi(os.Getenv("SERVER_CACHE_SIZE"))
	if err != nil {
		log.Println("Invalid env type for SERVER_CACHE_SIZE")
	}
	ucId, _ := lru.New[int, *db.User](cacheSize)
	qcId, _ := lru.New[int, *db.Quest](cacheSize)
	ucName, _ := lru.New[string, *db.User](cacheSize)
	qcName, _ := lru.New[string, *db.Quest](cacheSize)
	return &Rout{
		router:     g,
		db:         d,
		uCacheId:   ucId,
		qCacheId:   qcId,
		uCacheName: ucName,
		qCacheName: qcName,
	}
}
