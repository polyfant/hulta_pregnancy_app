#!/bin/bash
# Hetzner Server Preparation Script
# Version: 2.1
# Date: January 2025
# Maintainer: Hulta Pregnancy App Team

# Enhanced error handling and logging
set -euo pipefail

# Logging and error handling functions
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [PREPARE] $*" | tee -a ~/hulta_server_prep.log
}

error_exit() {
    log "ERROR: $1"
    echo "ERROR: $1" >&2
    exit 1
}

# Trap any errors
trap 'error_exit "Command failed: $BASH_COMMAND"' ERR

# Validate script requirements
validate_environment() {
    log "Validating system environment"
    
    # Check Ubuntu version
    if ! grep -q 'Ubuntu 22.04' /etc/os-release; then
        error_exit "Unsupported OS. Requires Ubuntu 22.04 LTS"
    fi

    # Check current user
    if [[ "$USER" != "deploy" ]]; then
        error_exit "Script must be run as 'deploy' user"
    fi

    # Check sudo access
    if ! sudo -v; then
        error_exit "User lacks sudo privileges"
    fi

    # Check available disk space
    if [[ $(df -h / | awk '/\// {print $5}' | sed 's/%//') -gt 80 ]]; then
        error_exit "Insufficient disk space. Need at least 20% free space."
    fi
}

# System update and preparation
prepare_system() {
    log "Updating system packages"
    sudo apt-get update || error_exit "Failed to update package lists"
    sudo apt-get upgrade -y || error_exit "System upgrade failed"
    
    log "Installing essential system tools"
    sudo apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release \
        software-properties-common \
        git \
        wget \
        unzip \
        || error_exit "Failed to install essential tools"
}

# Secure SSH configuration
configure_ssh() {
    log "Configuring SSH for enhanced security"
    sudo cp /etc/ssh/sshd_config /etc/ssh/sshd_config.backup

    # Advanced SSH hardening
    sudo sed -i 's/^#*Port 22/Port 22022/' /etc/ssh/sshd_config
    sudo sed -i 's/^#*PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config
    sudo sed -i 's/^#*PubkeyAuthentication.*/PubkeyAuthentication yes/' /etc/ssh/sshd_config
    sudo sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config
    sudo sed -i 's/^#*MaxAuthTries.*/MaxAuthTries 3/' /etc/ssh/sshd_config
    sudo sed -i 's/^#*AllowUsers.*/AllowUsers deploy/' /etc/ssh/sshd_config

    # Add additional SSH security
    echo "
# Additional Security Settings
Protocol 2
ClientAliveInterval 300
ClientAliveCountMax 0
IgnoreRhosts yes
RhostsRSAAuthentication no
HostbasedAuthentication no
" | sudo tee -a /etc/ssh/sshd_config > /dev/null

    sudo systemctl restart ssh || error_exit "Failed to restart SSH service"
}

# Firewall configuration
configure_firewall() {
    log "Configuring UFW firewall"
    sudo ufw default deny incoming
    sudo ufw default allow outgoing
    sudo ufw allow 22022/tcp  # SSH on non-standard port
    sudo ufw allow 80/tcp     # HTTP
    sudo ufw allow 443/tcp    # HTTPS
    sudo ufw --force enable || error_exit "Failed to enable UFW"
}

# Docker installation
install_docker() {
    log "Installing Docker and Docker Compose"
    # Remove any existing Docker installations
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
        sudo apt-get remove -y $pkg
    done

    # Official Docker installation
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg

    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin \
        || error_exit "Docker installation failed"

    # Ensure deploy user is in docker group
    sudo usermod -aG docker $USER
}

# Node.js and npm installation via NVM
install_nodejs() {
    log "Installing Node Version Manager (NVM) and Node.js"
    
    # Install NVM
    if [ ! -d "$HOME/.nvm" ]; then
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
        
        # Source NVM for the current script
        export NVM_DIR="$HOME/.nvm"
        [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
        [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
    fi

    # Install latest LTS version of Node.js
    nvm install --lts
    nvm use --lts
    
    # Install global npm packages
    npm install -g npm@latest
    npm install -g yarn pnpm

    # Verify installations
    node --version
    npm --version
    
    log "Node.js and npm installed successfully"
}

# Configure deployment environment
configure_deployment() {
    log "Configuring deployment environment"
    
    # Create project directory
    sudo mkdir -p /opt/hulta-pregnancy-app
    sudo chown $USER:$USER /opt/hulta-pregnancy-app

    # Secure SSH directory
    mkdir -p ~/.ssh
    chmod 700 ~/.ssh
    
    # Ensure authorized_keys exists with correct permissions
    touch ~/.ssh/authorized_keys
    chmod 600 ~/.ssh/authorized_keys
}

# Additional security tools
install_security_tools() {
    log "Installing security tools"
    sudo apt-get install -y \
        fail2ban \
        auditd \
        || error_exit "Failed to install security tools"

    # Fail2ban configuration
    sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
    sudo sed -i 's/bantime  = 10m/bantime  = 1h/' /etc/fail2ban/jail.local
    sudo sed -i 's/maxretry = 5/maxretry = 3/' /etc/fail2ban/jail.local

    sudo systemctl enable fail2ban
    sudo systemctl restart fail2ban
}

# Main execution
main() {
    log "Starting Hulta Pregnancy App Server Preparation"
    
    validate_environment
    prepare_system
    configure_ssh
    configure_firewall
    install_docker
    install_nodejs
    configure_deployment
    install_security_tools

    log "Server preparation completed successfully!"
    log "Next steps: Clone project, set up environment, deploy application"
}

# Run main function
main

exit 0
