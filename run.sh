#!/bin/bash

# Configuration
VENV_DIR="venv"
REQUIREMENTS="requirements.txt"

# ANSI Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${CYAN}🏴‍☠️  AnonBOX Launcher${NC}"
echo "======================"

# 1. Check for Python 3
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}❌ Error: Python 3 is not installed or not in PATH.${NC}"
    echo "Please install Python 3.10+."
    exit 1
fi

# 2. Check/Create Virtual Environment
if [ ! -d "$VENV_DIR" ]; then
    echo -e "${YELLOW}📦 Creating virtual environment...${NC}"
    python3 -m venv "$VENV_DIR"
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Failed to create virtual environment.${NC}"
        exit 1
    fi
fi

# 3. Activate Virtual Environment
source "$VENV_DIR/bin/activate"

# 4. Dependency Check
echo -e "🔍 Checking dependencies..."
if ! pip freeze | grep -q "customtkinter"; then
    echo -e "${YELLOW}⬇️  Installing dependencies...${NC}"
    pip install -r "$REQUIREMENTS"
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Failed to install dependencies.${NC}"
        echo "Check your internet connection or 'requirements.txt'."
        exit 1
    fi
    echo -e "${GREEN}✅ Dependencies installed.${NC}"
else
    echo -e "${GREEN}✅ System ready.${NC}"
fi

# 5. Main Loop
while true; do
    echo ""
    echo -e "${CYAN}--- Menu ---${NC}"
    echo "1) GUI Mode (Graphical Interface)"
    echo "2) CLI Mode (Command Line)"
    echo "3) Update System (Git Pull + Deps)"
    echo "4) Exit"
    echo ""
    read -p "Select option [1-4]: " choice

    case $choice in
        1)
            echo -e "${GREEN}🚀 Starting GUI...${NC}"
            python main.py gui
            ;;
        2)
            echo -e "${GREEN}🚀 Starting CLI...${NC}"
            read -p "Enter Vault Password (optional, press Enter for none): " pass
            if [ -z "$pass" ]; then
                python main.py cli
            else
                python main.py cli --password "$pass"
            fi
            ;;
        3)
            echo -e "${YELLOW}🔄 Updating System...${NC}"
            if command -v git &> /dev/null; then
                git pull
                if [ $? -eq 0 ]; then
                     echo -e "${GREEN}✅ Code updated.${NC}"
                     echo -e "⬇️  Updating dependencies..."
                     pip install -r "$REQUIREMENTS" --upgrade
                else
                     echo -e "${RED}❌ Git pull failed. (Are you in a git repo?)${NC}"
                fi
            else
                echo -e "${RED}❌ Git is not installed.${NC}"
            fi
            ;;
        4)
            echo "Exiting..."
            deactivate
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Invalid option. Please try again.${NC}"
            ;;
    esac
done
