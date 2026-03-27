package logger

// 这是一个示例文件，展示如何使用 Sugar Logger

/*
使用示例：

1. 基础用法（链式调用）：
   logger := NewSugarLogger("business")
   logger.WithUser(userID, username).
          WithTraceID(traceID).
          Info("UserID performed action", zap.String("action", "login"))

2. 从上下文获取 logger：
   func handler(c *gin.Context) {
       ctxLogger := GetLoggerFromContext(c)
       ctxLogger.Info("Processing request")
   }

3. 业务事件记录：
   logger.BusinessEvent("user_registration",
       zap.Int64("user_id", userID),
       zap.String("email", email),
   )

4. 数据库操作记录：
   logger.DatabaseOperation("SELECT", "users", duration,
       zap.Int64("rows_affected", count),
   )

5. 外部服务调用：
   logger.ExternalCall("payment_gateway", "/api/charge", duration,
       zap.String("order_id", orderID),
   )

6. 使用预定义的全局 logger：
   logger.RuntimeLogger.Info("Application started")
   logger.ErrorLogger.Error("Something went wrong", zap.Error(err))

7. 使用便捷字段函数：
   logger.Info("Request processed",
       WithUserID(userID),
       WithPath(c.Request.URL.Path),
       WithStatusCode(200),
       WithLatency(latency.Milliseconds()),
   )
*/
