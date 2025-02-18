#!/bin/bash

###############################
# Author: Oyefule Oluwatayo
# Date: 13/06/2024
#
# This script outputs the node health
#
# Version: v1
###############################

set -euo pipefail

echo "Initializing Nginx HTTP Server"

# Validate IP address
if [[ ! "$1" =~ ^([0-9]{1,3}\.){3}[0-9]{1,3}$ ]]; then
   echo "Error: Invalid IP address format"
   exit 1
fi

# Validate GitHub repository URL
if [[ ! "$5" =~ ^https://github.com/.* ]]; then
   echo "Error: Invalid GitHub repository URL"
   exit 1
fi

# Create environment file
create_env_file() {
   local env_content="$2"
   IFS=',' read -ra ENV_VARS <<< "$env_content"
   > .env
   for var in "${ENV_VARS[@]}"; do
      if [[ "$var" =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then
         echo "$var" >> .env
      else
         echo "Warning: Skipping invalid env var: $var"
      fi
   done
   chmod 600 .env
}

# Update and Install dependencies
sudo apt-get update -qq
sudo apt-get upgrade -y --with-new-pkgs
sudo apt-get install -y --no-install-recommends \
   curl nginx zip npm git

# Install Node.js using nvm
set +u
sudo curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
source ~/.profile
nvm --version || {
   echo "Failed to install Node Version Manager"
   exit 1
}

nvm install 22.1.0

# Verify Node installation
node -v || {
   echo "Failed to install Node.js"
   exit 1
}
set -u

# Setup app files
sudo mkdir -p /var/www/html/api
cd /var/www/html/api || exit 1

if ! sudo git clone $5 .; then
   echo "Error: Failed to clone repository"
   exit 1
fi

npm install || {
   echo "Failed to install package dependencies"
   exit 1
}

# Create environment file
# create_env_file "$2"

# Run startup script with error handling
if ! $3; then
   echo "Startup script failed"
   exit 1
fi

# Firewall configuration -> Enable firewall rules
sudo ufw allow 'Nginx HTTP'
sudo ufw --force enable

echo "Firewall rules configured:"
sudo ufw status

# Setup Nginx configuration with template
cd ~
nginx_config="/etc/nginx/sites-available/api"
sudo tee "$nginx_config" <<EOL
server {
   server_name $1;

   listen 80;
   listen [::]:80;

   # SSL configuration
   # listen 443 ssl
   # listen [::]:443 ssl

   location / {
      proxy_pass http://api:$4;
      # proxy_set_header Connection 'upgrade';
      # proxy_set_header Host $host;
      proxy_http_version 1.1;
      # proxy_cache_bypass $http_upgrade;
      # try_files $uri $uri/ =404;
   }
}
EOL

# Establish symbolic link
cd ~
sudo ln -sf $nginx_config /etc/nginx/sites-enabled/api

# Check if nginx is properly setup
if ! sudo nginx -t; then
   echo "Nginx configuration test failed"
   exit 1
fi

# Restart Nginx
sudo systemctl restart nginx

echo "Setup complete!!!.... Verify all configurations work fine"