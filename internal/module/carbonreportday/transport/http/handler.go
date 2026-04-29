package transport

import (
	"app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportday/domain"
	ipfs_application "app/internal/module/ipfs/application"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"context"
	"net/http"
	"strconv"

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
	var cmd application.CreateCarbonReportDayCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("carbon_report_day", "create carbon report day - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := platform_http.GetCurrentUser(c)
	if securityUser == nil {
		response.Forbidden(c, "用户信息异常")
		return
	}
	cmd.UserID = securityUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), &cmd)
	if err != nil {
		logger.Error("carbon_report_day", "create carbon report day failed",
			zap.Error(err),
		)
		response.InternalError(c, "创建失败")
		return
	}

	logger.Info("carbon_report_day", "carbon report day created successfully",
		zap.Int64("report_id", id),
	)
	response.Success(c, gin.H{"id": id})
}

// Update 更新碳日报
func (h *CarbonReportDayHandler) Update(c *gin.Context) {
	var cmd application.UpdateCarbonReportDayCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	user := platform_http.GetCurrentUser(c)
	if user == nil {
		response.Forbidden(c, "用户信息异常")
		return
	}
	cmd.UserID = user.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		response.InternalError(c, "更新失败")
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除碳日报
func (h *CarbonReportDayHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)
	if user == nil {
		response.Forbidden(c, "用户信息异常")
		return
	}

	if err := h.appService.Delete(platform_http.Ctx(c), id, user.ID); err != nil {
		logger.RuntimeL.Error("delete carbon report day failed", zap.Error(err))
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetByID 根据ID获取碳日报
func (h *CarbonReportDayHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	report, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get carbon report day failed", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, report)
}

// GetPage 分页查询碳日报
func (h *CarbonReportDayHandler) GetPage(c *gin.Context) {
	var query domain.CarbonReportDayPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	query.Fixed()

	res, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get carbon report day page failed", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, res)
}

func (h *CarbonReportDayHandler) ReportDay(c *gin.Context) {
	type CalcDirDto struct {
		RootDir    string `json:"rootDir" form:"rootDir"` // 要扫描的根目录，如 "/aibk/26/03/27"
		Date       string `json:"date" form:"date"`       // 日期，格式 "2026-03-27"
		ClientName string `json:"clientName" form:"clientName"`
	}

	var dto CalcDirDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(c, err.Error())
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

		turnover, err := ipfsService.CalcDir(ctx, dto.ClientName, dto.RootDir, dto.Date)
		if err != nil {
			logger.IpfsL.Error("calcDir failed", zap.String("rootDir", dto.RootDir), zap.String("date", dto.Date), zap.Error(err))
			return
		}

		logger.IpfsL.Info("calcDir completed", zap.String("rootDir", dto.RootDir), zap.String("date", dto.Date), zap.Any("turnover", turnover["turnover"]))
	}()
}
