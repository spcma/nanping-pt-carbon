@echo off
set APP_NAME=nanping-pt-carbon
set REMOTE_PATH=root@115.236.162.150:/home/yanping/app/nanping/yanping-carbon/tmp
set MAIN_FILE=./cmd/main.go
set GOOS=linux
set GOARCH=amd64

echo Building %APP_NAME% for Linux...
go build -o %APP_NAME% %MAIN_FILE%
if errorlevel 1 (
    echo Build failed!
    exit /b 1
)

echo Uploading to remote server...
C:\Windows\System32\OpenSSH\scp.exe .\%APP_NAME% %REMOTE_PATH%
if errorlevel 1 (
    echo Upload failed!
    del /f /q .\%APP_NAME%
    exit /b 1
)

echo Cleaning up...
del /f /q .\%APP_NAME%
echo Deploy successful!
