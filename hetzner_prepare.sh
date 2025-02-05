#!/bin/bash
# Hulta Pregnancy App Full Deployment Script
# Version: 3.0
# Date: January 2025
# Maintainer: Hulta Pregnancy App Team 🚀

# Enhanced error handling and logging
set -euo pipefail

# Configuration Variables
APP_DIR="/opt/hulta-pregnancy-app"
ENV_FILE="$APP_DIR/.env"
DOCKER_COMPOSE_FILE="$APP_DIR/docker-compose.yml"
GIT_REPO="https://github.com/polyfant/hulta_pregnancy_app.git"
BRANCH="main"

# Logging and error handling functions
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [DEPLOY] $*" | tee -a ~/hulta_deployment.log
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
    log "🔍 Validating system environment"
    
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
    log "🛠️ Updating system packages"
    sudo apt-get update || error_exit "Failed to update package lists"
    sudo apt-get upgrade -y || error_exit "System upgrade failed"
    
    log "📦 Installing essential system tools"
    sudo apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release \
        software-properties-common \
        git \
        wget \
        unzip \
        jq \
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
    sudo ufw allow 22022/tcp  # SSH
    sudo ufw allow 80/tcp     # HTTP
    sudo ufw allow 443/tcp    # HTTPS
    sudo ufw --force enable
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

# Configure environment
setup_environment() {
    log "🌍 Setting up environment variables"
    
    # Create the environment file if it doesn't exist
    touch $ENV_FILE
    chmod 600 $ENV_FILE
    
    # Function to get or generate a value
    get_or_generate_value() {
        local key=$1
        local default_value=$2
        local current_value=$(grep "^$key=" $ENV_FILE | cut -d '=' -f2)
        
        if [ -n "$current_value" ]; then
            echo "$current_value"
        else
            if [ -n "$default_value" ]; then
                echo "$default_value"
            else
                case $key in
                    POSTGRES_PASSWORD|REDIS_PASSWORD)
                        openssl rand -hex 16
                        ;;
                    *)
                        error_exit "Missing required environment variable: $key"
                        ;;
                esac
            fi
        fi
    }

    # Read or generate each value
    VITE_AUTH0_DOMAIN=$(get_or_generate_value "VITE_AUTH0_DOMAIN" "dev-r083cwkcv0pgz20x.eu.auth0.com")
    VITE_AUTH0_CLIENT_ID=$(get_or_generate_value "VITE_AUTH0_CLIENT_ID" "OBmEol1z4U49r3YI3priDdGbvF5i4O7d")
    VITE_AUTH0_AUDIENCE=$(get_or_generate_value "VITE_AUTH0_AUDIENCE" "https://api.hulta-foaltracker.app")
    VITE_API_URL=$(get_or_generate_value "VITE_API_URL" "https://api.hulta-foaltracker.app")
    VITE_WEBSOCKET_URL=$(get_or_generate_value "VITE_WEBSOCKET_URL" "wss://api.hulta-foaltracker.app/notifications")
    POSTGRES_DB=$(get_or_generate_value "POSTGRES_DB" "hulta_db")
    POSTGRES_USER=$(get_or_generate_value "POSTGRES_USER" "hulta_user")
    POSTGRES_PASSWORD=$(get_or_generate_value "POSTGRES_PASSWORD")
    REDIS_PASSWORD=$(get_or_generate_value "REDIS_PASSWORD")
    NODE_ENV=$(get_or_generate_value "NODE_ENV" "production")
    PORT=$(get_or_generate_value "PORT" "3000")

    # Write the environment file
    cat <<EOF > $ENV_FILE
# Auth0 Configuration
VITE_AUTH0_DOMAIN=${VITE_AUTH0_DOMAIN}
VITE_AUTH0_CLIENT_ID=${VITE_AUTH0_CLIENT_ID}
VITE_AUTH0_AUDIENCE=${VITE_AUTH0_AUDIENCE}

# API Configuration
VITE_API_URL=${VITE_API_URL}

# WebSocket Connection
VITE_WEBSOCKET_URL=${VITE_WEBSOCKET_URL}

# Database Configuration
POSTGRES_DB=${POSTGRES_DB}
POSTGRES_USER=${POSTGRES_USER}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# Redis Configuration
REDIS_PASSWORD=${REDIS_PASSWORD}

# App Configuration
NODE_ENV=${NODE_ENV}
PORT=${PORT}
EOF

    log "✅ Environment variables configured (existing values preserved)"
}

# Clone and setup repository
setup_repository() {
    log "📂 Setting up application repository"
    
    if [ -d "$APP_DIR" ]; then
        log "⚠️ Application directory exists, updating repository"
        cd $APP_DIR
        git fetch --all
        git reset --hard origin/$BRANCH
    else
        log "📥 Cloning repository"
        git clone -b $BRANCH $GIT_REPO $APP_DIR
        cd $APP_DIR
    fi
    
    log "✅ Repository setup complete"
}

# Start Docker services
start_services() {
    log "🐳 Starting Docker services"
    
    # Ensure Docker is running
    sudo systemctl start docker
    
    # Build and start containers
    docker compose -f $DOCKER_COMPOSE_FILE up -d --build || error_exit "Failed to start Docker services"
    
    log "✅ Services started successfully"
}

# Verify deployment
verify_deployment() {
    log "🔍 Verifying deployment"
    
    # Check if containers are running
    if ! docker compose -f $DOCKER_COMPOSE_FILE ps | grep -q "Up"; then
        error_exit "Some containers failed to start"
    fi
    
    # Check API health
    API_URL=$(grep VITE_API_URL $ENV_FILE | cut -d '=' -f2)
    if ! curl -s $API_URL/health | grep -q "OK"; then
        error_exit "API health check failed"
    fi
    
    log "🎉 Deployment verified successfully!"
}

# SSL setup
setup_ssl() {
    log "🔒 Setting up SSL with Let's Encrypt"
    
    # Install Certbot and Nginx plugin
    sudo apt-get update
    sudo apt-get install -y certbot python3-certbot-nginx
    
    # Get SSL certificate
    sudo certbot --nginx \
        -d hulta-foaltracker.app \
        -d www.hulta-foaltracker.app \
        -d api.hulta-foaltracker.app \
        --non-interactive \
        --agree-tos \
        --email hyvlaren@gmail.com \  
        --redirect \
        || error_exit "Failed to obtain SSL certificates"
        
    # Auto-renewal setup
    sudo systemctl enable certbot.timer
    sudo systemctl start certbot.timer
    
    log "✅ SSL certificates installed successfully"
}

# Main execution
main() {
    log "🚀 Starting Hulta Pregnancy App Full Deployment"
    
    validate_environment
    prepare_system
    configure_ssh
    configure_firewall
    install_docker
    install_nodejs
    configure_deployment
    install_security_tools
    setup_environment
    setup_repository
    setup_ssl
    start_services
    verify_deployment

    log "🎉🎉🎉 Deployment completed successfully! 🎉🎉🎉"
    log "🌐 Application is now live at https://api.hulta-foaltracker.app"
}

# Run main function
main

exit 0
