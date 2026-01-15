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
	incidents, err := h.service.Check.Check(c.Request.Context(), &service_dto.CheckInput{
		UserID:    input.UserID,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	nearbyIncidents := make([]handler_dto.GetIncidentResponse, len(incidents))
	for idx, i := range incidents {
		nearbyIncidents[idx] = handler_dto.GetIncidentResponse{
			IncidentID: i.ID,
			OperatorID: i.OperatorID,
			Latitude:   i.Latitude,
			Longitude:  i.Longitude,
			Radius:     i.Radius,
			CreatedAt:  i.CreatedAt,
		}
	}

	c.JSON(http.StatusCreated, handler_dto.CreateCheckResponse{
		Count:     uint16(len(nearbyIncidents)),
		Incidents: nearbyIncidents,
	})
}
