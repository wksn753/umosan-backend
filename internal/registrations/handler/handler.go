package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wksn753/umosan-backend/internal/registrations/services"
)

type RegistrationHandler struct {
	Service *services.RegistrationService
}

func NewRegistrationHandler(service *services.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{Service: service}
}

// HandleRegister processes incoming registration form submissions from the frontend
func (h *RegistrationHandler) HandleRegister(ctx *gin.Context) {
	var req services.RegisterMemberRequest

	// Bind and validate JSON payload against the request struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
		return
	}

	// Invoke the service layer to check/create member and save record
	if err := h.Service.RegisterOrCheckIn(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    "Registration received successfully.",
		"error":   nil,
	})
}

// HandleGetAttendance fetches attendance records based on a date query parameter (?date=YYYY-MM-DD)
func (h *RegistrationHandler) HandleGetAttendance(ctx *gin.Context) {
	dateStr := ctx.Query("date")
	if dateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"error":   "date query parameter (YYYY-MM-DD) is required",
		})
		return
	}

	records, err := h.Service.GetAttendanceByDate(dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    records,
		"error":   nil,
	})
}

// HandleGetStats fetches the overall application statistics
func (h *RegistrationHandler) HandleGetStats(ctx *gin.Context) {
	stats, err := h.Service.GetMemberStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"error":   nil,
	})
}
