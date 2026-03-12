@echo off
chcp 65001 >nul
cls
echo.
echo ================================================
echo   NPFS WebSocket - 功能演示
echo ================================================
echo.
echo 本系统提供两种运行模式：
echo.
echo   1. 交互式命令行模式（本地使用）
echo      - 直接在终端输入命令操作 NPFS
echo      - 适合本地测试和快速操作
echo.
echo   2. WebSocket 服务器模式（远程调用）
echo      - 启动 WebSocket 服务
echo      - 支持多个客户端同时连接
echo      - 适合远程操作和集成到其他系统
echo.
echo ================================================
echo.
echo 可用命令示例：
echo.
echo   基础操作:
echo     help                     - 查看帮助
echo     init                     - 初始化 NPFS 服务
echo     close                    - 关闭服务
echo.
echo   目录管理:
echo     checkdir /test           - 检查目录是否存在
echo     createdir /myfiles       - 创建目录
echo     ls /myfiles              - 列出目录内容
echo     delete /myfiles          - 删除目录
echo.
echo   文件操作:
echo     upload C:\test.txt /files test.txt   - 上传文件
echo     readfile /files/test.txt             - 读取文件
echo     savefile /files/test.txt C:\out.txt  - 保存到本地
echo     savecontent "Hello" /files hi.txt    - 保存文本
echo.
echo   测试工具:
echo     batch 100                - 批量创建 100 个文件
echo     test                     - 运行完整测试
echo.
echo ================================================
echo.
echo 请选择启动方式：
echo.
echo   [1] 交互式命令行模式
echo   [2] WebSocket 服务器模式
echo   [3] 查看完整文档
echo   [0] 退出
echo.
set /p CHOICE="请输入选项 (0-3): "

if "%CHOICE%"=="1" goto INTERACTIVE
if "%CHOICE%"=="2" goto WEBSOCKET
if "%CHOICE%"=="3" goto DOCS
if "%CHOICE%"=="0" goto END

echo 无效选项
pause
goto END

:INTERACTIVE
echo.
echo 启动交互式命令行...
echo.
go run . -mode interactive
goto END

:WEBSOCKET
echo.
echo 启动 WebSocket 服务器...
echo 监听地址：ws://127.0.0.1:8080/ws
echo.
echo 提示：在另一个终端窗口运行 client.bat 连接
echo.
go run . -mode websocket -addr :8080
goto END

:DOCS
echo.
start notepad USAGE.md
goto MENU

:MENU
cls
goto START

:END
echo.
echo 再见！
timeout /t 2 >nul
