@echo off
echo ====================================
echo   测试 NPFS HTTP API (DDD 架构)
echo ====================================
echo.

set BASE_URL=http://localhost:19870

echo [测试 1] 健康检查
curl -X GET "%BASE_URL%/health"
echo.
echo.

echo [测试 2] 检查目录是否存在
curl -X POST "%BASE_URL%/api/v1/dir/check" ^
     -H "Content-Type: application/json" ^
     -d "{\"path\":\"/test\"}"
echo.
echo.

echo [测试 3] 创建目录
curl -X POST "%BASE_URL%/api/v1/dir/create" ^
     -H "Content-Type: application/json" ^
     -d "{\"path\":\"/test/ddc_test\",\"recursive\":true}"
echo.
echo.

echo [测试 4] 列出目录
curl -X GET "%BASE_URL%/api/v1/dir/list?path=/test"
echo.
echo.

echo [测试 5] 保存文件内容
curl -X POST "%BASE_URL%/api/v1/file/save" ^
     -H "Content-Type: application/json" ^
     -d "{\"content\":\"Hello DDD World!\",\"dir\":\"/test\",\"filename\":\"hello.txt\"}"
echo.
echo.

echo [测试 6] 读取文件
curl -X GET "%BASE_URL%/api/v1/file/read?path=/test/hello.txt"
echo.
echo.

echo [测试 7] 删除文件
curl -X DELETE "%BASE_URL%/api/v1/file/delete?path=/test/hello.txt&force=true"
echo.
echo.

echo ====================================
echo   测试完成！
echo ====================================
pause
