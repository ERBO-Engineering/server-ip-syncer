#!/bin/bash

# Exit on error
set -e

echo "Installing IP Viber..."

# Create installation directory
sudo mkdir -p /opt/ip-viber

# Copy application files
sudo cp ip_viber /opt/ip-viber/
sudo cp config.json /opt/ip-viber/
sudo cp firestore.json /opt/ip-viber/

# Set permissions
sudo chmod 755 /opt/ip-viber/ip_viber
sudo chmod 644 /opt/ip-viber/config.json
sudo chmod 600 /opt/ip-viber/firestore.json

# Copy systemd service file
sudo cp ip-viber.service /etc/systemd/system/

# Reload systemd daemon
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable ip-viber
sudo systemctl start ip-viber

echo "Installation complete!"
echo "You can check the service status with: sudo systemctl status ip-viber"
echo "View logs with: sudo journalctl -u ip-viber -f" 