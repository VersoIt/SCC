package http

import (
	"github.com/gin-gonic/gin"
	"serverClientClient/internal/service"
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
		employee := api.Group("/employee", h.useCORS)
		{
			employee.GET(":id", h.getEmployeeById)
			employee.GET("/", h.getAllEmployees)
		}
	}
	return router
}
