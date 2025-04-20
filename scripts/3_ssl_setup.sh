#!/bin/bash
# Hulta Pregnancy App - SSL Setup Script
# This script handles SSL configuration with Let's Encrypt
# Version: 1.0

# Enhanced error handling and logging
set -euo pipefail

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [SSL] $*" | tee -a ~/hulta_ssl.log
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
DOMAIN=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -d|--domain)
            DOMAIN="$2"
            shift
            shift
            ;;
        *)
            error_exit "Unknown option: $1"
            ;;
    esac
done

# Check if domain is provided
if [ -z "$DOMAIN" ]; then
    error_exit "Domain name is required. Usage: $0 -d yourdomain.com"
fi

# Check if we're running as root
if [ "$(id -u)" = "0" ]; then
    log "WARNING: This script is running as root. Some commands might behave differently."
fi

# Install Certbot
install_certbot() {
    log "📦 Installing Certbot"
    
    if command -v certbot &> /dev/null; then
        log "Certbot is already installed: $(certbot --version)"
    else
        apt-get update
        apt-get install -y certbot python3-certbot-nginx || error_exit "Failed to install Certbot"
    fi
    
    log "✅ Certbot installation complete"
}

# Update Nginx configuration with domain
update_nginx_config() {
    log "🌐 Updating Nginx configuration with domain"
    
    # Update server_name in Nginx configuration
    sed -i "s/server_name _;/server_name $DOMAIN;/" /etc/nginx/sites-available/hulta-app || error_exit "Failed to update Nginx configuration"
    
    # Test Nginx configuration
    nginx -t || error_exit "Nginx configuration test failed"
    
    # Reload Nginx
    systemctl reload nginx || error_exit "Failed to reload Nginx"
    
    log "✅ Nginx configuration updated"
}

# Obtain SSL certificate
obtain_ssl_certificate() {
    log "🔒 Obtaining SSL certificate for $DOMAIN"
    
    # Run Certbot
    certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos --email "admin@$DOMAIN" --redirect || error_exit "Failed to obtain SSL certificate"
    
    log "✅ SSL certificate obtained successfully"
}

# Update frontend environment variables
update_frontend_env() {
    log "🔧 Updating frontend environment variables"
    
    # Update API URL in frontend .env file
    if [ -f "$APP_DIR/frontend-react/.env" ]; then
        sed -i "s|VITE_API_URL=.*|VITE_API_URL=https://$DOMAIN/api|" "$APP_DIR/frontend-react/.env" || error_exit "Failed to update frontend .env file"
        
        log "✅ Frontend environment variables updated"
        log "⚠️ You need to rebuild the frontend for changes to take effect"
    else
        log "⚠️ Frontend .env file not found. Skipping update."
    fi
}

# Rebuild and restart services
rebuild_services() {
    log "🚀 Rebuilding and restarting services"
    
    cd "$APP_DIR" || error_exit "Failed to change to app directory"
    
    # Rebuild and restart containers
    docker compose build frontend || error_exit "Failed to rebuild frontend"
    docker compose up -d || error_exit "Failed to restart services"
    
    log "✅ Services rebuilt and restarted"
}

# Verify SSL setup
verify_ssl() {
    log "🔍 Verifying SSL setup"
    
    # Check if HTTPS is working
    if curl -s -o /dev/null -w "%{http_code}" "https://$DOMAIN" | grep -q "200\|301\|302"; then
        log "HTTPS is working correctly"
    else
        log "⚠️ HTTPS check failed. Please verify manually."
    fi
    
    log "✅ SSL verification complete"
}

# Main function
main() {
    log "🚀 Starting SSL setup for $DOMAIN"
    
    install_certbot
    update_nginx_config
    obtain_ssl_certificate
    update_frontend_env
    rebuild_services
    verify_ssl
    
    log "🎉 SSL setup completed successfully!"
    log "Your application is now available at https://$DOMAIN"
}

# Run main function
main

exit 0
