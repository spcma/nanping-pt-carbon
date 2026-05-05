package transport

import (
	"app/internal/module/project/application"
	"app/internal/module/project/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// ProjectHandler 项目处理器（接口层）
type ProjectHandler struct {
	appService *application.ProjectAppService
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler(appService *application.ProjectAppService) *ProjectHandler {
	return &ProjectHandler{
		appService: appService,
	}
}

// Create 创建项目
func (h *ProjectHandler) Create(c *gin.Context) {

	securityUser := platform_http.GetCurrentUser(c)
	if securityUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.CreateProjectCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	cmd.UserID = securityUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.RuntimeL.Error("create project failed",
			zap.String("name", cmd.Name),
			zap.Error(err),
		)

		response.InternalError(c, "")
		return
	}

	logger.Info("project", "project created successfully",
		zap.Int64("project_id", id),
		zap.String("name", cmd.Name),
	)

	response.Success(c, nil)
}

// Update 更新项目
func (h *ProjectHandler) Update(c *gin.Context) {
	user := platform_http.GetCurrentUser(c)
	if user == nil {
		response.Unauthorized(c, "user is nil")
		return
	}

	var cmd application.UpdateProjectCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	cmd.UserID = user.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("update project failed",
			zap.Int64("project_id", cmd.ID),
			zap.Error(err),
		)

		response.InternalError(c, "")
		return
	}

	response.Success(c, nil)
}

// Delete 删除项目
func (h *ProjectHandler) Delete(c *gin.Context) {
	user := platform_http.GetCurrentUser(c)
	if user == nil {
		response.Unauthorized(c, "")
		return
	}

	var params map[string]any
	if err := c.ShouldBindJSON(&params); err != nil {
		response.BadRequest(c, "")
		return
	}

	var id int64
	if v, ok := params["id"].(int64); !ok {
		response.BadRequest(c, "invalid id")
		return
	} else {
		id = v
	}

	if err := h.appService.Delete(platform_http.Ctx(c), id, user.ID); err != nil {
		logger.RuntimeL.Error("delete project failed",
			zap.Int64("project_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "")
		return
	}

	response.Success(c, nil)
}

// GetByID 根据ID获取项目
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	project, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get project failed",
			zap.Int64("project_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, project)
}

// GetByQuery 根据条件查询项目
func (h *ProjectHandler) GetByQuery(c *gin.Context) {
	var query domain.ProjectQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	project, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get project failed",
			zap.Int64("project_id", query.ID),
			zap.Error(err),
		)
		response.InternalError(c, "")
		return
	}

	response.Success(c, project)
}

// ProjectDetail DTO（用于返回给前端）
type ProjectDetail struct {
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

// GetList 获取项目列表
func (h *ProjectHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		response.InternalError(c, "获取项目列表失败")
		return
	}

	var newList []*ProjectDetail
	for _, project := range list {
		newList = append(newList, &ProjectDetail{
			ID:          project.Id,
			Name:        project.Name,
			Code:        project.Code,
			Status:      string(project.Status),
			Description: project.Description,
			Icon:        project.Icon,
			StartDate:   project.StartDate.Format(time.DateOnly),
			EndDate:     project.EndDate.Format(time.DateOnly),
			CreateBy:    project.CreateBy,
			CreateTime:  project.CreateTime.Format(time.DateTime),
		})
	}

	response.Success(c, newList)
}

// GetPage 分页查询项目
func (h *ProjectHandler) GetPage(c *gin.Context) {
	var query domain.ProjectPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "查询参数格式错误")
		return
	}

	query.Fixed()

	projects, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.InternalError(c, "分页查询项目失败")
		return
	}

	response.Success(c, projects)
}

// ChangeStatus 变更项目状态
func (h *ProjectHandler) ChangeStatus(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "项目ID格式错误")
		return
	}

	var cmd struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数格式错误")
		return
	}

	user := platform_http.GetCurrentUser(c)
	changeCmd := application.ChangeProjectStatusCommand{
		ID:     id,
		Status: domain.ProjectStatus(cmd.Status),
		UserID: user.ID,
	}

	if err := h.appService.ChangeStatus(platform_http.Ctx(c), changeCmd); err != nil {
		response.InternalError(c, "变更项目状态失败")
		return
	}

	response.Success(c, nil)
}

// Activate 激活项目
func (h *ProjectHandler) Activate(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "项目ID格式错误")
		return
	}

	user := platform_http.GetCurrentUser(c)

	if err := h.appService.ActivateProject(platform_http.Ctx(c), id, user.ID); err != nil {
		response.InternalError(c, "激活项目失败")
		return
	}

	response.Success(c, nil)
}

// Complete 完成项目
func (h *ProjectHandler) Complete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "项目ID格式错误")
		return
	}

	user := platform_http.GetCurrentUser(c)

	if err := h.appService.CompleteProject(platform_http.Ctx(c), id, user.ID); err != nil {
		response.InternalError(c, "完成项目失败")
		return
	}

	response.Success(c, nil)
}

// Cancel 取消项目
func (h *ProjectHandler) Cancel(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "项目ID格式错误")
		return
	}

	user := platform_http.GetCurrentUser(c)

	if err := h.appService.CancelProject(platform_http.Ctx(c), id, user.ID); err != nil {
		response.InternalError(c, "取消项目失败")
		return
	}

	response.Success(c, nil)
}
