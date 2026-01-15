package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task-geochecker/pkg/logger"
)

const (
	authorizationHeader = "Authorization"
	operatorCtx         = "operatorId"
)

func (h *Handler) operatorIdentity(c *gin.Context) {
	apiKey, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	operator, err := h.service.Operator.ValidateAPIKey(c, apiKey)
	if err != nil {
		logger.Error(
			"error occurred while validating an api key",
			err,
			map[string]interface{}{
				"api_key": maskString(apiKey),
			},
		)
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(operatorCtx, operator.OperatorID)
}
func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		return "", errors.New("empty auth header")
	}

	return apiKey, nil
}
