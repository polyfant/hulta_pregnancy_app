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
GIT_REPO="git@github.com:polyfant/hulta_pregnancy_app.git"  # Changed to SSH URL
BRANCH="main"

# Use the deploy user's SSH key
GIT_SSH_KEY="/home/deploy/.ssh/id_ed25519"

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
    if ! grep -q 'Ubuntu' /etc/os-release || ! grep -q '24.04\|22.04' /etc/os-release; then
        error_exit "Unsupported OS. Requires Ubuntu 24.04 or 22.04 LTS"
    fi

    # Check sudo access
    if ! sudo -v; then
        error_exit "User lacks sudo privileges"
    fi

    # Check available disk space (requiring at least 20GB free)
    FREE_SPACE=$(df -BG / | awk 'NR==2 {print $4}' | sed 's/G//')
    if [ "$FREE_SPACE" -lt 20 ]; then
        error_exit "Insufficient disk space. Need at least 20GB free space, found ${FREE_SPACE}GB"
    fi

    # Check memory requirements (minimum 2GB)
    TOTAL_MEM=$(grep MemTotal /proc/meminfo | awk '{print int($2/1024/1024)}')
    if [ "$TOTAL_MEM" -lt 2 ]; then
        error_exit "Insufficient memory. Need at least 2GB RAM, found ${TOTAL_MEM}GB"
    fi

    # Check for required ports availability (excluding port 22 which should be in use for SSH)
    for port in 80 443 5432 3000 8080; do
        if netstat -tuln | grep -q ":$port "; then
            error_exit "Port $port is already in use"
        fi
    done

    # Verify hostname resolution
    if ! host hulta-foaltracker.app >/dev/null 2>&1; then
        log "⚠️  Warning: Cannot resolve hulta-foaltracker.app - DNS may not be properly configured"
    fi

    log "✅ Environment validation passed"
}

validate_docker_compose() {
    log "🔍 Validating Docker Compose configuration"
    
    if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
        error_exit "Docker Compose file not found at $DOCKER_COMPOSE_FILE"
    fi

    # Test Docker Compose configuration
    cd "$APP_DIR" || error_exit "Failed to change to app directory"
    if ! docker compose -f "$DOCKER_COMPOSE_FILE" config > /dev/null; then
        error_exit "Docker Compose configuration is invalid"
    fi
    
    # Verify required services are defined
    if ! docker compose -f "$DOCKER_COMPOSE_FILE" config --services | grep -q "frontend\|backend\|postgres"; then
        error_exit "Required services (frontend, backend, postgres) not found in docker-compose.yml"
    fi

    log "✅ Docker Compose configuration is valid"
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
    sudo ufw allow 22/tcp  # SSH
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

# ...existing code...

# Configure environment
setup_environment() {
    log "Setting up environment variables"
    
    # Create backend .env file
    cat <<EOF > "$APP_DIR/backend/.env"
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=horsetracker
DB_PASSWORD=R2,Y@B&wO46.Ln}Q
DB_NAME=horse_tracking_db

# Auth0 Configuration
AUTH0_AUDIENCE=https://api.hulta-foaltracker.app  
AUTH0_ALGORITHM=RS256
AUTH0_DOMAIN=dev-r083cwkcv0pgz20x.eu.auth0.com
AUTH0_ISSUER=https://dev-r083cwkcv0pgz20x.eu.auth0.com/
AUTH0_CLIENT_ID=OBmEol1z4U49r3YI3priDdGbvF5i4O7d
AUTH0_CLIENT_SECRET=yTYWEdEogBYo9zpI9G8fD2s-i0FUfu8sg1HbHvpoikkwS3I8zh8kkuj-_8F2IZNP
AUTH0_CALLBACK_URL=https://hulta-foaltracker.app/callback

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# JWT Configuration
JWT_SECRET=hL8n7X2mK9pQ4vR3tY6uJ1cF5bN8mW0aS4dG7hK9lP2nX5vC8qE3wT6yU0iO4pA7sD9fG2hJ4kL7mN0bV5xC8zQ1wE3rT6yU

# Deployment Specific
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=horse_tracking_db
POSTGRES_USER=horsetracker
POSTGRES_PASSWORD=R2,Y@B&wO46.Ln}Q

