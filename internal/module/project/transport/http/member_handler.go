package transport

import (
	"app/internal/module/project/application"
	"app/internal/module/project/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// ProjectMembersHandler 项目成员处理器（接口层）
type ProjectMembersHandler struct {
	appService *application.ProjectMembersAppService
}

// NewProjectMembersHandler 创建项目成员处理器
func NewProjectMembersHandler(appService *application.ProjectMembersAppService) *ProjectMembersHandler {
	return &ProjectMembersHandler{
		appService: appService,
	}
}

// Create 创建项目成员
func (h *ProjectMembersHandler) Create(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.CreateProjectMemberCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("project", "create project member - invalid request",
			zap.String("error", err.Error()),
		)
		response.BadRequest(c, "请求参数错误")
		return
	}

	cmd.CreateBy = currentUser.ID

	id, err := h.appService.CreateProjectMember(platform_http.Ctx(c), &cmd)
	if err != nil {
		logger.Error("project", "create project member failed",
			zap.Int64("project_id", cmd.ProjectId),
			zap.Int64("user_id", cmd.UserId),
			zap.Error(err),
		)
		response.InternalError(c, "创建项目成员失败")
		return
	}

	logger.Info("project", "project member created successfully",
		zap.Int64("member_id", id),
		zap.Int64("project_id", cmd.ProjectId),
		zap.Int64("user_id", cmd.UserId),
	)
	response.Success(c, gin.H{"id": id})
}

// Update 更新项目成员
func (h *ProjectMembersHandler) Update(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.BadRequest(c, "ID不能为空")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var cmd application.UpdateProjectMemberCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	cmd.ID = id

	// 获取操作员用户信息
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}
	cmd.CreateBy = currentUser.ID

	if err := h.appService.UpdateProjectMember(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("update project member failed",
			zap.Int64("member_id", cmd.ID),
			zap.Error(err),
		)
		response.InternalError(c, "更新项目成员失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除项目成员
func (h *ProjectMembersHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.BadRequest(c, "ID不能为空")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	user := platform_http.GetCurrentUser(c)
	cmd := application.DeleteProjectMemberCommand{
		ID:       id,
		CreateBy: user.ID,
	}

	if err := h.appService.DeleteProjectMember(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("delete project member failed",
			zap.Int64("member_id", cmd.ID),
			zap.Error(err),
		)
		response.InternalError(c, "删除项目成员失败")
		return
	}

	response.Success(c, nil)
}

// GetList 获取项目成员列表
func (h *ProjectMembersHandler) GetList(c *gin.Context) {
	type GetByProjectIDRequest struct {
		ProjectID int64 `json:"projectId"`
		UserID    int64 `json:"userId"`
	}

	var query GetByProjectIDRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	var members []*domain.ProjectMembers
	var err error

	if query.ProjectID > 0 {
		members, err = h.appService.GetProjectMembersByProjectID(platform_http.Ctx(c), query.ProjectID)
		if err != nil {
			logger.RuntimeL.Error("get project members by project id failed",
				zap.Int64("project_id", query.ProjectID),
				zap.Error(err),
			)
			response.InternalError(c, "获取项目成员失败")
			return
		}
	}

	if query.UserID > 0 {
		members, err = h.appService.GetProjectMembersByUserID(platform_http.Ctx(c), query.UserID)
		if err != nil {
			logger.RuntimeL.Error("get project members by user id failed",
				zap.Int64("user_id", query.UserID),
				zap.Error(err),
			)
			response.InternalError(c, "获取项目成员失败")
			return
		}
	}

	response.Success(c, members)
}

// GetByUserID 根据用户 ID 获取项目成员
func (h *ProjectMembersHandler) GetByUserID(c *gin.Context) {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		response.BadRequest(c, "userId不能为空")
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的userId")
		return
	}

	members, err := h.appService.GetProjectMembersByUserID(platform_http.Ctx(c), userID)
	if err != nil {
		logger.RuntimeL.Error("get project members by user id failed",
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
		response.InternalError(c, "获取项目成员失败")
		return
	}

	response.Success(c, members)
}

// GetPage 分页查询项目成员
func (h *ProjectMembersHandler) GetPage(c *gin.Context) {
	var query domain.ProjectMembersPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	query.Fixed()

	members, err := h.appService.GetProjectMemberPage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get project member page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "分页查询项目成员失败")
		return
	}

	response.Success(c, members)
}
