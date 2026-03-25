package http

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/domain"
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

	securityUser := platform_http.GetCurrentUser(c)
	cmd.UserID = securityUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
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
	idStr := c.Query("id")
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

	user := platform_http.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a methodology
func (h *MethodologyHandler) Delete(c *gin.Context) {
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

// GetById gets methodology by ID
func (h *MethodologyHandler) GetById(c *gin.Context) {
	var query application.MethodologyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if query.ID == 0 {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	methodology, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, methodology)
}

// GetByQuery queries methodologies with conditions
func (h *MethodologyHandler) GetByQuery(c *gin.Context) {
	var query application.MethodologyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	methodology, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, methodology)
}

type MethodologyDetail struct {
	entity.BaseEntity
	Name        string                   `json:"name" gorm:"column:name"`
	Code        string                   `json:"code" gorm:"column:code"`
	Status      domain.MethodologyStatus `json:"status" gorm:"column:status"`
	Description string                   `json:"description" gorm:"column:description"`
	Icon        string                   `json:"icon" gorm:"column:icon"`
	StartDate   string                   `json:"startDate" gorm:"column:start_date"`
	EndDate     string                   `json:"endDate" gorm:"column:end_date"`
}

func (h *MethodologyHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var newList []*MethodologyDetail
	for _, methodology := range list {
		newList = append(newList, &MethodologyDetail{
			BaseEntity:  methodology.BaseEntity,
			Name:        methodology.Name,
			Code:        methodology.Code,
			Status:      methodology.Status,
			Description: methodology.Description,
			Icon:        methodology.Icon,
			StartDate:   methodology.StartDate.Format(time.DateOnly),
			EndDate:     methodology.EndDate.Format(time.DateOnly),
		})
	}

	response.Success(c, newList)
}

// GetPage queries methodologies with pagination
func (h *MethodologyHandler) GetPage(c *gin.Context) {
	var query domain.MethodologyPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query.Fixed()

	result, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ChangeStatus changes methodology status
func (h *MethodologyHandler) ChangeStatus(c *gin.Context) {
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
	changeCmd := application.ChangeMethodologyStatusCommand{
		ID:     id,
		Status: domain.MethodologyStatus(cmd.Status),
		UserID: user.ID,
	}

	if err := h.appService.ChangeStatus(platform_http.Ctx(c), changeCmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
