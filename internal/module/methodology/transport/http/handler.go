package transport

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"strconv"
	"time"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// MethodologyHandler 方法学处理器（接口层）
type MethodologyHandler struct {
	appService *application.MethodologyAppService
}

// NewMethodologyHandler 创建方法学处理器
func NewMethodologyHandler(appService *application.MethodologyAppService) *MethodologyHandler {
	return &MethodologyHandler{
		appService: appService,
	}
}

// Create 创建方法学
func (h *MethodologyHandler) Create(c *gin.Context) {
	var cmd application.CreateMethodologyCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("methodology", "create methodology - invalid request",
			zap.String("error", err.Error()),
		)
		response.BadRequest(c, "请求参数错误")
		return
	}

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Forbidden(c, "user not found")
		return
	}
	cmd.UserID = currentUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.RuntimeL.Error("create methodology failed",
			zap.String("name", cmd.Name),
			zap.Error(err),
		)
		response.InternalError(c, "创建失败")
		return
	}

	logger.Info("methodology", "methodology created successfully",
		zap.Int64("methodology_id", id),
		zap.String("name", cmd.Name),
	)
	response.Success(c, gin.H{"id": id})
}

// Update 更新方法学
func (h *MethodologyHandler) Update(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var cmd application.UpdateMethodologyCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	cmd.ID = id

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Forbidden(c, "user not found")
		return
	}
	cmd.UserID = currentUser.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("update methodology failed",
			zap.Int64("methodology_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除方法学
func (h *MethodologyHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Forbidden(c, "user not found")
		return
	}

	if err := h.appService.Delete(platform_http.Ctx(c), id, currentUser.ID); err != nil {
		logger.RuntimeL.Error("delete methodology failed",
			zap.Int64("methodology_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetById 根据ID获取方法学
func (h *MethodologyHandler) GetById(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.BadRequest(c, "id is required")
		return
	}

	id := cast.ToInt64(idStr)
	if id == 0 {
		response.BadRequest(c, "invalid id")
		return
	}

	methodology, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get methodology failed",
			zap.Int64("methodology_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "获取失败")
		return
	}

	response.Success(c, methodology)
}

// GetByQuery 综合查询
func (h *MethodologyHandler) GetByQuery(c *gin.Context) {
	var query domain.MethodologyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	methodology, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get methodology by query failed",
			zap.Error(err),
		)
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, methodology)
}

// MethodologyDetail DTO（用于返回给前端）
type MethodologyDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	CreateBy    int64  `json:"createBy"`
	CreateTime  string `json:"createTime"`
}

// GetList 获取方法学列表
func (h *MethodologyHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		logger.RuntimeL.Error("get methodology list failed",
			zap.Error(err),
		)
		response.InternalError(c, "获取列表失败")
		return
	}

	var newList []*MethodologyDetail
	for _, methodology := range list {
		newList = append(newList, &MethodologyDetail{
			ID:          methodology.Id,
			Name:        methodology.Name,
			Code:        methodology.Code,
			Status:      string(methodology.Status),
			Description: methodology.Description,
			Icon:        methodology.Icon,
			StartDate:   methodology.StartDate.Format(time.DateOnly),
			EndDate:     methodology.EndDate.Format(time.DateOnly),
			CreateBy:    methodology.CreateBy,
			CreateTime:  methodology.CreateTime.Format(time.DateTime),
		})
	}

	response.Success(c, newList)
}

// GetPage 分页查询方法学
func (h *MethodologyHandler) GetPage(c *gin.Context) {
	var query domain.MethodologyPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	query.Fixed()

	result, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get methodology page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "分页查询失败")
		return
	}

	response.Success(c, result)
}

// ChangeStatus 变更方法学状态
func (h *MethodologyHandler) ChangeStatus(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var cmd struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Forbidden(c, "user not found")
		return
	}

	changeCmd := application.ChangeMethodologyStatusCommand{
		ID:     id,
		Status: domain.MethodologyStatus(cmd.Status),
		UserID: currentUser.ID,
	}

	if err := h.appService.ChangeStatus(platform_http.Ctx(c), changeCmd); err != nil {
		logger.RuntimeL.Error("change methodology status failed",
			zap.Int64("methodology_id", id),
			zap.String("status", cmd.Status),
			zap.Error(err),
		)
		response.InternalError(c, "变更状态失败")
		return
	}

	response.Success(c, nil)
}

// Activate 启用方法学
func (h *MethodologyHandler) Activate(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Forbidden(c, "user not found")
		return
	}

	if err := h.appService.Activate(platform_http.Ctx(c), id, currentUser.ID); err != nil {
		logger.RuntimeL.Error("activate methodology failed",
			zap.Int64("methodology_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "启用失败")
		return
	}

	response.Success(c, nil)
}

// Deactivate 禁用方法学
func (h *MethodologyHandler) Deactivate(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Forbidden(c, "user not found")
		return
	}

	if err := h.appService.Deactivate(platform_http.Ctx(c), id, currentUser.ID); err != nil {
		logger.RuntimeL.Error("deactivate methodology failed",
			zap.Int64("methodology_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "禁用失败")
		return
	}

	response.Success(c, nil)
}
