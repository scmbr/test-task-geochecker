package v1

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		v1.GET("/system/health", func(c *gin.Context) {
			dbStatus := "ok"
			redisStatus := "ok"
			if err := h.db.Ping(); err != nil {
				dbStatus = "fail"
			}
			if err := h.redis.Ping(context.Background()).Err(); err != nil {
				redisStatus = "fail"
			}
			status := http.StatusOK
			if dbStatus != "ok" {
				status = http.StatusServiceUnavailable
			}

			c.JSON(status, gin.H{
				"service": "ok",
				"db":      dbStatus,
				"redis":   redisStatus,
				"time":    time.Now(),
			})
		})
		h.initIncidentsRoutes(v1)
		h.initCheckRoutes(v1)
		h.initOperatorRoutes(v1)
	}
}
