package reports

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"green-compass-backend/pkg/httpx"
)

type ReportTypeHandler struct{}

func NewReportTypeHandler() *ReportTypeHandler {
	return &ReportTypeHandler{}
}

func (h *ReportTypeHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	group := router.Group("/v1")
	group.Use(authMiddleware)
	group.GET("/report-types", h.ListReportTypes)
}

func (h *ReportTypeHandler) ListReportTypes(c *gin.Context) {
	reportTypes := []gin.H{
		{"key": "water_change", "label": "Water has changed"},
		{"key": "flooding_visible", "label": "Flooding is visible"},
		{"key": "unusually_dry", "label": "Unusually dry"},
		{"key": "vegetation_stress", "label": "Vegetation stress"},
		{"key": "heat_impact", "label": "Heat impact"},
		{"key": "infrastructure_change", "label": "Infrastructure change"},
		{"key": "incorrect_information", "label": "Incorrect information"},
		{"key": "other", "label": "Other"},
	}

	requestID := httpx.GetRequestID(c)
	c.JSON(http.StatusOK, httpx.Success(gin.H{"report_types": reportTypes}, requestID))
}
