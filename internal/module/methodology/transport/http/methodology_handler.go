package http

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/domain"
	http2 "app/internal/platform/http"
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
	var cmd application.CreateMethodologyCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("methodology", "create methodology - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := http2.GetCurrentUser(c)
	cmd.UserID = securityUser.ID

	id, err := h.appService.CreateMethodology(http2.Ctx(c), cmd)
	if err != nil {
		logger.Error("methodology", "create methodology failed",
			zap.String("name", cmd.Name),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("methodology", "methodology created successfully",
		zap.Int64("methodology_id", id),
		zap.String("name", cmd.Name),
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

	var cmd application.UpdateMethodologyCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	user := http2.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.UpdateMethodology(http2.Ctx(c), cmd); err != nil {
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

	user := http2.GetCurrentUser(c)

	if err := h.appService.DeleteMethodology(http2.Ctx(c), id, user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID gets methodology by ID
func (h *MethodologyHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	methodology, err := h.appService.GetMethodologyByID(http2.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, methodology)
}

// GetByCode gets methodology by code
func (h *MethodologyHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	methodology, err := h.appService.GetMethodologyByCode(http2.Ctx(c), code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, methodology)
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

	methodologies, total, err := h.appService.GetMethodologyPage(http2.Ctx(c), query)
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

	user := http2.GetCurrentUser(c)
	changeCmd := application.ChangeMethodologyStatusCommand{
		ID:     id,
		Status: domain.MethodologyStatus(cmd.Status),
		UserID: user.ID,
	}

	if err := h.appService.ChangeMethodologyStatus(http2.Ctx(c), changeCmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
