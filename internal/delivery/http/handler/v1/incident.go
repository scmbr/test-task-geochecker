package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_dto "github.com/scmbr/test-task-geochecker/internal/delivery/http/dto"
	"github.com/scmbr/test-task-geochecker/internal/service"
	service_dto "github.com/scmbr/test-task-geochecker/internal/service/dto"
	"github.com/scmbr/test-task-geochecker/pkg/logger"
)

func (h *Handler) initIncidentsRoutes(api *gin.RouterGroup) {
	incidents := api.Group("/incidents", h.operatorIdentity)
	{
		incidents.POST("", h.createIncident)
		incidents.GET("", h.getAllIncidents)
		incidents.GET("/:id", h.getIncidentById)
		incidents.GET("/stats/:id", h.getIncidentStatsById)
		incidents.PUT("/:id", h.updateIncidentById)
		incidents.DELETE("/:id", h.deleteIncidentById)
	}
}
func (h *Handler) createIncident(c *gin.Context) {
	operatorIDRaw, ok := c.Get(operatorCtx)
	if !ok {
		newResponse(c, http.StatusUnauthorized, "operator not authorized")
		return
	}

	operatorID, ok := operatorIDRaw.(string)
	if !ok {
		newResponse(c, http.StatusInternalServerError, "invalid operator context")
		return
	}
	var input handler_dto.CreateIncidentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.Incident.Create(c.Request.Context(), &service_dto.CreateIncidentInput{
		OperatorID: operatorID,
		Longitude:  input.Longitude,
		Latitude:   input.Latitude,
		Radius:     input.Radius,
	})
	if err != nil {
		logger.Error(
			"error occurred while creating an incident",
			err,
			map[string]interface{}{
				"operator_id": operatorID,
				"longitude":   input.Longitude,
				"latitude":    input.Latitude,
				"radius":      input.Radius,
			},
		)
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusCreated)
}
func (h *Handler) getAllIncidents(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		newResponse(c, http.StatusBadRequest, "invalid limit")
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		newResponse(c, http.StatusBadRequest, "invalid offset")
		return
	}
	res, err := h.service.Incident.GetAll(c.Request.Context(), &service_dto.GetAllIncidentsInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error(
			"error occurred while getting all incidents",
			err,
			map[string]interface{}{
				"limit":  limit,
				"offset": offset,
			},
		)
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	incidents := make([]handler_dto.GetIncidentResponse, len(res.Incidents))
	for idx, i := range res.Incidents {
		incidents[idx] = handler_dto.GetIncidentResponse{
			IncidentID: i.ID,
			Latitude:   i.Latitude,
			Longitude:  i.Longitude,
			Radius:     i.Radius,
			CreatedAt:  i.CreatedAt,
		}
	}
	c.JSON(http.StatusOK, handler_dto.GetAllIncidentsResponse{
		Total:     int32(res.Total),
		Incidents: incidents,
	})

}
func (h *Handler) getIncidentById(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	res, err := h.service.Incident.GetById(c.Request.Context(), id)
	if err != nil {
		logger.Error(
			"error occurred while getting incident by id",
			err,
			map[string]interface{}{
				"incident_id": id,
			},
		)
		if errors.Is(err, service.ErrIncidentNotFound) {
			newResponse(c, http.StatusNotFound, "incident not found")
			return
		}
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	incident := handler_dto.GetIncidentResponse{
		IncidentID: res.ID,
		Latitude:   res.Latitude,
		Longitude:  res.Longitude,
		Radius:     res.Radius,
		CreatedAt:  res.CreatedAt,
	}
	c.JSON(http.StatusOK, incident)

}
func (h *Handler) updateIncidentById(c *gin.Context) {
	var input handler_dto.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	if err := h.service.Incident.Update(c.Request.Context(), id, &service_dto.UpdateIncidentInput{
		OperatorID: input.OperatorID,
		Latitude:   input.Latitude,
		Longitude:  input.Longitude,
		Radius:     input.Radius,
	}); err != nil {
		logger.Error(
			"error occurred while updating incident by id",
			err,
			map[string]interface{}{
				"incident_id": id,
				"operator_id": input.OperatorID,
				"longitude":   input.Longitude,
				"latitude":    input.Latitude,
				"radius":      input.Radius,
			},
		)
		if errors.Is(err, service.ErrIncidentNotFound) {
			newResponse(c, http.StatusNotFound, "incident not found")
			return
		}
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) deleteIncidentById(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	if err := h.service.Incident.Delete(c.Request.Context(), id); err != nil {
		logger.Error(
			"error occurred while deleting incident by id",
			err,
			map[string]interface{}{
				"incident_id": id,
			},
		)
		if errors.Is(err, service.ErrIncidentNotFound) {
			newResponse(c, http.StatusNotFound, "incident not found")
			return
		}
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) getIncidentStatsById(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	sinceMinutesStr := c.DefaultQuery("since", "20")
	sinceMinutes, err := strconv.Atoi(sinceMinutesStr)
	if err != nil || sinceMinutes <= 0 {
		newResponse(c, http.StatusBadRequest, "invalid since_minutes")
		return
	}

	since := time.Now().Add(-time.Duration(sinceMinutes) * time.Minute)
	count, err := h.service.Incident.GetStats(c.Request.Context(), id, since)
	if err != nil {
		logger.Error(
			"error occurred while get incident's stats by id",
			err,
			map[string]interface{}{
				"incident_id": id,
			},
		)
		if errors.Is(err, service.ErrIncidentNotFound) {
			newResponse(c, http.StatusNotFound, "incident not found")
			return
		}
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.JSON(http.StatusOK, handler_dto.GetIncidentStatsByIdResponse{
		IncidentID:   id,
		UserCount:    count,
		SinceMinutes: since,
	})
}
