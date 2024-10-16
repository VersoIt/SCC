package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getEmployeeById(c *gin.Context) {
	employeeIdString := c.Param("id")
	employeeId, err := strconv.Atoi(employeeIdString)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid employee id"))
		return
	}

	employee, err := h.service.GetById(employeeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) getAllEmployees(c *gin.Context) {
	employees, err := h.service.Employee.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, employees)
}
