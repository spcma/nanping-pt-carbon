package http

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// MethodologyHandler 方法学处理器
type MethodologyHandler struct {
	appService *application.MethodologyAppService
}

// NewMethodologyHandler creates methodology handler
func NewMethodologyHandler(appService *application.MethodologyAppService) *MethodologyHandler {
	return &MethodologyHandler{
		appService: appService,
	}
}

// Create creates a methodology
func (h *MethodologyHandler) Create(c *gin.Context) {
	var param application.CreateMethodologyParam
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.Warn("methodology", "create methodology - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := platform_http.GetCurrentUser(c)
	param.UserID = securityUser.ID

	id, err := h.appService.CreateMethodology(platform_http.Ctx(c), param)
	if err != nil {
		logger.Error("methodology", "create methodology failed",
			zap.String("name", param.Name),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("methodology", "methodology created successfully",
		zap.Int64("methodology_id", id),
		zap.String("name", param.Name),
	)
	response.Success(c, gin.H{"id": id})
}

// Update updates a methodology
func (h *MethodologyHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd application.UpdateMethodologyParam
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	user := platform_http.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.UpdateMethodology(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a methodology
func (h *MethodologyHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user := platform_http.GetCurrentUser(c)

	if err := h.appService.DeleteMethodology(platform_http.Ctx(c), id, user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID gets methodology by ID
func (h *MethodologyHandler) GetByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	methodology, err := h.appService.GetMethodologyByID(platform_http.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, methodology)
}

// GetByCode gets methodology by code
func (h *MethodologyHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	methodology, err := h.appService.GetMethodologyByCode(platform_http.Ctx(c), code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, methodology)
}

func (h *MethodologyHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, list)
}

// GetPage queries methodologies with pagination
func (h *MethodologyHandler) GetPage(c *gin.Context) {
	var query domain.MethodologyPageQuery
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

	methodologies, total, err := h.appService.GetMethodologyPage(platform_http.Ctx(c), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  methodologies,
		"total": total,
	})
}

// ChangeStatus changes methodology status
func (h *MethodologyHandler) ChangeStatus(c *gin.Context) {
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
	changeCmd := application.ChangeMethodologyStatusCommand{
		ID:     id,
		Status: domain.MethodologyStatus(cmd.Status),
		UserID: user.ID,
	}

	if err := h.appService.ChangeMethodologyStatus(platform_http.Ctx(c), changeCmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
