#!/bin/sh
echo "Hello there!!!... initializing your Node app with Nginx HTTP Proxy"
export IP_ADDR=${{ secrets.IP_ADDR }}
export GITHUB_REPO=${{ secrets.GITHUB_REPO }}
export ENV=${{ secrets.ENV }}
export STARTUP_SCRIPT=${{ secrets.STARTUP_SCRIPT }}
export API_PORT=${{ secrets.API_PORT }}

# Install dependencies
sudo apt update -y
sudo apt upgrade -y
sudo apt-get install -y curl
sudo apt-get install -y nodejs
sudo apt install -y nginx vim nano zip npm

# Setup your app files
sudo mkdir /var/www/html/api
cd /var/www/html/api
git clone $GITHUB_REPO
echo $ENV > .env
$STARTUP_SCRIPT

# Enable firewall rules
sudo ufw allow 'OpenSSH'
sudo ufw allow 'Nginx HTTP'
sudo ufw enable

# Setup Nginx Proxy Config
cd ~
sudo touch /etc/nginx/sites-available/api
echo "server {
   server_name $IP_ADDR;

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
}" > api

# Establish symbolic link
cd ~
sudo ln -s /etc/nginx/sites-available/api /etc/nginx/sites-enabled/api

# Check if nginx is properly setup
sudo nginx -t
sudo systemctl restart nginx

echo "Setup complete!!!.... Verify all configurations work fine"