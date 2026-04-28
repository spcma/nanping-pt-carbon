package transport

import (
	"app/internal/module/project/application"
	"app/internal/module/project/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"net/http"
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
	var cmd application.CreateProjectMemberCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("project", "create project member - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取操作员用户信息
	user := platform_http.GetCurrentUser(c)
	cmd.CreateBy = user.ID

	id, err := h.appService.CreateProjectMember(platform_http.Ctx(c), &cmd)
	if err != nil {
		logger.Error("project", "create project member failed",
			zap.Int64("project_id", cmd.ProjectId),
			zap.Int64("user_id", cmd.UserId),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
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
		response.Error(c, http.StatusBadRequest, "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd application.UpdateProjectMemberCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	// 获取操作员用户信息
	user := platform_http.GetCurrentUser(c)
	cmd.CreateBy = user.ID

	if err := h.appService.UpdateProjectMember(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除项目成员
func (h *ProjectMembersHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)
	cmd := application.DeleteProjectMemberCommand{
		ID:       id,
		CreateBy: user.ID,
	}

	if err := h.appService.DeleteProjectMember(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var members []*domain.ProjectMembers
	var err error

	if query.ProjectID > 0 {
		members, err = h.appService.GetProjectMembersByProjectID(platform_http.Ctx(c), query.ProjectID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if query.UserID > 0 {
		members, err = h.appService.GetProjectMembersByUserID(platform_http.Ctx(c), query.UserID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.Success(c, members)
}

// GetByUserID 根据用户 ID 获取项目成员
func (h *ProjectMembersHandler) GetByUserID(c *gin.Context) {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, "userId is required")
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid userId")
		return
	}

	members, err := h.appService.GetProjectMembersByUserID(platform_http.Ctx(c), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, members)
}

// GetPage 分页查询项目成员
func (h *ProjectMembersHandler) GetPage(c *gin.Context) {
	var query domain.ProjectMembersPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query.Fixed()

	members, err := h.appService.GetProjectMemberPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, members)
}
