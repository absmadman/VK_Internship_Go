package server

import (
	"VK_Internship_Go/db"
	"database/sql"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
	"sync"
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
	mu         sync.Mutex
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
	//cacheSize, err := strconv.Atoi(os.Getenv("SERVER_CACHE_SIZE"))
	//if {

	//}
	ucId, _ := lru.New[int, *db.User](128)
	qcId, _ := lru.New[int, *db.Quest](128)
	ucName, _ := lru.New[string, *db.User](128)
	qcName, _ := lru.New[string, *db.Quest](128)
	return &Rout{
		router:     g,
		db:         d,
		uCacheId:   ucId,
		qCacheId:   qcId,
		uCacheName: ucName,
		qCacheName: qcName,
	}
}
