@echo off
echo ====================================
echo   启动 NPFS HTTP 服务器 (DDD 架构)
echo ====================================
echo.

REM 检查可执行文件是否存在
if not exist "nanping-pt-carbon.exe" (
    echo [错误] 未找到可执行文件，请先运行构建
    echo 运行：go build -o nanping-pt-carbon.exe
    pause
    exit /b 1
)

echo [信息] 正在启动 HTTP 服务器...
echo.

REM 启动服务器
nanping-pt-carbon.exe -mode http -addr :19870

pause
