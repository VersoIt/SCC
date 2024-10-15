package handler

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, err error) {
	logrus.Errorf(err.Error())
	if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{Message: "element not found"})
		return
	}
	c.AbortWithStatusJSON(statusCode, errorResponse{err.Error()})
}
