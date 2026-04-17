@echo off
setlocal enabledelayedexpansion

if not exist ".env" (
    echo Error: .env file not found!
    exit /b 1
)

for /f "tokens=1,2 delims==" %%a in (.env) do (
    set "key=%%a"
    set "value=%%b"
    set "key=!key: =!"
    set "value=!value: =!"
    set "value=!value:"=!"
    
    if "!key!"=="APP_NAME" (
        set "APP_NAME=!value!"
    )
)

if "!APP_NAME!"=="" (
    echo Error: APP_NAME not found in .env file!
    exit /b 1
)

echo Building application: !APP_NAME!

go build -o !APP_NAME! ./cmd/main.go && echo Build successful: !APP_NAME! || (echo Build failed! & exit /b 1)

endlocal
