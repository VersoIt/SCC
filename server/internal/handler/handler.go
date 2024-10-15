package handler

import (
	"github.com/gin-gonic/gin"
	"serverClientClient/server/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		employee := api.Group("/employee")
		{
			employee.GET(":id", h.getEmployeeById)
		}
	}
	return router
}
