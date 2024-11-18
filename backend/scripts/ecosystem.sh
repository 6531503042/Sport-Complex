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
        # macOS - fixed AppleScript syntax
        osascript <<EOF
            tell application "Terminal"
                do script "cd \"${current_dir}\" && echo \"Starting ${service} service...\" && go run main.go \"${env_file}\""
                tell window 1
                    set custom title to "${service}"
                end tell
            end tell
EOF
    elif [ "$OS" = "Linux" ]; then
        # Linux with gnome-terminal
        gnome-terminal --title="$service" -- bash -c "cd \"$current_dir\" && echo -e '${GREEN}Starting $service service...${NC}' && go run main.go \"$env_file\"; exec bash"
    else
        # Fallback for other systems
        echo -e "${RED}Unsupported operating system for terminal windows${NC}"
        echo -e "${GREEN}Starting $service service in background...${NC}"
        cd "$current_dir" && go run main.go "$env_file" &
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

# Optimize the start_service function
start_service() {
    local service=$1
    local env_file=$(get_env_file "$service")
    
    if [ ! -f "$env_file" ]; then
        echo -e "\n${RED}${CROSS_MARK} Environment file not found: $env_file${NC}"
        return 1
    fi
    
    echo -e "\n${CYAN}${ROCKET} Starting ${BOLD}$service${NC}"
    show_loading_bar 0.5 "Initializing"
    
    open_terminal "$service" "$env_file" &
    sleep 0.5
    
    if pgrep -f "go run main.go.*$service" > /dev/null; then
        echo -e "${GREEN}${CHECK_MARK} ${service}${NC} ${DIM}ready${NC}"
    else
        echo -e "${RED}${CROSS_MARK} ${service}${NC} ${DIM}failed${NC}"
    fi
}

# Optimize the status function to be simpler and faster
status() {
    echo -e "\n${BLUE}╭───────────────────────╮${NC}"
    echo -e "${BLUE}│${NC} ${BOLD}Services Status${NC}      ${BLUE}│${NC}"
    echo -e "${BLUE}╰───────────────────────╯${NC}\n"
    
    local running=0
    local total=${#SERVICES[@]}
    
    for service in "${SERVICES[@]}"
    do
        if pgrep -f "go run main.go.*$service" > /dev/null; then
            echo -e "${GREEN}${CHECK_MARK}${NC} ${BOLD}${service}${NC} ${DIM}(running)${NC}"
            ((running++))
        else
            echo -e "${RED}${CROSS_MARK}${NC} ${BOLD}${service}${NC} ${DIM}(stopped)${NC}"
        fi
    done
    
    echo -e "\n${DIM}Running: ${GREEN}$running${NC}${DIM} of ${total}${NC}"
}

# Optimize the stop_service function
stop_service() {
    local service=$1
    echo -e "\n${RED}${STOP_SIGN} Stopping ${BOLD}$service${NC}"
    
    if [ "$OS" = "Darwin" ]; then
        osascript <<EOF
            tell application "Terminal"
                set windowList to every window whose name contains "${service}"
                repeat with windowItem in windowList
                    close windowItem
                end repeat
            end tell
EOF
    else
        pkill -f "go run main.go.*$service"
    fi
    
    sleep 0.5
    
    if ! pgrep -f "go run main.go.*$service" > /dev/null; then
        echo -e "${GREEN}${CHECK_MARK} ${service}${NC} ${DIM}stopped${NC}"
    else
        echo -e "${RED}${WARNING} ${service}${NC} ${DIM}failed to stop${NC}"
    fi
    
    # Close terminal window after stopping the service (macOS only)
    if [ "$OS" = "Darwin" ]; then
        osascript -e 'tell application "Terminal" to quit'
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