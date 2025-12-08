package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tariff-api/internal/model"
	"tariff-api/internal/service"

	"github.com/gin-gonic/gin"
)

type StationHandler struct {
	service *service.StationService
}

func NewStationHandler(service *service.StationService) *StationHandler {
	return &StationHandler{service: service}
}

func (h *StationHandler) AddStation(c *gin.Context) {
	var input struct {
		Kod       string `json:"stan_kod"`
		Name      string `json:"stan_name"`
		Priznak   int    `json:"stan_priznak"`
		Paragraph string `json:"paragraph"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "stan_name is required"})
		return
	}

	station, err := h.service.AddStation(model.Station{
		Kod:       input.Kod,
		Name:      input.Name,
		Priznak:   input.Priznak,
		Paragraph: strings.TrimSpace(input.Paragraph),
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidStationData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid station data",
				"message": "Некорректный формат данных"})
			return
		case errors.Is(err, service.ErrDuplicateStation):
			c.JSON(http.StatusConflict, gin.H{"error": "station already exists",
				"message": "Станция уже существует"})
			return
		default:
			respondInternalError(c)
			return
		}
	}

	c.JSON(http.StatusCreated, station)
}

func (h *StationHandler) GetStations(c *gin.Context) {
	filters, err := parseFilters(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stations, metadata, err := h.service.GetStations(filters)
	if err != nil {
		respondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stations": stations,
		"metadata": metadata,
	})
}

func (h *StationHandler) GetStation(c *gin.Context) {
	kod := strings.TrimSpace(c.Param("kod"))
	if kod == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid station code"})
		return
	}

	station, err := h.service.GetStation(kod)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "station not found"})
			return
		default:
			respondInternalError(c)
			return
		}
	}

	c.JSON(http.StatusOK, station)
}

func (h *StationHandler) UpdateStation(c *gin.Context) {
	kod := strings.TrimSpace(c.Param("kod"))
	if kod == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid station code"})
		return
	}

	var input struct {
		Kod       string `json:"stan_kod"`
		Name      string `json:"stan_name"`
		Priznak   int    `json:"stan_priznak"`
		Paragraph string `json:"paragraph"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	station, err := h.service.UpdateStation(kod,
		model.Station{
			Kod:       input.Kod,
			Name:      input.Name,
			Priznak:   input.Priznak,
			Paragraph: strings.TrimSpace(input.Paragraph),
		})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidStationData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid station data"})
			return
		case errors.Is(err, service.ErrDuplicateStation):
			c.JSON(http.StatusConflict, gin.H{"error": "station already exists"})
			return
		case errors.Is(err, service.ErrStationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "station not found"})
			return
		default:
			respondInternalError(c)
			return
		}
	}

	c.JSON(http.StatusOK, station)
}

func (h *StationHandler) DeleteStation(c *gin.Context) {
	kod := strings.TrimSpace(c.Param("kod"))
	if kod == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid station code"})
		return
	}

	if err := h.service.DeleteStation(kod); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidStationData):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid station code"})
			return
		case errors.Is(err, service.ErrStationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "station not found"})
			return
		default:
			respondInternalError(c)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func parseFilters(qs url.Values) (model.Filters, error) {
	filters := model.Filters{
		Page:     1,
		PageSize: 20,
		Sort:     "id",
		SortSafeList: []string{
			"id", "-id",
			"stan_kod", "-stan_kod",
			"stan_name", "-stan_name",
			"paragraph", "-paragraph",
		},
	}

	if v := qs.Get("page"); v != "" {
		page, err := strconv.Atoi(v)
		if err != nil || page < 1 {
			return filters, ErrInvalidFilter("page must be a positive integer")
		}
		filters.Page = page
	}

	if v := qs.Get("page_size"); v != "" {
		size, err := strconv.Atoi(v)
		if err != nil || size < 1 || size > 200 {
			return filters, ErrInvalidFilter("page_size must be between 1 and 200")
		}
		filters.PageSize = size
	}

	if v := qs.Get("sort"); v != "" {
		filters.Sort = v
	}

	if v := qs.Get("name"); v != "" {
		filters.Name = strings.TrimSpace(v)
	}

	if v := qs.Get("paragraph"); v != "" {
		filters.Paragraph = strings.TrimSpace(v)
	}

	if !isSafeSort(filters.Sort, filters.SortSafeList) {
		return filters, ErrInvalidFilter("invalid sort value")
	}

	return filters, nil
}

func isSafeSort(sort string, safelist []string) bool {
	for _, safe := range safelist {
		if sort == safe {
			return true
		}
	}
	return false
}

type ErrInvalidFilter string

func (e ErrInvalidFilter) Error() string { return string(e) }

func respondInternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "internal_server_error",
		"message": "Внутренняя ошибка сервера",
	})
}
