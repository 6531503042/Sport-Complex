# Colors and formatting
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
CYAN := \033[0;36m
GRAY := \033[0;90m
NC := \033[0m
BOLD := \033[1m
DIM := \033[2m

# Symbols
ROCKET := 🚀
CHECK := ✓
CROSS := ✗
GEAR := ⚙️
STOP := 🛑

# Progress bar characters
LOADING := ▓
EMPTY := ░

.PHONY: start stop restart status build

# Default target
all: help

# Function to show loading animation
define show_loading
	@for i in 1 2 3; do \
		printf "\r${CYAN}$(1) [${LOADING}${LOADING}   ]${NC}"; sleep 0.05; \
		printf "\r${CYAN}$(1) [${LOADING}${LOADING}${LOADING}  ]${NC}"; sleep 0.05; \
		printf "\r${CYAN}$(1) [${LOADING}${LOADING}${LOADING}${LOADING}]${NC}"; sleep 0.05; \
	done
	@echo "\n"
endef

# Make script executable
init:
	@echo "${BLUE}${GEAR} Initializing ecosystem...${NC}"
	@chmod +x ./scripts/ecosystem.sh
	@echo "${GREEN}${CHECK} Ecosystem ready!${NC}"

# Start all services or a specific service
start:
ifdef service
	@echo -e "\n${CYAN}${ROCKET} Starting ${BOLD}${service}${NC}"
	$(call show_loading,"Initializing")
	@./scripts/ecosystem.sh start $(service)
else
	@echo -e "\n${CYAN}${ROCKET} ${BOLD}Launching Services${NC}"
	@echo -e "${DIM}Order: user → auth → facility → booking → payment${NC}\n"
	$(call show_loading,"Preparing cluster")
	@./scripts/ecosystem.sh start
endif

# Stop all services or a specific service
stop:
ifdef service
	@echo -e "\n${RED}${STOP} Stopping ${BOLD}${service}${NC}"
	$(call show_loading,"Shutting down")
	@./scripts/ecosystem.sh stop $(service)
else
	@echo -e "\n${RED}${STOP} ${BOLD}Stopping All Services${NC}"
	$(call show_loading,"Initiating shutdown")
	@./scripts/ecosystem.sh stop
endif

# Restart all services or a specific service
restart:
ifdef service
	@echo -e "\n${YELLOW}${GEAR} Restarting ${BOLD}${service}${NC}"
	$(call show_loading,"Restarting")
	@./scripts/ecosystem.sh restart $(service)
else
	@echo -e "\n${YELLOW}${GEAR} ${BOLD}Restarting All Services${NC}"
	$(call show_loading,"Restarting cluster")
	@./scripts/ecosystem.sh restart
endif

# Show status of all services
status:
	@echo "\n${BLUE}╭─────────────────────────────────╮${NC}"
	@echo "${BLUE}│${NC}  ${BOLD}Services Status Overview${NC}         ${BLUE}│${NC}"
	@echo "${BLUE}╰─────────────────────────────────╯${NC}\n"
	@./scripts/ecosystem.sh status

# Help command
help:
	@echo "\n${YELLOW}╔════════════════════════════════════════╗${NC}"
	@echo "${YELLOW}║       ${BOLD}Service Management System${NC}        ${YELLOW}║${NC}"
	@echo "${YELLOW}╚════════════════════════════════════════╝${NC}\n"
	@echo "${BOLD}Available Commands:${NC}"
	@echo "  ${GREEN}make init${NC}              ${DIM}Initialize the ecosystem${NC}"
	@echo "  ${GREEN}make start${NC}             ${DIM}Start all services${NC}"
	@echo "  ${GREEN}make start service=X${NC}    ${DIM}Start specific service${NC}"
	@echo "  ${GREEN}make stop${NC}              ${DIM}Stop all services${NC}"
	@echo "  ${GREEN}make stop service=X${NC}     ${DIM}Stop specific service${NC}"
	@echo "  ${GREEN}make restart${NC}           ${DIM}Restart all services${NC}"
	@echo "  ${GREEN}make restart service=X${NC}  ${DIM}Restart specific service${NC}"
	@echo "  ${GREEN}make status${NC}            ${DIM}Show services status${NC}\n"
	@echo "${BOLD}Examples:${NC}"
	@echo "  ${CYAN}$$ make start service=auth${NC}    ${DIM}Start auth service${NC}"
	@echo "  ${CYAN}$$ make stop service=user${NC}     ${DIM}Stop user service${NC}"
	@echo "  ${CYAN}$$ make status${NC}               ${DIM}Show all services status${NC}\n" 