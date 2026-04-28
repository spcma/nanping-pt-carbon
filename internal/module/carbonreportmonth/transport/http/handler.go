package transport

import (
	"app/internal/module/carbonreportmonth/application"
	"app/internal/module/carbonreportmonth/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CarbonReportMonthHandler 碳月报处理器（接口层）
type CarbonReportMonthHandler struct {
	appService *application.CarbonReportMonthAppService
}

// NewCarbonReportMonthHandler 创建碳月报处理器
func NewCarbonReportMonthHandler(appService *application.CarbonReportMonthAppService) *CarbonReportMonthHandler {
	return &CarbonReportMonthHandler{
		appService: appService,
	}
}

// Create 创建碳月报
func (h *CarbonReportMonthHandler) Create(c *gin.Context) {
	var cmd application.CreateCarbonReportMonthCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("carbon_report_month", "create carbon report month - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := platform_http.GetCurrentUser(c)
	if securityUser == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	cmd.UserID = securityUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.Error("carbon_report_month", "create carbon report month failed",
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("carbon_report_month", "carbon report month created successfully",
		zap.Int64("report_id", id),
	)
	response.Success(c, gin.H{"id": id})
}

// Update 更新碳月报
func (h *CarbonReportMonthHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd application.UpdateCarbonReportMonthCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	user := platform_http.GetCurrentUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	cmd.UserID = user.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除碳月报
func (h *CarbonReportMonthHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.appService.Delete(platform_http.Ctx(c), id, user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 根据ID获取碳月报
func (h *CarbonReportMonthHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	report, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, report)
}

// GetPage 分页查询碳月报
func (h *CarbonReportMonthHandler) GetPage(c *gin.Context) {
	var query domain.CarbonReportMonthPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query.Fixed()

	res, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, res)
}
