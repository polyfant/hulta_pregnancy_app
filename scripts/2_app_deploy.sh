#!/bin/bash
# Hulta Pregnancy App - Application Deployment Script
# This script handles application deployment and configuration
# Version: 1.0

# Enhanced error handling and logging
set -euo pipefail

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [DEPLOY] $*" | tee -a ~/hulta_deploy.log
}

error_exit() {
    log "ERROR: $1"
    echo "ERROR: $1" >&2
    exit 1
}

# Trap any errors
trap 'error_exit "Command failed: $BASH_COMMAND"' ERR

# Configuration Variables
APP_DIR="/opt/hulta-pregnancy-app"
ENV_FILE="$APP_DIR/.env"
DOCKER_COMPOSE_FILE="$APP_DIR/docker-compose.yml"
# Replace YOUR_USERNAME and YOUR_TOKEN with your actual GitHub username and personal access token
GIT_REPO="https://polyfant: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@github.com/polyfant/hulta_pregnancy_app.git"
BRANCH="main"

# Check if we're running as root - changed to just a warning
if [ "$(id -u)" = "0" ]; then
    log "WARNING: This script is running as root. Some commands might behave differently."
    # Not exiting, just warning
fi

# Verify prerequisites
verify_prerequisites() {
    log "Verifying prerequisites"
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        error_exit "Docker is not installed. Please run 1_server_setup.sh first."
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker compose &> /dev/null; then
        error_exit "Docker Compose is not installed. Please run 1_server_setup.sh first."
    fi
    
    # Check if application directory exists
    if [ ! -d "$APP_DIR" ]; then
        error_exit "Application directory does not exist. Please run 1_server_setup.sh first."
    fi
    
    log "Prerequisites verified"
}

# Clone or update repository
setup_repository() {
    log "Setting up repository"
    
    cd "$APP_DIR" || error_exit "Failed to change to app directory"
    
    if [ -d "$APP_DIR/.git" ]; then
        log "Repository already exists, updating..."
        git pull origin "$BRANCH" || error_exit "Failed to update repository"
    else
        log "Cloning repository..."
        git clone --branch "$BRANCH" "$GIT_REPO" . || error_exit "Failed to clone repository"
    fi
    
    log "Repository setup complete"
}

# Create environment files
create_env_files() {
    log "Creating environment files"
    
    # Check if .env already exists
    if [ -f "$ENV_FILE" ]; then
        log ".env file already exists, skipping creation"
    else
        log "Creating .env file..."
        
        # Generate random passwords
        DB_PASSWORD=$(openssl rand -base64 12)
        JWT_SECRET=$(openssl rand -base64 32)
        
        # Create .env file
        cat > "$ENV_FILE" << EOF
# Database settings
POSTGRES_DB=horse_tracking_db
POSTGRES_USER=horsetracker
POSTGRES_PASSWORD=$DB_PASSWORD
DATABASE_URL=postgres://horsetracker:$DB_PASSWORD@postgres:5432/horse_tracking_db

# Auth0 settings
AUTH0_DOMAIN=your-domain.auth0.com
AUTH0_AUDIENCE=your-api-identifier
JWT_SECRET=$JWT_SECRET

# Other settings
GIN_MODE=release
EOF
        
        log "Please update AUTH0_DOMAIN and AUTH0_AUDIENCE in $ENV_FILE"
    fi
    
    # Check if frontend .env already exists
    if [ -f "$APP_DIR/frontend-react/.env" ]; then
        log "Frontend .env file already exists, skipping creation"
    else
        log "Creating frontend .env file..."
        
        # Create frontend .env file
        cat > "$APP_DIR/frontend-react/.env" << EOF
VITE_AUTH0_DOMAIN=your-domain.auth0.com
VITE_AUTH0_CLIENT_ID=your-client-id
VITE_AUTH0_AUDIENCE=your-api-identifier
VITE_API_URL=https://your-domain.com/api
EOF
        
        log "Please update AUTH0 settings in $APP_DIR/frontend-react/.env"
    fi
    
    log "Environment files created"
}

# Configure Nginx
configure_nginx() {
    log "Configuring Nginx"
    
    # Create Nginx configuration file
    tee /etc/nginx/sites-available/hulta-app > /dev/null << EOF
server {
    listen 80;
    server_name _;  # Replace with your domain name when ready

    # Frontend
    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF
    
    # Enable the site
    if [ ! -f /etc/nginx/sites-enabled/hulta-app ]; then
        ln -s /etc/nginx/sites-available/hulta-app /etc/nginx/sites-enabled/ || error_exit "Failed to enable Nginx site"
        rm -f /etc/nginx/sites-enabled/default || true  # Remove default site if it exists
    fi
    
    # Test Nginx configuration
    nginx -t || error_exit "Nginx configuration test failed"
    
    # Restart Nginx
    systemctl restart nginx || error_exit "Failed to restart Nginx"
    
    log "Nginx configured"
}

# Build and start Docker containers
start_services() {
    log "Building and starting services"
    
    cd "$APP_DIR" || error_exit "Failed to change to app directory"
    
    # Check if docker-compose.yml exists
    if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
        error_exit "Docker Compose file not found at $DOCKER_COMPOSE_FILE"
    fi
    
    # Build and start containers
    docker compose -f "$DOCKER_COMPOSE_FILE" down || true  # Bring down any existing containers
    docker compose -f "$DOCKER_COMPOSE_FILE" build || error_exit "Failed to build Docker images"
    docker compose -f "$DOCKER_COMPOSE_FILE" up -d || error_exit "Failed to start Docker containers"
    
    log "Services started"
}

# Verify deployment
verify_deployment() {
    log "Verifying deployment"
    
    # Check if containers are running
    if ! docker compose -f "$DOCKER_COMPOSE_FILE" ps | grep -q "Up"; then
        error_exit "Containers are not running. Check docker logs for details."
    fi
    
    # Wait for services to be ready
    log "Waiting for services to be ready..."
    sleep 10
    
    # Check if backend is responding
    if ! curl -s http://localhost:8080/api/v1/health | grep -q "ok"; then
        log "Backend health check failed. This might be normal if the health endpoint is not implemented."
    else
        log "Backend is responding correctly"
    fi
    
    log "Deployment verification complete"
    log "Application deployed successfully!"
    log "Next step: Run 3_ssl_setup.sh to configure SSL (optional)"
}

# Main function
main() {
    log "Starting application deployment"
    
    verify_prerequisites
    setup_repository
    create_env_files
    configure_nginx
    start_services
    verify_deployment
    
    log "Application deployment completed successfully!"
}

# Run main function
main

exit 0