# WebSocket Configuration
WEBSOCKET_HOST=localhost
WEBSOCKET_PORT=8081
WEBSOCKET_PATH=/notifications
EOF

    # Create frontend .env file
    cat <<EOF > "$APP_DIR/frontend-react/.env"
# Auth0 Configuration
VITE_AUTH0_DOMAIN=dev-r083cwkcv0pgz20x.eu.auth0.com
VITE_AUTH0_CLIENT_ID=OBmEol1z4U49r3YI3priDdGbvF5i4O7d
VITE_AUTH0_AUDIENCE=https://api.hulta-foaltracker.app

# API Configuration
VITE_API_URL=https://api.hulta-foaltracker.app

# WebSocket Connection
VITE_WEBSOCKET_URL=wss://api.hulta-foaltracker.app/notifications
EOF

    # Create a combined .env file for Docker Compose at the root
    cat <<EOF > "$APP_DIR/.env"
# Database Configuration
POSTGRES_DB=horse_tracking_db
POSTGRES_USER=horsetracker
POSTGRES_PASSWORD=R2,Y@B&wO46.Ln}Q
DATABASE_URL=postgresql://horsetracker:R2,Y@B&wO46.Ln}Q@postgres:5432/horse_tracking_db

# JWT Configuration
JWT_SECRET=hL8n7X2mK9pQ4vR3tY6uJ1cF5bN8mW0aS4dG7hK9lP2nX5vC8qE3wT6yU

# Auth0 Configuration
AUTH0_DOMAIN=dev-r083cwkcv0pgz20x.eu.auth0.com
AUTH0_AUDIENCE=https://api.hulta-foaltracker.app
AUTH0_CLIENT_ID=OBmEol1z4U49r3YI3priDdGbvF5i4O7d
AUTH0_CLIENT_SECRET=yTYWEdEogBYo9zpI9G8fD2s-i0FUfu8sg1HbHvpoikkwS3I8zh8kkuj-_8F2IZNP

# API Configuration
API_URL=https://api.hulta-foaltracker.app
WEBSOCKET_URL=wss://api.hulta-foaltracker.app/ws
EOF

    # Ensure correct permissions
    chmod 600 "$APP_DIR/.env" "$APP_DIR/backend/.env" "$APP_DIR/frontend-react/.env"
    
    log "✅ Environment setup complete with hardcoded values"
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

# Update setup_ssh_keys() with better error handling
setup_ssh_keys() {
    log "🔑 Checking SSH configuration"

    # Start ssh-agent and ensure it's running
    eval "$(ssh-agent -s)"
    
    # Use the existing deploy user's SSH key
    key="/home/deploy/.ssh/id_ed25519"
    
    if [ ! -f "$key" ]; then
        error_exit "SSH key not found at $key"
    fi

    # Add the key to the agent
    if ! ssh-add "$key" 2>/dev/null; then
        log "⚠️ Warning: Could not add SSH key to agent, trying alternative method..."
        export GIT_SSH_COMMAND="ssh -i $key -o StrictHostKeyChecking=accept-new -o IdentitiesOnly=yes"
    fi

    # Display the key fingerprint for verification
    log "🔑 SSH key fingerprint: $(ssh-keygen -lf $key)"

    # Test GitHub SSH connection (ignoring exit code as it always returns 1)
    log "🔍 Testing GitHub SSH connection..."
    if ! ssh -T -o StrictHostKeyChecking=accept-new git@github.com 2>&1 | grep -q "successfully authenticated"; then
        error_exit "SSH authentication to GitHub failed. Please ensure your SSH key is added to GitHub"
    fi

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
        log "📥 Cloning repository using SSH"
        # Test SSH connection first
        if ! ssh -T git@github.com 2>&1 | grep -q "successfully authenticated"; then
            error_exit "SSH authentication to GitHub failed. Please ensure your SSH key is added to GitHub"
        fi
        
        git clone -b $BRANCH $GIT_REPO $APP_DIR || error_exit "Failed to clone repository"
        cd $APP_DIR || error_exit "Failed to change to app directory"
    fi
    
    log "✅ Repository setup complete"
    setup_environment
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

    # Check for environment files
    if [ ! -f "$APP_DIR/.env.frontend" ]; then
        error_exit ".env.frontend not found in $APP_DIR"
    fi

    if [ ! -f "$APP_DIR/.env.backend" ]; then
        error_exit ".env.backend not found in $APP_DIR"
    fi
    
    log "✅ Project structure verified"
}

