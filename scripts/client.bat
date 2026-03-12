@echo off
echo ================================================
echo   NPFS WebSocket 客户端测试
echo ================================================
echo.

set /p ADDR="请输入 WebSocket 服务器地址 (默认：ws://127.0.0.1:8080/ws): "
if "%ADDR%"=="" set ADDR=ws://127.0.0.1:8080/ws

echo.
echo 连接到：%ADDR%
echo.

go run client_example.go %ADDR%
