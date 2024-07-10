#!/bin/sh
echo "Hello there!!!... Updating your Node app from Nginx HTTP to HTTPS Proxy"
export DOMAIN_NAME=${{ secrets.DOMAIN_NAME }}

# Enable firewall rules
sudo ufw allow 'Nginx FULL'

# Setup Nginx Proxy Config for HTTPS
cd ~
cd /etc/nginx/sites-available
echo "server {
   server_name $DOMAIN_NAME;

   listen 80;
   listen [::]:80;

   # SSL configuration
   # listen 443 ssl
   # listen [::]:443 ssl

   location / {
      proxy_pass http://api:5000;
      # proxy_set_header Connection 'upgrade';
      # proxy_set_header Host $host;
      proxy_http_version 1.1;
      # proxy_cache_bypass $http_upgrade;
      # try_files $uri $uri/ =404;
   }
}" > api

# Enable SSL certification using Cerbot
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot --nginx

# Check if nginx is properly setup
sudo nginx -t

# Disable HTTP
sudo ufw delete allow 'Nginx HTTP'
sudo systemctl restart nginx

echo "Setup complete!!!.... Verify all configurations work fine"