package carbonreportday

import (
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
	service *CarbonReportDayService
}

// NewCarbonReportDayHandler creates carbon report day handler
func NewCarbonReportDayHandler(service *CarbonReportDayService) *CarbonReportDayHandler {
	return &CarbonReportDayHandler{
		service: service,
	}
}

// Create creates a carbon report day
func (h *CarbonReportDayHandler) Create(c *gin.Context) {
	var cmd CreateCarbonReportDayCommand
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

	id, err := h.service.CreateCarbonReportDay(platform_http.Ctx(c), cmd)
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

// Update updates a carbon report day
func (h *CarbonReportDayHandler) Update(c *gin.Context) {
	var cmd UpdateCarbonReportDayCommand
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

	if err := h.service.UpdateCarbonReportDay(platform_http.Ctx(c), cmd); err != nil {
		response.InternalError(c, "更新失败")
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
		response.BadRequest(c, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)
	if user == nil {
		response.Forbidden(c, "用户信息异常")
		return
	}

	if err := h.service.DeleteCarbonReportDay(platform_http.Ctx(c), id, user.ID); err != nil {
		logger.RuntimeL.Error("delete carbon report day failed", zap.Error(err))
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetByID gets carbon report day by ID
func (h *CarbonReportDayHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	report, err := h.service.GetCarbonReportDayByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get carbon report day failed", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, report)
}

// GetPage queries carbon report days with pagination
func (h *CarbonReportDayHandler) GetPage(c *gin.Context) {
	var query CarbonReportDayPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	query.Fixed()

	res, err := h.service.GetCarbonReportDayPage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get carbon report day page failed", zap.Error(err))
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, res)
}
