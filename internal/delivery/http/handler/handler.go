package handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	v1 "github.com/scmbr/test-task-geochecker/internal/delivery/http/handler/v1"
	"github.com/scmbr/test-task-geochecker/internal/service"
)

type Handler struct {
	service *service.Service
	db      *sql.DB
	redis   *redis.Client
}

func NewHandler(service *service.Service, db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{
		service: service,
		db:      db,
		redis:   redis,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery(), gin.Logger())
	h.initAPI(router)
	return router
}
func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.service, h.db, h.redis)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
