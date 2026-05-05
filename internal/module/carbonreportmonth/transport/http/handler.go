package transport

import (
	"app/internal/module/carbonreportmonth/application"
	"app/internal/module/carbonreportmonth/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
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
		response.BadRequest(c, "请求参数错误")
		return
	}

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "未授权访问")
		return
	}
	cmd.UserID = currentUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.RuntimeL.Error("create carbon report month failed",
			zap.String("collection_date", cmd.CollectionDate),
			zap.Error(err),
		)
		response.InternalError(c, "创建失败")
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
		response.BadRequest(c, "无效的ID")
		return
	}

	var cmd application.UpdateCarbonReportMonthCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	cmd.ID = id

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "未授权访问")
		return
	}
	cmd.UserID = currentUser.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("update carbon report month failed",
			zap.Int64("report_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除碳月报
func (h *CarbonReportMonthHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "未授权访问")
		return
	}

	if err := h.appService.Delete(platform_http.Ctx(c), id, currentUser.ID); err != nil {
		logger.RuntimeL.Error("delete carbon report month failed",
			zap.Int64("report_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetByID 根据ID获取碳月报
func (h *CarbonReportMonthHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	report, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get carbon report month failed",
			zap.Int64("report_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, report)
}

// GetPage 分页查询碳月报
func (h *CarbonReportMonthHandler) GetPage(c *gin.Context) {
	var query domain.CarbonReportMonthPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	query.Fixed()

	res, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get carbon report month page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, res)
}
