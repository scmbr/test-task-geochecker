package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	handler_dto "github.com/scmbr/test-task-geochecker/internal/delivery/http/dto"
	"github.com/scmbr/test-task-geochecker/internal/service"
	service_dto "github.com/scmbr/test-task-geochecker/internal/service/dto"
	"github.com/scmbr/test-task-geochecker/pkg/logger"
)

func (h *Handler) initOperatorRoutes(api *gin.RouterGroup) {
	operators := api.Group("/operators", h.operatorIdentity)
	{
		operators.POST("", h.createOperator)
		operators.DELETE("/:id", h.revokeOperator)
	}
}
func (h *Handler) createOperator(c *gin.Context) {
	var input handler_dto.CreateOperatorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Operator.Create(c.Request.Context(), &service_dto.CreateOperatorInput{
		Name:   input.Name,
		APIKey: input.APIKey,
	}); err != nil {
		logger.Error(
			"error occurred while creating an operator",
			err,
			map[string]interface{}{
				"name":    input.Name,
				"api_key": maskString(input.APIKey),
			},
		)
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusCreated)
}
func (h *Handler) revokeOperator(c *gin.Context) {
	id := c.Param("id")
	callerID := c.GetString(operatorCtx)
	if callerID == id {
		newResponse(c, http.StatusForbidden, "cannot revoke yourself")
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}
	if err := h.service.Operator.Revoke(c.Request.Context(), id); err != nil {
		logger.Error(
			"error occurred while revoking an operator",
			err,
			map[string]interface{}{
				"operator_id": id,
				"caller_id":   callerID,
			},
		)
		if errors.Is(err, service.ErrOperatorNotFound) {
			newResponse(c, http.StatusNotFound, "operator not found")
			return
		}
		newResponse(c, http.StatusInternalServerError, "something went wrong")
		return
	}
	c.Status(http.StatusNoContent)
}
