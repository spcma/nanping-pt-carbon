@echo off
chcp 65001 >nul
echo.
echo ================================================
echo   测试 NPFS WebSocket 功能
echo ================================================
echo.
echo 运行参数测试...
echo.

go run main.go np_fs.go np_fs_refactored.go np_fs_cli.go -h

echo.
echo ================================================
echo 按任意键退出...
pause >nul
