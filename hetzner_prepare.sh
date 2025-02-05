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
GIT_REPO="git@github.com:polyfant/hulta_pregnancy_app.git"
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

# Add this function after prepare_system()
cleanup_system() {
    log "🧹 Cleaning up unnecessary packages"
    sudo apt-get autoremove -y
    sudo apt-get clean
    log "✅ System cleanup complete"
}

# Secure SSH configuration
configure_ssh() {
    log "🔒 Configuring SSH for enhanced security"
    
    # Backup original config if backup doesn't exist
    if [ ! -f /etc/ssh/sshd_config.backup ]; then
        sudo cp /etc/ssh/sshd_config /etc/ssh/sshd_config.backup
    fi

    # Create a new minimal SSH config
    cat <<EOF | sudo tee /etc/ssh/sshd_config > /dev/null
# Basic SSH Configuration
Port 22
PermitRootLogin no
PasswordAuthentication yes
PubkeyAuthentication yes
ChallengeResponseAuthentication no
UsePAM yes
X11Forwarding no
PrintMotd no
AcceptEnv LANG LC_*
Subsystem sftp /usr/lib/openssh/sftp-server
EOF

    # Test configuration before applying
    sudo sshd -t || {
        log "SSH configuration test failed. Restoring backup..."
        sudo cp /etc/ssh/sshd_config.backup /etc/ssh/sshd_config
        error_exit "SSH configuration test failed"
    }
    
    sudo systemctl restart ssh || error_exit "Failed to restart SSH service"
    log "✅ SSH configuration updated successfully"
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

# Docker installation with better error handling
install_docker() {
    log "🐳 Installing Docker and Docker Compose"
    
    # Check if Docker is already installed
    if command -v docker &> /dev/null; then
        log "Docker is already installed: $(docker --version)"
        return 0
    fi

    # Remove any conflicting packages
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
        sudo apt-get remove -y $pkg &> /dev/null || true
    done

    # Add Docker's official GPG key
    if [ ! -f /etc/apt/keyrings/docker.gpg ]; then
        sudo install -m 0755 -d /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        sudo chmod a+r /etc/apt/keyrings/docker.gpg
    fi

    # Add the repository to Apt sources
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
        $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
        sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin \
        || error_exit "Docker installation failed"

    # Ensure deploy user is in docker group
    sudo usermod -aG docker $USER
    
    # Verify Docker installation
    docker --version || error_exit "Docker installation verification failed"
    
    log "✅ Docker installation complete"
}

# Node.js and npm installation via NVM
install_nodejs() {
    log "Installing Node Version Manager (NVM) and Node.js"
    
    # Install NVM if not already installed
    if [ ! -d "$HOME/.nvm" ]; then
        # Download and run the NVM installation script
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash || error_exit "Failed to install NVM"
        
        # Ensure NVM is loaded in current shell
        export NVM_DIR="$HOME/.nvm"
        [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    fi
    
    # Reload shell environment to ensure NVM is available
    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    
    # Check if Node.js is already installed
    if command -v node &> /dev/null; then
        log "Node.js is already installed: $(node --version)"
    else
        # Install latest LTS version of Node.js
        set +u  # Temporarily disable undefined variable checking
        nvm install --lts || error_exit "Failed to install Node.js"
        nvm use --lts || error_exit "Failed to use Node.js"
        set -u  # Re-enable undefined variable checking
    fi
    
    # Verify installations
    if ! command -v node &> /dev/null; then
        error_exit "Node.js installation verification failed"
    fi
    
    if ! command -v npm &> /dev/null; then
        error_exit "npm installation verification failed"
    fi
    
    # Show versions
    log "Node.js version: $(node --version)"
    log "npm version: $(npm --version)"
    
    log "✅ Node.js and npm installation verified"
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
    log "Setting up environment variables"
    
    # Create .env file if it doesn't exist
    if [ ! -f "$ENV_FILE" ]; then
        cat <<EOF > "$ENV_FILE"
# Auth0 Configuration
VITE_AUTH0_DOMAIN=dev-r083cwkcv0pgz20x.eu.auth0.com
VITE_AUTH0_CLIENT_ID=OBmEol1z4U49r3YI3priDdGbvF5i4O7d
VITE_AUTH0_AUDIENCE=https://api.hulta-foaltracker.app

# API Configuration
VITE_API_URL=https://api.hulta-foaltracker.app
VITE_WEBSOCKET_URL=wss://api.hulta-foaltracker.app/ws

# Database Configuration
POSTGRES_DB=hulta_db
POSTGRES_USER=hulta_user
POSTGRES_PASSWORD=$(openssl rand -base64 32)
DATABASE_URL=postgresql://hulta_user:${POSTGRES_PASSWORD}@postgres:5432/hulta_db

# JWT Configuration
JWT_SECRET=$(openssl rand -base64 32)
EOF
    fi
    
    log "✅ Environment setup complete"
}

# Add after configure_deployment()
setup_ssh_keys() {
    log "🔑 Checking SSH configuration"
    
    # Check for any existing SSH keys
    EXISTING_KEYS=( "$HOME/.ssh/id_rsa" "$HOME/.ssh/id_ed25519" )
    SSH_KEY_FILE=""
    
    for key in "${EXISTING_KEYS[@]}"; do
        if [ -f "$key" ]; then
            SSH_KEY_FILE="$key"
            log "✅ Found existing SSH key: $SSH_KEY_FILE"
            break
        fi
    done
    
    # If no key was found, error out since we want to use existing keys
    if [ -z "$SSH_KEY_FILE" ]; then
        error_exit "No existing SSH keys found. Please add your SSH key to the server first."
    fi

    # Add GitHub to known_hosts if not already there
    if ! grep -q "github.com" ~/.ssh/known_hosts 2>/dev/null; then
        ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null
    fi
    
    # Since we know the key is already in GitHub, just verify it exists
    log "✅ Using existing SSH key: ${SSH_KEY_FILE}"
    log "🔑 Key fingerprint: $(ssh-keygen -lf $SSH_KEY_FILE)"
    
    # Set proper permissions just in case
    chmod 600 $SSH_KEY_FILE
    
    log "✅ SSH configuration verified"
}

# Update setup_repository() with better error handling
setup_repository() {
    log "📂 Setting up application repository"
    
    # Ensure git is configured
    git config --global user.email "deploy@hulta-foaltracker.app"
    git config --global user.name "Deploy Bot"
    
    if [ -d "$APP_DIR/.git" ]; then
        log "⚠️ Application directory exists, updating repository"
        cd $APP_DIR || error_exit "Failed to change to app directory"
        
        # Backup any local changes
        if [ -n "$(git status --porcelain)" ]; then
            BACKUP_DIR="$APP_DIR/backup_$(date +%Y%m%d_%H%M%S)"
            log "📦 Creating backup of local changes to $BACKUP_DIR"
            mkdir -p "$BACKUP_DIR"
            git diff > "$BACKUP_DIR/local_changes.patch"
        fi
        
        # Update repository
        git fetch --all || error_exit "Failed to fetch repository updates"
        git reset --hard origin/$BRANCH || error_exit "Failed to reset to $BRANCH"
    else
        log "📥 Cloning repository"
        git clone -b $BRANCH $GIT_REPO $APP_DIR || error_exit "Failed to clone repository"
        cd $APP_DIR || error_exit "Failed to change to app directory"
    fi
    
    log "✅ Repository setup complete"
}

# Add Nginx configuration
setup_nginx() {
    log "🌐 Setting up Nginx"
    
    sudo apt-get install -y nginx || error_exit "Failed to install Nginx"
    
    # Create Nginx configuration
    cat <<EOF | sudo tee /etc/nginx/sites-available/hulta-foaltracker.conf > /dev/null
server {
    listen 80;
    server_name hulta-foaltracker.app www.hulta-foaltracker.app;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}

server {
    listen 80;
    server_name api.hulta-foaltracker.app;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF

    # Enable the site
    sudo ln -sf /etc/nginx/sites-available/hulta-foaltracker.conf /etc/nginx/sites-enabled/
    sudo rm -f /etc/nginx/sites-enabled/default
    
    # Test and reload Nginx
    sudo nginx -t || error_exit "Nginx configuration test failed"
    sudo systemctl reload nginx
    
    log "✅ Nginx configured successfully"
}

# Update SSL setup to run after Nginx
setup_ssl() {
    log "🔒 Setting up SSL with Let's Encrypt"
    
    # Install Certbot and Nginx plugin
    sudo apt-get update
    sudo apt-get install -y certbot python3-certbot-nginx || error_exit "Failed to install Certbot"
    
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

# Add this function to verify project structure
verify_project_structure() {
    log "🔍 Verifying project structure"
    
    # Ensure docker-compose.yml exists
    if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
        error_exit "docker-compose.yml not found in $APP_DIR"
    fi
    
    log "✅ Project structure verified"
}

# Update start_services() to handle environment files better
start_services() {
    log "🐳 Starting Docker services"
    
    # Verify project structure first
    verify_project_structure
    
    # Ensure Docker is running
    sudo systemctl start docker
    
    # Create or update environment files
    if [ ! -f "$APP_DIR/.env" ]; then
        # If .env doesn't exist in APP_DIR, create it from ENV_FILE
        cp "$ENV_FILE" "$APP_DIR/.env"
    else
        # If both exist and are different, update APP_DIR/.env
        if ! cmp -s "$ENV_FILE" "$APP_DIR/.env"; then
            cp "$ENV_FILE" "$APP_DIR/.env"
        fi
    fi
    
    # Handle frontend environment file
    if [ -d "$APP_DIR/frontend-react" ]; then
        if [ ! -f "$APP_DIR/frontend-react/.env" ] || ! cmp -s "$ENV_FILE" "$APP_DIR/frontend-react/.env"; then
            cp "$ENV_FILE" "$APP_DIR/frontend-react/.env"
        fi
    fi
    
    # Export required environment variables
    if [ -f "$APP_DIR/.env" ]; then
        export DATABASE_URL=$(grep '^DATABASE_URL=' "$APP_DIR/.env" | cut -d '=' -f2 || echo "")
        export JWT_SECRET=$(grep '^JWT_SECRET=' "$APP_DIR/.env" | cut -d '=' -f2 || echo "")
    fi
    
    # Build and start containers with environment variables
    cd $APP_DIR
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

# Main execution
main() {
    log "🚀 Starting Hulta Pregnancy App Full Deployment"
    
    validate_environment
    prepare_system
    cleanup_system
    configure_ssh
    setup_ssh_keys
    configure_firewall
    install_docker
    install_nodejs
    configure_deployment
    install_security_tools
    setup_environment
    setup_repository
    verify_project_structure
    setup_nginx
    start_services
    setup_ssl
    verify_deployment

    log "🎉🎉🎉 Deployment completed successfully! 🎉🎉🎉"
    log "🌐 Application is now live at https://hulta-foaltracker.app"
    log "⚠️ Important: Make sure to add the displayed SSH public key to GitHub!"
}

# Run main function
main

exit 0