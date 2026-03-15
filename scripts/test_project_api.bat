@echo off
chcp 65001 >nul
echo =====================================
echo    Project Module API Tests
echo =====================================
echo.

REM 设置基础 URL 和 Token
set BASE_URL=http://localhost:8080/api
set TOKEN=your_jwt_token_here

echo [1] 创建项目
curl -X POST "%BASE_URL%/project" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -d "{\"name\":\"测试项目\",\"code\":\"TEST001\",\"description\":\"这是一个测试项目\"}"
echo.
echo.

echo [2] 获取项目列表
curl -X GET "%BASE_URL%/projects?pageNum=1&pageSize=10" ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.

echo [3] 根据 ID 获取项目
curl -X GET "%BASE_URL%/project/1" ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.

echo [4] 根据编码获取项目
curl -X GET "%BASE_URL%/project/code/TEST001" ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.

echo [5] 更新项目
curl -X PUT "%BASE_URL%/project/1" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -d "{\"name\":\"更新的测试项目\",\"description\":\"这是更新后的描述\"}"
echo.
echo.

echo [6] 变更项目状态
curl -X PUT "%BASE_URL%/project/1/status" ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -d "{\"status\":\"2\"}"
echo.
echo.

echo [7] 删除项目
curl -X DELETE "%BASE_URL%/project/1" ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.

echo =====================================
echo    测试完成
echo =====================================
pause
