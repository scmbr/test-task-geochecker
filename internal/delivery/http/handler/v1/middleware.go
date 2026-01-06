package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(operatorCtx, operator.OperatorID)
}
func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Api-Key" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
