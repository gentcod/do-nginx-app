#!/bin/bash
set -euo pipefail

echo "Initializing Nginx HTTP Server"
export IP_ADDR=$1
export ENV=$2
export STARTUP_SCRIPT=$3
export API_PORT=$4
export GITHUB_REPO=$5

HOST=$(echo "$IP_ADDR" | grep -E '^([0-9]{1,3}\.){3}[0-9]{1,3}$')
if [ -z "$HOST" ]; then
   echo "Invalid IP address format"
   exit 1
fi

create_env_file() {
   local env_content="$ENV"
   IFS=',' read -ra ENV_VARS <<< "$env_content"
   > .env
   for var in "${ENV_VARS[@]}"; do
      if [[ "$var" =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then
         echo "$var" >> .env
      else
         echo "Warning: Skipping invalid env var: $var"
      fi
   done
}

# Update and Install dependencies
sudo apt-get update -qq
sudo apt-get upgrade -y
sudo apt-get install -y --no-install-recommends \
    curl nodejs nginx vim nano zip npm

# Setup your app files
sudo mkdir /var/www/html/api
cd /var/www/html/api || exit 1

if ! git clone "$GITHUB_REPO"; then
   echo "Error: Failed to clone repository"
   exit 1
fi

# Create environment file
create_env_file "$ENV"

# Run startup script with error handling
if ! bash "$STARTUP_SCRIPT"; then
   echo "Startup script failed"
   exit 1
fi

# Firewall configuration -> Enable firewall rules
sudo ufw allow 'OpenSSH'
sudo ufw allow 'Nginx HTTP'
sudo ufw --force enable

echo "Firewall rules configured:"
sudo ufw status

# Setup Nginx configuration with template
cd ~
nginx_config="/etc/nginx/sites-available/api"
sudo touch /etc/nginx/sites-available/api
echo "server {
   server_name $HOST;

   listen 80;
   listen [::]:80;

   # SSL configuration
   # listen 443 ssl
   # listen [::]:443 ssl

   location / {
      proxy_pass http://api:$API_PORT;
      # proxy_set_header Connection 'upgrade';
      # proxy_set_header Host $host;
      proxy_http_version 1.1;
      # proxy_cache_bypass $http_upgrade;
      # try_files $uri $uri/ =404;
   }
}" | sudo tee "$nginx_config"

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