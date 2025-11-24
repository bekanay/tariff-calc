package handler

import (
	"net/http"
	"net/url"
	"strconv"
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

func (h *StationHandler) GetStations(c *gin.Context) {
	filters, err := parseFilters(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stations, metadata, err := h.service.GetStations(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stations": stations,
		"metadata": metadata,
	})
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
