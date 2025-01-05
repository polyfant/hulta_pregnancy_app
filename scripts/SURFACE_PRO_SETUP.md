# Surface Pro 3 Server Setup Script

## Prerequisites
- Fresh Ubuntu Server LTS installation
- Internet connection
- Sudo access

## Usage Instructions

1. After initial Ubuntu Server installation:
```bash
# Transfer the script to the server
scp surface_pro_setup.sh your_username@server_ip:~/

# Make the script executable
chmod +x surface_pro_setup.sh

# Run with sudo
sudo ./surface_pro_setup.sh
```

## What This Script Does
- Updates system packages
- Configures 8GB swap space
- Installs Go 1.23.4
- Installs Docker and Docker Compose
- Hardens SSH security
- Configures firewall (UFW)
- Sets up automatic security updates
- Installs monitoring tools

## Post-Installation Steps
1. Generate SSH keys
```bash
ssh-keygen -t ed25519
```

2. Copy your public key to the server
```bash
ssh-copy-id your_username@server_ip
```

3. Clone your project
```bash
cd /opt/hulta_equestrian
git clone [YOUR_REPO_URL]
```

## Customization
- Modify Go version in the script if needed
- Adjust swap size based on your specific Surface Pro 3 model

## Troubleshooting
- Ensure you have a stable internet connection
- Check log messages during script execution
- Reboot after installation if prompted

## Security Notes
- The script disables password login
- Root login is disabled
- Fail2Ban is configured to prevent brute-force attacks
