package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler_dto "github.com/scmbr/test-task-geochecker/internal/delivery/http/dto"
	service_dto "github.com/scmbr/test-task-geochecker/internal/service/dto"
)

func (h *Handler) initCheckRoutes(api *gin.RouterGroup) {
	checks := api.Group("/checks")
	{
		checks.POST("", h.createCheck)
	}
}
func (h *Handler) createCheck(c *gin.Context) {
	var input handler_dto.CreateCheckRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Check.Check(c.Request.Context(), &service_dto.CheckInput{
		UserID:    input.UserID,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusCreated)
}
