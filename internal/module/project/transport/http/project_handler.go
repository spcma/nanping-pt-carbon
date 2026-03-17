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

// ProjectHandler 项目处理器
type ProjectHandler struct {
	appService *application.ProjectAppService
}

// NewProjectHandler creates project handler
func NewProjectHandler(appService *application.ProjectAppService) *ProjectHandler {
	return &ProjectHandler{
		appService: appService,
	}
}

// Create creates a project
func (h *ProjectHandler) Create(c *gin.Context) {
	var cmd application.CreateProjectCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("project", "create project - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := platform_http.GetCurrentUser(c)
	cmd.UserID = securityUser.ID

	id, err := h.appService.CreateProject(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.Error("project", "create project failed",
			zap.String("name", cmd.Name),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("project", "project created successfully",
		zap.Int64("project_id", id),
		zap.String("name", cmd.Name),
	)
	response.Success(c, gin.H{"id": id})
}

// Update updates a project
func (h *ProjectHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd application.UpdateProjectCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	user := platform_http.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.UpdateProject(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a project
func (h *ProjectHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)

	if err := h.appService.DeleteProject(platform_http.Ctx(c), id, user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID gets project by ID
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	project, err := h.appService.GetProjectByID(platform_http.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, project)
}

// GetByCode gets project by code
func (h *ProjectHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	project, err := h.appService.GetProjectByCode(platform_http.Ctx(c), code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, project)
}

// GetPage queries projects with pagination
func (h *ProjectHandler) GetPage(c *gin.Context) {
	var query domain.ProjectPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if query.PageNum == 0 {
		query.PageNum = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	projects, total, err := h.appService.GetProjectPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  projects,
		"total": total,
	})
}

// ChangeStatus changes project status
func (h *ProjectHandler) ChangeStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := platform_http.GetCurrentUser(c)
	changeCmd := application.ChangeProjectStatusCommand{
		ID:     id,
		Status: domain.ProjectStatus(cmd.Status),
		UserID: user.ID,
	}

	if err := h.appService.ChangeProjectStatus(platform_http.Ctx(c), changeCmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
