@echo off
setlocal

echo 🏴‍☠️  AnonBOX Build Script 🏴‍☠️
echo =================================

REM Check for Go
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed! Please install Go from https://go.dev/dl/
    pause
    exit /b 1
)
echo ✅ Go is installed.

REM Check for GCC (Required for GUI)
gcc --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  GCC is not installed! GUI build will fail.
    echo    Please install TDM-GCC ^(https://jmeubank.github.io/tdm-gcc/^) or MinGW.
    echo    Proceeding with CLI build only...
    set BUILD_GUI=0
) else (
    echo ✅ GCC is installed.
    set BUILD_GUI=1
)

echo.
echo 📦 Building CLI...
go build -o anonbox-cli.exe ./cmd/cli
if %errorlevel% neq 0 (
    echo ❌ CLI Build Failed!
    exit /b 1
)
echo ✅ CLI Built: anonbox-cli.exe

if %BUILD_GUI%==1 (
    echo.
    echo 🎨 Building GUI...
    go build -o anonbox-gui.exe ./cmd/gui
    if %errorlevel% neq 0 (
        echo ❌ GUI Build Failed!
    ) else (
        echo ✅ GUI Built: anonbox-gui.exe
    )
)

echo.
echo 🎉 Build Complete!
pause
