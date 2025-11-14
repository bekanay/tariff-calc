package handler

import (
	"net/http"
	"tariff-api/internal/service"

	"github.com/gin-gonic/gin"
)

type StationHandler struct {
	service *service.StationService
}

func NewStationHandler(service *service.StationService) *StationHandler {
	return &StationHandler{service: service}
}

func (h *StationHandler) GetStations(c *gin.Context) {
	stations, err := h.service.GetStations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stations)
}
