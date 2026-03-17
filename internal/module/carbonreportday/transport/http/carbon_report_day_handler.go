package http

import (
	"app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportday/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CarbonReportDayHandler 碳报告日报处理器
type CarbonReportDayHandler struct {
	appService *application.CarbonReportDayAppService
}

// NewCarbonReportDayHandler creates carbon report day handler
func NewCarbonReportDayHandler(appService *application.CarbonReportDayAppService) *CarbonReportDayHandler {
	return &CarbonReportDayHandler{
		appService: appService,
	}
}

// Create creates a carbon report day
func (h *CarbonReportDayHandler) Create(c *gin.Context) {
	var cmd application.CreateCarbonReportDayCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("carbon_report_day", "create carbon report day - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := platform_http.GetCurrentUser(c)
	cmd.UserID = securityUser.ID

	id, err := h.appService.CreateCarbonReportDay(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.Error("carbon_report_day", "create carbon report day failed",
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("carbon_report_day", "carbon report day created successfully",
		zap.Int64("report_id", id),
	)
	response.Success(c, gin.H{"id": id})
}

// Update updates a carbon report day
func (h *CarbonReportDayHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd application.UpdateCarbonReportDayCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	user := platform_http.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.UpdateCarbonReportDay(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a carbon report day
func (h *CarbonReportDayHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)

	if err := h.appService.DeleteCarbonReportDay(platform_http.Ctx(c), id, user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID gets carbon report day by ID
func (h *CarbonReportDayHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	report, err := h.appService.GetCarbonReportDayByID(platform_http.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, report)
}

// GetPage queries carbon report days with pagination
func (h *CarbonReportDayHandler) GetPage(c *gin.Context) {
	var query domain.CarbonReportDayPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if query.PageNum == 0 {
		query.PageNum = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	reports, total, err := h.appService.GetCarbonReportDayPage(platform_http.Ctx(c), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  reports,
		"total": total,
	})
}
