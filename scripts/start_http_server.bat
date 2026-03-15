@echo off
echo ================================================
echo   NPFS HTTP 服务器启动脚本
echo ================================================
echo.

REM 检查是否已安装 gin
go mod tidy

echo.
echo 正在启动 NPFS HTTP 服务器...
echo 服务将监听在 http://localhost:8080
echo.

REM 运行 HTTP 服务器示例
go run . -http-server

pause
