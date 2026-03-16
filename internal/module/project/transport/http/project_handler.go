package http

import (
	"app/internal/module/project/application"
	"app/internal/module/project/domain"
	http2 "app/internal/platform/http"
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

	securityUser := http2.GetCurrentUser(c)
	cmd.UserID = securityUser.ID

	id, err := h.appService.CreateProject(http2.Ctx(c), cmd)
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

	user := http2.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.UpdateProject(http2.Ctx(c), cmd); err != nil {
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

	user := http2.GetCurrentUser(c)

	if err := h.appService.DeleteProject(http2.Ctx(c), id, user.ID); err != nil {
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

	project, err := h.appService.GetProjectByID(http2.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, project)
}

// GetByCode gets project by code
func (h *ProjectHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	project, err := h.appService.GetProjectByCode(http2.Ctx(c), code)
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

	projects, total, err := h.appService.GetProjectPage(http2.Ctx(c), &query)
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

	user := http2.GetCurrentUser(c)
	changeCmd := application.ChangeProjectStatusCommand{
		ID:     id,
		Status: domain.ProjectStatus(cmd.Status),
		UserID: user.ID,
	}

	if err := h.appService.ChangeProjectStatus(http2.Ctx(c), changeCmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
