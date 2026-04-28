package scheduler

import (
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// routes 路由注册器
type routes struct {
	handler *SchedulerHandler
}

// NewSchedulerRoutes 创建定时任务模块的路由注册器
// 注意: 此函数会在 InitRouter() 中被调用,此时所有业务模块的 Service 已经初始化完成
// 调度任务的 Start() 会在所有路由注册完成后才调用,确保任务执行时依赖的 Service 已就绪
func NewSchedulerRoutes() shared_http.RouteRegistry {
	scheduler := Default()

	handler := NewSchedulerHandler(scheduler)

	return &routes{
		handler: handler,
	}
}

// RegisterRoutes 注册路由
func (r *routes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 定时任务管理路由 - /api/scheduler/*
		schedulerGroup := authTypeRequiredRoute.Group("/scheduler")
		{
			// 任务列表
			schedulerGroup.GET("/tasks", r.handler.ListTasks)

			// 单个任务操作
			schedulerGroup.POST("/task", r.handler.AddTask)
			schedulerGroup.GET("/tasks/:name", r.handler.GetTaskStatus)
			schedulerGroup.DELETE("/tasks/:name", r.handler.RemoveTask)
			schedulerGroup.PUT("/tasks/:name/enable", r.handler.EnableTask)
			schedulerGroup.PUT("/tasks/:name/disable", r.handler.DisableTask)
		}
	}
}
