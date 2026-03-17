@echo off

rem 本地应用程序名
set "local_file=npc"

rem 上传到服务器的目录
set "target_path=root@192.168.1.10:/yanping/app/nanping/yanping-carbon/tmp"

if /i "%~1"=="1" call :1
if /i "%~1"=="build" call :build
if /i "%~1"=="clear" call :clear
if /i "%~1"=="rscp" call :rscp
if /i "%~1"=="cleargen" call :cleargen
if /i "%~1"=="cc" call :clearall

exit /b

rem -------------------------------------------------------------
rem 打包并发布到线上环境
:1

go env -w GOOS=linux
go env -w GOARCH=amd64
go env -w CGO_ENABLED=0

if exist "%local_file%" (
    echo -- del %local_file%
    del "%local_file%"
    echo -- del complete
)

echo -- go build -o %local_file%

go build -o %local_file%

echo -- build complete

echo scp "C:\Windows\System32\OpenSSH\scp.exe" %local_file% %target_path%

"C:\Windows\System32\OpenSSH\scp.exe" %local_file% %target_path%

echo -- scp complete

if exist "%local_file%" (
    echo -- del %local_file%
    del "%local_file%"
    echo -- del complete
)

echo current time ~ %time%

exit /b

rem -------------------------------------------------------------
rem 打包项目
:build

go env -w GOOS=linux
go env -w GOARCH=amd64
go env -w CGO_ENABLED=0

if exist "%local_file%" (
    echo -- del %local_file%
    del "%local_file%"
    echo -- del complete
) else (
    echo -- file not found : %local_file%
)

echo -- go build -o %local_file%

go build -o %local_file%

echo -- build complete

echo current time ~ %time%

exit /b

rem -------------------------------------------------------------
rem 清除打包文件
:clear

if exist "%local_file%" (
    echo -- del %local_file%
    del "%local_file%"
    echo -- del complete
) else (
    echo -- file not found : %local_file%
)

echo current time ~ %time%

exit /b

rem -------------------------------------------------------------
rem 将打包文件上传到服务器
:rscp

echo "C:\Windows\System32\OpenSSH\scp.exe" %local_file% %target_path%

"C:\Windows\System32\OpenSSH\scp.exe" %local_file% %target_path%

echo current time ~ %time%

exit /b

rem -------------------------------------------------------------
rem 清除生成的代码
:cleargen

rmdir /s /q "./gen/models"

echo current time ~ %time%

exit /b

rem -------------------------------------------------------------
rem 清除所有缓存
:clearall

go clean -cache