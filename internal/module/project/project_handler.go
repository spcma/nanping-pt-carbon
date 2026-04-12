package project

import (
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/entity"
	"app/internal/shared/logger"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	appService *ProjectAppService
}

// NewProjectHandler creates project handler
func NewProjectHandler(appService *ProjectAppService) *ProjectHandler {
	return &ProjectHandler{
		appService: appService,
	}
}

// Create creates a project
func (h *ProjectHandler) Create(c *gin.Context) {
	var cmd CreateProjectCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("project", "create project - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := platform_http.GetCurrentUser(c)
	cmd.UserID = securityUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
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
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd UpdateProjectCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	user := platform_http.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a project
func (h *ProjectHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)

	if err := h.appService.Delete(platform_http.Ctx(c), id, user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ProjectHandler) GetById(c *gin.Context) {
	var query ProjectQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if query.ID == 0 {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	project, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, project)
}

// GetByQuery 根据条件查询项目
func (h *ProjectHandler) GetByQuery(c *gin.Context) {
	var query ProjectQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	project, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, project)
}

type ProjectDetail struct {
	entity.BaseEntity
	Name        string        `json:"name" gorm:"column:name"`
	Code        string        `json:"code" gorm:"column:code"`
	Status      ProjectStatus `json:"status" gorm:"column:status"`
	Description string        `json:"description" gorm:"column:description"`
	Icon        string        `json:"icon" gorm:"column:icon"`
	StartDate   string        `json:"startDate" gorm:"column:start_date"`
	EndDate     string        `json:"endDate" gorm:"column:end_date"`
}

func (h *ProjectHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newList []*ProjectDetail
	for _, project := range list {
		newList = append(newList, &ProjectDetail{
			BaseEntity:  project.BaseEntity,
			Name:        project.Name,
			Code:        project.Code,
			Status:      project.Status,
			Description: project.Description,
			Icon:        project.Icon,
			StartDate:   project.StartDate.Format(time.DateOnly),
			EndDate:     project.EndDate.Format(time.DateOnly),
		})
	}

	response.Success(c, newList)
}

// GetPage queries projects with pagination
func (h *ProjectHandler) GetPage(c *gin.Context) {
	var query ProjectPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query.Fixed()

	projects, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, projects)
}

// ChangeStatus changes project status
func (h *ProjectHandler) ChangeStatus(c *gin.Context) {
	idStr := c.Query("id")
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
	changeCmd := ChangeProjectStatusCommand{
		ID:     id,
		Status: ProjectStatus(cmd.Status),
		UserID: user.ID,
	}

	if err := h.appService.ChangeStatus(platform_http.Ctx(c), changeCmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
