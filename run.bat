@echo off
setlocal EnableDelayedExpansion

title AnonBOX Launcher
cls

echo 🏴‍☠️  AnonBOX Launcher
echo ======================

:: 1. Check for Python
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Error: Python is not installed or not in PATH.
    echo Please install Python 3.10+ from https://python.org/
    pause
    exit /b
)

:: 2. Create Virtual Environment
if not exist "venv" (
    echo 📦 Creating virtual environment...
    python -m venv venv
    if !errorlevel! neq 0 (
        echo ❌ Failed to create venv.
        pause
        exit /b
    )
)

:: 3. Check Dependencies
if not exist "venv\Lib\site-packages\customtkinter" (
    echo ⬇️  Installing dependencies...
    venv\Scripts\pip install -r requirements.txt
    if !errorlevel! neq 0 (
        echo ❌ Failed to install dependencies.
        pause
        exit /b
    )
    echo ✅ Dependencies installed.
)

:menu
cls
echo 🏴‍☠️  AnonBOX Launcher
echo ======================
echo.
echo [1] GUI Mode (Graphical Interface)
echo [2] CLI Mode (Command Line)
echo [3] Update System (Git Pull + Deps)
echo [4] Exit
echo.
set /p choice="Select option [1-4]: "

if "%choice%"=="1" (
    echo 🚀 Starting GUI...
    venv\Scripts\python main.py gui
    goto menu
)

if "%choice%"=="2" (
    echo 🚀 Starting CLI...
    set /p "pass=Enter Vault Password (optional, press Enter for none): "
    if "!pass!"=="" (
        venv\Scripts\python main.py cli
    ) else (
        venv\Scripts\python main.py cli --password "!pass!"
    )
    pause
    goto menu
)

if "%choice%"=="3" (
    echo 🔄 Updating System...
    git --version >nul 2>&1
    if %errorlevel% equ 0 (
        git pull
        if !errorlevel! equ 0 (
            echo ✅ Code updated.
            echo ⬇️  Updating dependencies...
            venv\Scripts\pip install -r requirements.txt --upgrade
            pause
        ) else (
            echo ❌ Git pull failed.
            pause
        )
    ) else (
        echo ❌ Git is not installed.
        pause
    )
    goto menu
)

if "%choice%"=="4" (
    exit /b
)

echo ❌ Invalid option.
pause
goto menu
