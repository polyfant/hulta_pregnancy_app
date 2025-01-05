#!/bin/bash

# Hulta Equestrian Server Setup Script for Surface Pro 3
# Version: 1.0
# Designed for Ubuntu Server LTS

set -e  # Exit immediately if a command exits with a non-zero status
set -o pipefail  # Fail on any part of a pipe

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*"
}

# Check for root/sudo
if [[ $EUID -ne 0 ]]; then
   log "This script must be run with sudo" 
   exit 1
fi

# Update and Upgrade
log "Updating system packages..."
apt update && apt upgrade -y

# Install essential utilities
log "Installing essential utilities..."
apt install -y \
    curl \
    wget \
    git \
    htop \
    glances \
    net-tools \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release

# Optimize Swap and Memory
log "Configuring swap space..."
# Remove existing swapfile if it exists
swapoff -a
rm -f /swapfile

# Create new swapfile (8GB for 4GB RAM)
fallocate -l 8G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile

# Permanent swap configuration
if ! grep -q '/swapfile' /etc/fstab; then
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
fi

# Kernel and Performance Tuning
log "Tuning kernel parameters..."
cat << EOF >> /etc/sysctl.conf
# Surface Pro 3 Optimizations
vm.swappiness=60
vm.overcommit_memory=1
vm.dirty_ratio=10
vm.dirty_background_ratio=5
EOF

sysctl -p

# Install Go
log "Installing Go..."
GO_VERSION="1.23.4"
wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH
cat << EOF >> /etc/profile.d/golang.sh
export GOROOT=/usr/local/go
export GOPATH=\$HOME/go
export PATH=\$PATH:\$GOROOT/bin:\$GOPATH/bin
EOF

# Install Docker
log "Installing Docker..."
# Remove any existing Docker installations
apt-get remove -y docker docker-engine docker.io containerd runc

# Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add current user to docker group
USER=$(logname)
usermod -aG docker $USER

# Install PostgreSQL Client
log "Installing PostgreSQL client..."
apt install -y postgresql-client

# Security Enhancements
log "Configuring SSH and security..."
# Harden SSH
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# Install and configure Fail2Ban
apt install -y fail2ban
cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sed -i 's/bantime  = 10m/bantime  = 1h/' /etc/fail2ban/jail.local
systemctl enable fail2ban
systemctl start fail2ban

# Unattended Security Updates
log "Configuring automatic security updates..."
apt install -y unattended-upgrades
dpkg-reconfigure -f noninteractive unattended-upgrades

# Firewall Configuration
log "Configuring UFW firewall..."
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https
ufw --force enable

# Create project directory
log "Creating project directory..."
mkdir -p /opt/hulta_equestrian
chown $USER:$USER /opt/hulta_equestrian

# Final system cleanup
log "Performing system cleanup..."
apt autoremove -y
apt autoclean

log "Setup complete! Please reboot the system."

# Optional: Reboot
read -p "Do you want to reboot now? (y/n) " response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    reboot
fi