# ...existing code...

# Update start_services to handle separate env files
start_services() {
    log "🐳 Starting Docker services"
    
    # Verify project structure first
    verify_project_structure
    
    # Ensure Docker is running
    sudo systemctl start docker

    # Create necessary directories if they don't exist
    mkdir -p "$APP_DIR/backend"
    mkdir -p "$APP_DIR/frontend-react"
    
    # Handle backend environment file
    if [ -f "$APP_DIR/.env.backend" ]; then
        cp "$APP_DIR/.env.backend" "$APP_DIR/backend/.env"
        # Also copy to root for docker-compose
        cp "$APP_DIR/.env.backend" "$APP_DIR/.env"
    fi
    
    # Handle frontend environment file
    if [ -f "$APP_DIR/.env.frontend" ]; then
        cp "$APP_DIR/.env.frontend" "$APP_DIR/frontend-react/.env"
    fi
    
    # Export required environment variables for docker-compose
    if [ -f "$APP_DIR/.env" ]; then
        export DATABASE_URL=$(grep '^DATABASE_URL=' "$APP_DIR/.env" | cut -d '=' -f2 || echo "")
        export JWT_SECRET=$(grep '^JWT_SECRET=' "$APP_DIR/.env" | cut -d '=' -f2 || echo "")
        export POSTGRES_DB=$(grep '^POSTGRES_DB=' "$APP_DIR/.env" | cut -d '=' -f2 || echo "")
        export POSTGRES_USER=$(grep '^POSTGRES_USER=' "$APP_DIR/.env" | cut -d '=' -f2 || echo "")
        export POSTGRES_PASSWORD=$(grep '^POSTGRES_PASSWORD=' "$APP_DIR/.env" | cut -d '=' -f2 || echo "")
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

# Main execution with improved error handling
main() {
    log "🚀 Starting Hulta Pregnancy App Full Deployment"
    
    # Create APP_DIR if it doesn't exist
    mkdir -p "$APP_DIR" || error_exit "Failed to create application directory"
    
    validate_environment || error_exit "Environment validation failed"
    prepare_system || error_exit "System preparation failed"
    cleanup_system || error_exit "System cleanup failed"
    configure_ssh || error_exit "SSH configuration failed"
    setup_ssh_keys || error_exit "SSH key setup failed"
    configure_firewall || error_exit "Firewall configuration failed"
    install_docker || error_exit "Docker installation failed"
    install_nodejs || error_exit "Node.js installation failed"
    install_security_tools || error_exit "Security tools installation failed"
    setup_repository || error_exit "Repository setup failed"
    setup_environment || error_exit "Environment setup failed"
    verify_project_structure || error_exit "Project structure verification failed"
    validate_docker_compose || error_exit "Docker Compose validation failed"
    setup_nginx || error_exit "Nginx setup failed"
    start_services || error_exit "Service startup failed"
    setup_ssl || error_exit "SSL setup failed"
    verify_deployment || error_exit "Deployment verification failed"

    log "🎉🎉🎉 Deployment completed successfully! 🎉🎉🎉"
    log "🌐 Application is now live at https://hulta-foaltracker.app"
    log "⚠️ Important: Make sure to add the displayed SSH public key to GitHub!"
}

# Run main function
main

exit 0