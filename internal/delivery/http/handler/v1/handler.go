package v1

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task-geochecker/internal/service"
)

type Handler struct {
	service *service.Service
	db      *sql.DB
}

func NewHandler(service *service.Service, db *sql.DB) *Handler {
	return &Handler{
		service: service,
		db:      db,
	}
}
func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		v1.GET("/system/health", func(c *gin.Context) {
			dbStatus := "ok"

			if err := h.db.Ping(); err != nil {
				dbStatus = "fail"
			}

			status := http.StatusOK
			if dbStatus != "ok" {
				status = http.StatusServiceUnavailable
			}

			c.JSON(status, gin.H{
				"service": "ok",
				"db":      dbStatus,
				"time":    time.Now(),
			})
		})
		h.initIncidentsRoutes(v1)
		h.initCheckRoutes(v1)
		h.initOperatorRoutes(v1)
	}
}
