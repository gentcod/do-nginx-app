#!/bin/bash
set -euo pipefail

echo "Hello there!!!... Updating your Node app from Nginx HTTP to HTTPS Proxy"

if [ $# -ne 4 ]; then
   echo "Usage: $0 <IP_ADDR> <ENV> <STARTUP_SCRIPT> <API_PORT> <GITHUB_REPO>"
   exit 1
fi

export DOMAIN_NAME=$1
export API_PORT=$2

if ! [[ "$API_PORT" =~ ^[0-9]+$ ]]; then
   echo "Invalid API port: $API_PORT"
   exit 1
fi

if ! [[ "$DOMAIN_NAME" =~ ^[a-zA-Z0-9.-]+$ ]]; then
   echo "Invalid domain name: $DOMAIN_NAME"
   exit 1
fi

# Enable firewall rules
sudo ufw allow 'Nginx FULL'

# Setup Nginx Proxy Config for HTTPS
cd ~
nginx_config="/etc/nginx/sites-available/api"
echo "server {
   server_name $DOMAIN_NAME;

   # SSL configuration
   listen 443 ssl
   listen [::]:443 ssl

   location / {
      proxy_pass http://api:$API_PORT;
      # proxy_set_header Connection 'upgrade';
      # proxy_set_header Host $host;
      proxy_http_version 1.1;
      # proxy_cache_bypass $http_upgrade;
      # try_files $uri $uri/ =404;
   }
}" | sudo tee "$nginx_config"

sudo ln -sf "$nginx_config" /etc/nginx/sites-enabled/api

# Enable SSL certification using Cerbot
if ! command -v snap &> /dev/null; then
   echo "Snap is not installed. Please install snap before proceeding."
   exit 1
fi

sudo snap install --classic certbot

for cmd in snap certbot nginx; do
   if ! command -v $cmd &> /dev/null; then
      echo "Error: $cmd is not installed."
      exit 1
   fi
done

sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot --nginx

# Check if nginx is properly setup
if ! sudo nginx -t; then
   echo "Nginx configuration test failed"
   exit 1
fi

# Disable HTTP
sudo ufw delete allow 'Nginx HTTP'
sudo systemctl restart nginx

echo "Setup complete!!!.... Verify all configurations work fine"