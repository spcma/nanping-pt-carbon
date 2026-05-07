package transport

import (
	"app/internal/module/carbonreportmonth/application"
	"app/internal/module/carbonreportmonth/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"

	"github.com/spf13/cast"
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

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.CreateCarbonReportMonthCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	cmd.UserID = currentUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.RuntimeL.Error("create carbon_report_month failed",
			zap.String("collection_date", cmd.CollectionDate),
			zap.Error(err),
		)
		response.InternalError(c, "创建失败")
		return
	}

	logger.RuntimeL.Info("create carbon_report_month success",
		zap.String("collection_date", cmd.CollectionDate),
		zap.Int64("report_id", id),
	)

	response.Success(c, gin.H{"id": id})
}

// Update 更新碳月报
func (h *CarbonReportMonthHandler) Update(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.UpdateCarbonReportMonthCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	cmd.UserID = currentUser.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("update carbon_report_month failed",
			zap.Int64("report_id", cmd.ID),
			zap.Int64("update_by", cmd.UserID),
			zap.Error(err),
		)
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除碳月报
func (h *CarbonReportMonthHandler) Delete(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	type deleteRequest struct {
		ID int64 `json:"id"`
	}

	var req deleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "")
		return
	}

	if err := h.appService.Delete(platform_http.Ctx(c), req.ID, currentUser.ID); err != nil {
		logger.RuntimeL.Error("delete carbon_report_month failed",
			zap.Int64("report_id", req.ID),
			zap.Int64("delete_by", currentUser.ID),
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
	id := cast.ToInt64(idStr)
	if id == 0 {
		response.BadRequest(c, "invalid id")
		return
	}

	report, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get carbon_report_month failed",
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
		response.BadRequest(c, "")
		return
	}

	query.Fixed()

	res, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get carbon_report_month page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, res)
}
