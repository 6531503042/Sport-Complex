#!/bin/bash

# Enhanced colors and styles
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
GRAY='\033[0;90m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# Unicode symbols for better visualization
CHECK_MARK="✓"
CROSS_MARK="✗"
ARROW="➜"
GEAR="⚙"
ROCKET="🚀"
WARNING="⚠"
STOP_SIGN="🛑"

# Service names and their corresponding env files
SERVICES=("user" "auth" "facility" "booking" "payment")
ENV_FILES=("env/dev/.env.user" "env/dev/.env.auth" "env/dev/.env.facility" "env/dev/.env.booking" "env/dev/.env.payment")

# Detect OS
OS="$(uname)"

# Add these new spinner styles at the top with other constants
SPINNER_STYLE=("⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏")
LOADING_BAR="▓"
EMPTY_BAR="░"

# Create logs directory if it doesn't exist
mkdir -p ./logs

# Function to get env file for service
get_env_file() {
    local service=$1
    local index=0
    for s in "${SERVICES[@]}"
    do
        if [ "$s" = "$service" ]; then
            echo "${ENV_FILES[$index]}"
            return
        fi
        index=$((index + 1))
    done
}

# Function to check if service exists
service_exists() {
    local service=$1
    for s in "${SERVICES[@]}"
    do
        if [ "$s" = "$service" ]; then
            return 0
        fi
    done
    return 1
}

# Function to open a new terminal window based on OS
open_terminal() {
    local service=$1
    local env_file=$2
    local current_dir=$(pwd)
    
    if [ "$OS" = "Darwin" ]; then
        # macOS - Run in background without new terminal window
        nohup go run main.go "$env_file" > "./logs/${service}.log" 2>&1 &
    elif [ "$OS" = "Linux" ]; then
        # Linux - Run in background
        nohup go run main.go "$env_file" > "./logs/${service}.log" 2>&1 &
    else
        # Fallback for other systems
        go run main.go "$env_file" > "./logs/${service}.log" 2>&1 &
    fi
}

