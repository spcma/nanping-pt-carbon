package transport

import (
	"app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportday/domain"
	ipfs_application "app/internal/module/ipfs/application"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"app/internal/shared/timeutil"
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CarbonReportDayHandler 碳日报处理器（接口层）
type CarbonReportDayHandler struct {
	appService *application.CarbonReportDayService
}

// NewCarbonReportDayHandler 创建碳日报处理器
func NewCarbonReportDayHandler(appService *application.CarbonReportDayService) *CarbonReportDayHandler {
	return &CarbonReportDayHandler{
		appService: appService,
	}
}

// Create 创建碳日报
func (h *CarbonReportDayHandler) Create(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.CreateCarbonReportDayCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	cmd.UserID = currentUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), &cmd)
	if err != nil {
		logger.RuntimeL.Error("create carbon_report_day failed",
			zap.Error(err),
		)
		response.InternalError(c, "创建失败")
		return
	}

	logger.RuntimeL.Info("create carbon_report_day success",
		zap.Int64("id", id),
		zap.Int64("create_by", cmd.UserID),
	)

	response.Success(c, gin.H{"id": id})
}

// Update 更新碳日报
func (h *CarbonReportDayHandler) Update(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.UpdateCarbonReportDayCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	cmd.UserID = currentUser.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("update carbon_report_day failed",
			zap.Int64("report_id", cmd.ID),
			zap.Int64("update_by", cmd.UserID),
			zap.Error(err),
		)
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除碳日报
func (h *CarbonReportDayHandler) Delete(c *gin.Context) {
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
		logger.RuntimeL.Error("delete carbon_report_day failed",
			zap.Error(err))
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetByID 根据ID获取碳日报
func (h *CarbonReportDayHandler) GetByID(c *gin.Context) {

	idstr := c.Query("id")
	id := cast.ToInt64(idstr)
	if id == 0 {
		response.BadRequest(c, "invalid id")
		return
	}

	report, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get carbon_report_day failed", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, report)
}

// GetPage 分页查询碳日报
func (h *CarbonReportDayHandler) GetPage(c *gin.Context) {
	var query domain.CarbonReportDayLatestPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "")
		return
	}

	query.Fixed()

	res, err := h.appService.GetLatestByDatePage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get carbon report day page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, res)
}

// GetLatestByDatePage 按日期分组查询每天最新的记录
func (h *CarbonReportDayHandler) GetLatestByDatePage(c *gin.Context) {
	var query domain.CarbonReportDayLatestPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "")
		return
	}

	query.Fixed()

	res, err := h.appService.GetLatestByDatePage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get latest carbon report day by date page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, res)
}

func (h *CarbonReportDayHandler) ReportDay(c *gin.Context) {
	type calcDirRequest struct {
		RootDir    string `json:"rootDir" form:"rootDir"` // 要扫描的根目录，如 "/aibk/26/03/27"
		Date       string `json:"date" form:"date"`       // 日期，格式 "2026-03-27"
		ClientName string `json:"clientName" form:"clientName"`
	}

	var dto calcDirRequest
	if err := c.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(c, "")
		return
	}

	if dto.RootDir == "" {
		response.BadRequest(c, "请指定目录")
		return
	}

	if dto.Date == "" {
		response.BadRequest(c, "请指定日期")
		return
	}

	go func() {
		ctx := context.Background()

		ipfsService := ipfs_application.Ipfs()

		report, err := ipfsService.CalcDir(ctx, dto.ClientName, dto.RootDir, dto.Date)
		if err != nil {
			logger.IpfsL.Error("calcDir failed",
				zap.String("rootDir", dto.RootDir),
				zap.String("date", dto.Date),
				zap.Error(err))
			return
		}

		cmd := &application.CreateCarbonReportDayCommand{}

		if val, ok := report["turnover"].(float64); ok {
			cmd.Turnover = val
		}
		if val, ok := report["baseline"].(float64); ok {
			cmd.Baseline = val
		}
		if val, ok := report["carbonReduce"].(float64); ok {
			cmd.CarbonReduction = val
		}
		if val, ok := report["hash"].(string); ok {
			cmd.Hash = val
		}
		if val, ok := report["traceCode"].(string); ok {
			cmd.TraceCode = val
		}
		if val, ok := report["collectionDate"].(timeutil.Time); ok {
			cmd.CollectionDate = val
		}

		_, err = h.appService.Create(ctx, cmd)
		if err != nil {
			logger.SchedulerL.Error("碳日报创建失败", zap.Error(err))
			return
		}

		logger.IpfsL.Info("calcDir completed",
			zap.String("rootDir", dto.RootDir),
			zap.String("date", dto.Date),
			zap.Any("turnover", report["turnover"]))
	}()

	response.Success(c, "计算任务已启动，请稍后查看结果")
}
