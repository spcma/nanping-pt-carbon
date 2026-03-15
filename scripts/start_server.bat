@echo off
echo ================================================
echo   NPFS WebSocket 服务器启动脚本
echo ================================================
echo.
echo 正在启动 WebSocket 服务器...
echo 监听地址：ws://127.0.0.1:8080/ws
echo.
echo 提示：在另一个终端运行 client.bat 连接测试
echo 按 Ctrl+C 停止服务器
echo ================================================
echo.

go run . -mode websocket -addr :8080
