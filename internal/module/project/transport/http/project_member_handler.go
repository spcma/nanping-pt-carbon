package http

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

// ProjectMembersHandler 项目成员处理器
type ProjectMembersHandler struct {
	appService *application.ProjectMembersService
}

// NewProjectMembersHandler 创建项目成员处理器
func NewProjectMembersHandler(appService *application.ProjectMembersService) *ProjectMembersHandler {
	return &ProjectMembersHandler{
		appService: appService,
	}
}

// Create creates a project member
func (h *ProjectMembersHandler) Create(c *gin.Context) {
	var param application.CreateProjectMemberParam
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Warn("project", "create project member - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	//	获取操作员用户信息
	user := platform_http.GetCurrentUser(c)
	param.CreateBy = user.ID

	id, err := h.appService.CreateProjectMember(platform_http.Ctx(c), param)
	if err != nil {
		logger.Error("project", "create project member failed",
			zap.Int64("project_id", param.ProjectId),
			zap.Int64("user_id", param.UserId),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("project", "project member created successfully",
		zap.Int64("member_id", id),
		zap.Int64("project_id", param.ProjectId),
		zap.Int64("user_id", param.UserId),
	)
	response.Success(c, gin.H{"id": id})
}

// Update updates a project member
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

	var param application.UpdateProjectMemberParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	param.ID = id

	//	获取操作员用户信息
	user := platform_http.GetCurrentUser(c)
	param.CreateBy = user.ID

	if err := h.appService.UpdateProjectMember(platform_http.Ctx(c), param); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a project member
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
	param := application.DeleteProjectMemberParam{
		ID:       id,
		CreateBy: user.ID,
	}

	if err := h.appService.DeleteProjectMember(platform_http.Ctx(c), param); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

type GetByProjectIDRequest struct {
	ProjectID int64 `json:"projectId"`
	UserID    int64 `json:"userId"`
}

// GetList gets project members by project ID
func (h *ProjectMembersHandler) GetList(c *gin.Context) {
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

// GetByUserID gets project members by user ID
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

// GetPage queries project members with pagination
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
