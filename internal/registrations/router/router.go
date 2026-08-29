package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wksn753/umosan-backend/internal/registrations/handler"
)

// RegisterRouter wires up your V1 routes and injects the handler
func RegisterRouter(router *gin.RouterGroup, regHandler *handler.RegistrationHandler) {

	// Registrations and attendance group
	registrations := router.Group("/registrations")
	{
		// POST /api/v1/registrations - Handles frontend form submissions
		registrations.POST("", regHandler.HandleRegister)

		// GET /api/v1/registrations/attendance?date=YYYY-MM-DD - Fetches attendance by date
		registrations.GET("/attendance", regHandler.HandleGetAttendance)

		// GET /api/v1/registrations/stats - Fetches global counts and metrics
		registrations.GET("/stats", regHandler.HandleGetStats)
	}
}
