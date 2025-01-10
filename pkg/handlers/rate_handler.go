package handlers

import (
	"net/http"
	"time"

	"github.com/PavelBradnitski/Rates/pkg/services"
	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	service services.RateServiceInterface
}

func NewRateHandler(service *services.RateService) *RateHandler {
	return &RateHandler{service: service}
}

func (h *RateHandler) RegisterRoutes(router *gin.Engine) {
	rateGroup := router.Group("/rate")
	{
		rateGroup.GET("/", h.GetAllRates)
		rateGroup.GET("/:date", h.GetRateByDate)
	}
}

func (h *RateHandler) GetAllRates(c *gin.Context) {
	rates, err := h.service.GetAllRates(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rates"})
		return
	}
	if len(rates) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rates not found"})
		return
	}
	c.JSON(http.StatusOK, rates)
}

func (h *RateHandler) GetRateByDate(c *gin.Context) {
	date := c.Param("date")
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect date format. Expect YYYY-MM-DD"})
		return
	}
	foundRate, err := h.service.GetRateByDate(c, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rates not found"})
		return
	}
	if len(foundRate) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rates not found"})
		return
	}
	c.JSON(http.StatusOK, foundRate)
}