# Function to print fancy headers
print_header() {
    local text=$1
    local length=${#text}
    local padding=$((50 - length))
    local half_padding=$((padding / 2))
    
    echo
    echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
    printf "${BLUE}║${NC}%*s${BOLD}%s${NC}%*s${BLUE}║${NC}\n" $half_padding "" "$text" $half_padding ""
    echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
    echo
}

# Function to show spinner
show_spinner() {
    local pid=$1
    local message=$2
    local spin='⣾⣽⣻⢿⡿⣟⣯⣷'
    local i=0
    
    while kill -0 $pid 2>/dev/null; do
        i=$(( (i+1) % ${#spin} ))
        printf "\r${CYAN}${spin:$i:1}${NC} ${message}"
        sleep .1
    done
    printf "\r"
}

# Add this new function for fancy loading bars
show_loading_bar() {
    local text=$2
    local bar_size=20
    
    echo -ne "\n"
    for ((i=0; i<=$bar_size; i++)); do
        local percentage=$((i * 100 / bar_size))
        local filled_bars=$((i))
        local empty_bars=$((bar_size - i))
        
        echo -ne "\r${CYAN}${text}${NC} ["
        for ((j=0; j<filled_bars; j++)); do
            echo -ne "${GREEN}${LOADING_BAR}${NC}"
        done
        for ((j=0; j<empty_bars; j++)); do
            echo -ne "${GRAY}${EMPTY_BAR}${NC}"
        done
        echo -ne "] ${percentage}%"
        sleep 0.02
    done
    echo -ne "\n"
}

# Enhanced spinner function
show_fancy_spinner() {
    local pid=$1
    local message=$2
    local i=0
    
    while kill -0 $pid 2>/dev/null; do
        echo -ne "\r${CYAN}${SPINNER_STYLE[i]} ${message}${NC}"
        i=$(( (i+1) % ${#SPINNER_STYLE[@]} ))
        sleep 0.1
    done
    echo -ne "\r"
}

# Add this function at the top with other functions
get_service_port() {
    local service=$1
    local env_file=$(get_env_file "$service")
    
    if [ -f "$env_file" ]; then
        # Extract port from env file
        local port=$(grep "PORT=" "$env_file" | cut -d'=' -f2)
        echo "$port"
    fi
}

# Function to kill process using a specific port
kill_port_process() {
    local port=$1
    
    # Kill process using port with sudo if needed
    if [ "$OS" = "Darwin" ]; then
        # For macOS - try without sudo first
        local pids=$(lsof -i -Pn | grep ":${port}" | awk '{print $2}')
        if [ -z "$pids" ]; then
            # If no PIDs found, try with sudo
            pids=$(sudo lsof -i -Pn | grep ":${port}" | awk '{print $2}')
        fi
        if [ ! -z "$pids" ]; then
            echo -e "${YELLOW}${WARNING} Port ${port} is in use. Stopping existing process...${NC}"
            echo "$pids" | xargs kill -9 2>/dev/null || echo "$pids" | xargs sudo kill -9 2>/dev/null
        fi
    else
        # For Linux
        local pids=$(lsof -i -Pn | grep ":${port}" | awk '{print $2}' || sudo lsof -i -Pn | grep ":${port}" | awk '{print $2}')
        if [ ! -z "$pids" ]; then
            echo -e "${YELLOW}${WARNING} Port ${port} is in use. Stopping existing process...${NC}"
            echo "$pids" | xargs kill -9 2>/dev/null || echo "$pids" | xargs sudo kill -9 2>/dev/null
        fi
    fi
    
    sleep 1
}

# Modified start_service function
start_service() {
    local service=$1
    local env_file=$(get_env_file "$service")
    local port=$(get_service_port "$service")
    
    if [ ! -f "$env_file" ]; then
        echo -e "\n${RED}${CROSS_MARK} Environment file not found: $env_file${NC}"
        return 1
    fi
    
    echo -e "\n${CYAN}${ROCKET} Starting ${BOLD}$service${NC}"
    show_loading_bar 0.5 "Initializing"
    
    # Kill any existing instances first
    stop_service "$service" > /dev/null 2>&1
    
    # Start the service
    open_terminal "$service" "$env_file"
    
    # Wait for service to start
    for i in {1..10}; do
        sleep 1
        if pgrep -f "go run main.go.*$service" > /dev/null; then
            if [ ! -z "$port" ]; then
                if [ "$OS" = "Darwin" ]; then
                    if lsof -i :${port} | grep "main" > /dev/null; then
                        echo -e "${GREEN}${CHECK_MARK} ${service}${NC} ${DIM}ready on port ${port}${NC}"
                        return 0
                    fi
                else
                    if netstat -tlpn 2>/dev/null | grep ":${port}" | grep "main" > /dev/null; then
                        echo -e "${GREEN}${CHECK_MARK} ${service}${NC} ${DIM}ready on port ${port}${NC}"
                        return 0
                    fi
                fi
            else
                echo -e "${GREEN}${CHECK_MARK} ${service}${NC} ${DIM}ready${NC}"
                return 0
            fi
        fi
    done
    
    echo -e "${RED}${CROSS_MARK} ${service}${NC} ${DIM}failed to start${NC}"
    return 1
}

# Add this helper function to check if a port is in use
is_port_in_use() {
    local port=$1
    if [ "$OS" = "Darwin" ]; then
        lsof -i :${port} >/dev/null 2>&1
        return $?
    else
        netstat -tuln | grep ":${port} " >/dev/null 2>&1
        return $?
    fi
}

# Modified status function to check ports
status() {
    echo -e "\n${BLUE}╭───────────────────────────────────╮${NC}"
    echo -e "${BLUE}│${NC} ${BOLD}Services Status${NC}                  ${BLUE}│${NC}"
    echo -e "${BLUE}╰───────────────────────────────────╯${NC}\n"
    
    local running=0
    local total=${#SERVICES[@]}
    
    # Define service ports
    declare -A SERVICE_PORTS=(
        ["user"]="1323"
        ["auth"]="1326"
        ["facility"]="1335"
        ["booking"]="1327"
        ["payment"]="1325"
    )
    
    for service in "${SERVICES[@]}"; do
        local port="${SERVICE_PORTS[$service]}"
        if [ ! -z "$port" ] && is_port_in_use "$port"; then
            echo -e "${GREEN}${CHECK_MARK}${NC} ${BOLD}${service}${NC} ${DIM}(running on port ${port})${NC}"
            ((running++))
        else
            echo -e "${RED}${CROSS_MARK}${NC} ${BOLD}${service}${NC} ${DIM}(stopped)${NC}"
        fi
    done
    
    echo -e "\n${DIM}Running: ${GREEN}$running${NC}${DIM} of ${total} services${NC}"
    
    # Show detailed port status
    echo -e "\n${BLUE}Port Status:${NC}"
    for service in "${SERVICES[@]}"; do
        local port="${SERVICE_PORTS[$service]}"
        if [ ! -z "$port" ]; then
            if is_port_in_use "$port"; then
                echo -e "${DIM}  ├─ ${NC}${GREEN}:${port}${NC} ${DIM}(${service})${NC}"
            else
                echo -e "${DIM}  ├─ ${NC}${RED}:${port}${NC} ${DIM}(${service} - not active)${NC}"
            fi
        fi
    done
    echo -e "${DIM}  └─${NC}"
}

# Modified stop_service function
stop_service() {
    local service=$1
    local port="${SERVICE_PORTS[$service]}"
    echo -e "\n${RED}${STOP_SIGN} Stopping ${BOLD}$service${NC}"

    # First: Kill IDE terminal processes (more aggressive)
    if [ ! -z "$port" ]; then
        # Kill all processes on the port (IDE terminal specific)
        sudo kill -9 $(sudo lsof -t -i:${port}) 2>/dev/null
        
        # Force kill any remaining port processes
        sudo lsof -i :${port} | awk 'NR>1 {print $2}' | sudo xargs -r kill -9 2>/dev/null
        
        # Additional port killing methods
        sudo fuser -k -n tcp ${port} 2>/dev/null
        sudo netstat -tlpn 2>/dev/null | grep ":${port}" | awk '{print $7}' | cut -d'/' -f1 | sudo xargs -r kill -9 2>/dev/null
    fi

    # Second: Kill all possible process variations
    local patterns=(
        "go run main.go.*${service}"
        "go build.*${service}"
        "__debug_bin.*${service}"
        "dlv.*${service}"          # Debug processes
        "gopls.*${service}"        # Go language server
        ".*${service}.*"           # Any process containing service name
    )

    for pattern in "${patterns[@]}"; do
        # Kill with different signals
        sudo pkill -SIGTERM -f "${pattern}" 2>/dev/null
        sleep 1
        sudo pkill -SIGKILL -f "${pattern}" 2>/dev/null
    done

    # Third: Kill parent processes that might be keeping the service alive
    ps -ef | grep "${service}" | grep -v grep | awk '{print $3}' | sudo xargs -r kill -9 2>/dev/null

    # Fourth: Kill terminal sessions running the service
    ps aux | grep "[t]erminal.*${service}" | awk '{print $2}' | sudo xargs -r kill -9 2>/dev/null
    
    # Wait for processes to die
    sleep 2

    # Final verification
    local is_running=0
    
    # Check port
    if [ ! -z "$port" ] && (sudo lsof -i :${port} >/dev/null 2>&1); then
        is_running=1
        # One last attempt to kill port
        sudo kill -9 $(sudo lsof -t -i:${port}) 2>/dev/null
    fi

    # Check processes
    if pgrep -f ".*${service}.*" >/dev/null 2>&1; then
        is_running=1
        # One last attempt to kill processes
        sudo pkill -9 -f ".*${service}.*" 2>/dev/null
    fi

    # Reset terminal if needed
    if [ $is_running -eq 1 ]; then
        # Reset terminal
        reset >/dev/null 2>&1
        clear
        echo -e "${RED}${WARNING} ${service}${NC} ${DIM}could not be stopped completely${NC}"
        echo -e "${DIM}You may need to restart your IDE terminal${NC}"
        return 1
    else
        echo -e "${GREEN}${CHECK_MARK} ${service}${NC} ${DIM}stopped successfully${NC}"
        return 0
    fi
}

# Function to restart a service
restart_service() {
    local service=$1
    stop_service "$service"
    sleep 2
    start_service "$service"
}

# Enhanced main command handler
case "$1" in
    start)
        print_header "Starting Services"
        if [ "$2" ]; then
            if service_exists "$2"; then
                start_service "$2"
            else
                echo -e "${RED}${CROSS_MARK} Invalid service: $2${NC}"
                echo -e "${GRAY}${ARROW} Available services: ${SERVICES[*]}${NC}"
                exit 1
            fi
        else
            echo -e "${CYAN}${ROCKET} Launching all services...${NC}"
            for service in "${SERVICES[@]}"
            do
                start_service "$service"
            done
        fi
        ;;
    stop)
        if [ "$2" ]; then
            if service_exists "$2"; then
                stop_service "$2"
            else
                echo -e "${RED}${CROSS_MARK} Invalid service: $2${NC}"
                echo -e "${GRAY}${ARROW} Available services: ${SERVICES[*]}${NC}"
                exit 1
            fi
        else
            for service in "${SERVICES[@]}"
            do
                stop_service "$service"
            done
        fi
        ;;
    restart)
        if [ "$2" ]; then
            if service_exists "$2"; then
                restart_service "$2"
            else
                echo -e "${RED}${CROSS_MARK} Invalid service: $2${NC}"
                echo -e "${GRAY}${ARROW} Available services: ${SERVICES[*]}${NC}"
                exit 1
            fi
        else
            for service in "${SERVICES[@]}"
            do
                restart_service "$service"
            done
        fi
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status} [service_name]"
        echo "Available services: ${SERVICES[*]}"
        exit 1
        ;;
esac

exit 0 