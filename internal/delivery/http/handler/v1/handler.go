package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task-geochecker/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initIncidentsRoutes(v1)
		h.initCheckRoutes(v1)
		h.initOperatorRoutes(v1)
	}
}
