@echo off
echo ================================================
echo   NPFS HTTP API 快速测试
echo ================================================
echo.

REM 设置 API 地址
set API_URL=http://localhost:8080

echo 正在检查服务是否运行...
curl -s %API_URL%/health >nul 2>&1
if errorlevel 1 (
    echo [错误] 服务未启动！请先运行：go run . -http-server
    pause
    exit /b 1
)

echo [成功] 服务正常运行
echo.

echo ================================================
echo 测试 1: 创建目录
echo ================================================
echo.
curl -X POST %API_URL%/api/v1/dir/create ^
  -H "Content-Type: application/json" ^
  -d "{\"path\":\"/test_http\",\"recursive\":true}"
echo.
echo.

echo ================================================
echo 测试 2: 检查目录
echo ================================================
echo.
curl -X POST %API_URL%/api/v1/dir/check ^
  -d "path=/test_http"
echo.
echo.

echo ================================================
echo 测试 3: 保存文件内容
echo ================================================
echo.
curl -X POST %API_URL%/api/v1/file/save ^
  -H "Content-Type: application/json" ^
  -d "{\"content\":\"Hello from HTTP API!\",\"dir\":\"/test_http\",\"filename\":\"hello.txt\"}"
echo.
echo.

echo ================================================
echo 测试 4: 读取文件
echo ================================================
echo.
curl %API_URL%/api/v1/file/read?path=/test_http/hello.txt
echo.
echo.

echo ================================================
echo 测试 5: 列出目录
echo ================================================
echo.
curl %API_URL%/api/v1/dir/list?path=/test_http
echo.
echo.

echo ================================================
echo 测试 6: 删除文件
echo ================================================
echo.
curl -X DELETE "%API_URL%/api/v1/file/delete?path=/test_http/hello.txt&force=true"
echo.
echo.

echo ================================================
echo 测试完成!
echo ================================================
echo.
echo 提示：
echo - 打开 np_fs_web.html 使用图形界面
echo - 查看 NPFS_HTTP_API.md 了解完整 API 文档
echo.

pause
