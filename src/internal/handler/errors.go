package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"shopAPI/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type errorResponse struct {
	Message string `json:"error"`
}

func sendError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var notFound *domain.NotFoundError
	if errors.As(err, &notFound) {
		c.JSON(http.StatusNotFound, errorResponse{Message: err.Error()})
		return
	}

	var validation *domain.ValidationError
	var emptyPayload *domain.EmptyPayloadError
	if errors.As(err, &validation) || errors.As(err, &emptyPayload) {
		c.JSON(http.StatusBadRequest, errorResponse{Message: err.Error()})
		return
	}

	slog.Error("internal server error", "error", err)
	c.JSON(http.StatusInternalServerError, errorResponse{Message: "internal server error"})
}

func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		message := err.Error()
		if errors.Is(err, io.EOF) {
			message = "request body is required"
		}
		c.JSON(http.StatusBadRequest, errorResponse{Message: message})
		return false
	}
	return true
}

func parseIDFromParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: fmt.Sprintf("invalid %s parameter: expected uuid", name)})
		return uuid.Nil, false
	}
	return id, true
}

func parsePagination(c *gin.Context) (limit int, offset int, ok bool) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "limit must be an integer"})
		return 0, 0, false
	}

	offset, err = strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "offset must be an integer"})
		return 0, 0, false
	}

	return limit, offset, true
}
