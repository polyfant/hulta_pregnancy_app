#!/bin/bash
# Diagnostic script to check deployment status

echo "=== System Status ==="
echo "Date: $(date)"
echo "Hostname: $(hostname)"
echo "Kernel: $(uname -r)"

echo -e "\n=== Directory Structure ==="
echo "Application directory exists: $(if [ -d /opt/hulta-pregnancy-app ]; then echo "Yes"; else echo "No"; fi)"
echo "Nginx sites available: $(ls -la /etc/nginx/sites-available/ 2>/dev/null || echo "Directory not found")"
echo "Nginx sites enabled: $(ls -la /etc/nginx/sites-enabled/ 2>/dev/null || echo "Directory not found")"

echo -e "\n=== Service Status ==="
echo "Docker status: $(systemctl is-active docker 2>/dev/null || echo "Not installed")"
echo "Nginx status: $(systemctl is-active nginx 2>/dev/null || echo "Not installed")"

echo -e "\n=== Docker Status ==="
echo "Docker containers:"
docker ps -a 2>/dev/null || echo "Docker not running or not installed"

echo -e "\n=== Network Status ==="
echo "Open ports:"
netstat -tulpn 2>/dev/null | grep LISTEN || echo "netstat not available"

echo -e "\n=== Log Files ==="
echo "Setup log (last 5 lines):"
tail -n 5 ~/hulta_setup.log 2>/dev/null || echo "Log file not found"
echo "Deploy log (last 5 lines):"
tail -n 5 ~/hulta_deploy.log 2>/dev/null || echo "Log file not found"
