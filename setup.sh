#!/bin/bash
set -e
C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

info() {
    echo -e "${C_BLUE}[INFO]${C_RESET} $1"
}

success() {
    echo -e "${C_GREEN}[SUCCESS]${C_RESET} $1"
}

warn() {
    echo -e "${C_YELLOW}[WARNING]${C_RESET} $1"
}

error() {
    echo -e "${C_RED}[ERROR]${C_RESET} $1"
    exit 1
}

setup_pterodactyl() {
    info "Starting..."
    
    GO_VERSION="1.24.6"
    GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
    GO_INSTALL_DIR="/home/container/golang"
    GO_BIN_DIR="${GO_INSTALL_DIR}/bin"
    GO_TAR_FILE="go.tar.gz"
    
    if command -v go &> /dev/null; then
        success "Go already Installed : $(go version)"
    else
        warn "Go not found, Installing..."
        
		curl -s -L -o "${GO_TAR_FILE}" "${GO_URL}"
        
		mkdir -p "${GO_INSTALL_DIR}"
        info "Extracting..."
		tar -C "${GO_INSTALL_DIR}" -xzf "${GO_TAR_FILE}" --strip-components=1
        
		info "Cleaning..."
        rm "${GO_TAR_FILE}"
        
        info "Setup PATH"

        if ! grep -q "export PATH=${GO_BIN_DIR}" ~/.bashrc; then
            echo "" >> ~/.bashrc
            echo "export PATH=${GO_BIN_DIR}:\$PATH" >> ~/.bashrc
        fi

		export PATH="${GO_BIN_DIR}:$PATH"
        
        success "Go Installed!"
        info "Go version : $(go version)"
    fi

	# setup temporary dir
    info "Setup temporary directory..."
    TMP_DIR="/home/container/.tmp"
	mkdir -p "${TMP_DIR}"
    export GOTMPDIR="${TMP_DIR}"
    export TMPDIR="${TMP_DIR}"

	# add tmpdir to .bashrc
 	if ! grep -q "export GOTMPDIR=${TMP_DIR}" ~/.bashrc; then
        echo "export GOTMPDIR=${TMP_DIR}" >> ~/.bashrc
        echo "export TMPDIR=${TMP_DIR}" >> ~/.bashrc
    fi
	source ~/.bashrc
    
    # copy env example to .env
    cp .env.example .env
    
    success "Done!"
    echo "$(echo -e ${C_YELLOW}'Execute? [1] Run (development) / [2] Build (production): '${C_RESET})"
    read -p "" run_mode
    
    case "$run_mode" in
        1)
            info "Starting Sora with 'go run .'"
            info "This might take a while, please wait..."
            go run .
            ;;
        2)
            info "Building binary..."
            go build -ldflags="-s -w" .
            success "Done!"
            info "Executing './sora'..."
            ./sora
            ;;
        *)
            error "Invalid, exiting..."
            ;;
    esac
}

setup_termux() {
    warn "Coming soon."
    exit 0
}

setup_vps() {
    warn "Coming soon."
    exit 0
}


clear
echo -e "${C_BLUE}=====================================${C_RESET}"
echo -e "${C_BLUE}   Sora Bot Wangsaf Setup Script    ${C_RESET}"
echo -e "${C_BLUE}=====================================${C_RESET}"
echo -e "${C_YELLOW}Choose one where you will run Sora?\n[1] Panel Pterodactyl \n[2] Termux \n[3] VPS\n1 / 2 /3 ? ${C_RESET}"

read -p "$(echo -e ${C_YELLOW}'1 / 2 / 3 : '${C_RESET})" choice

case "$choice" in
    1)
        setup_pterodactyl
        ;;
    2)
        setup_termux
        ;;
    3)
        setup_vps
        ;;
    *)
        error "Invalid, Exiting..."
        ;;
esac

success "Setup Completed!"
