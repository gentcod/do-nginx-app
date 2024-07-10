# do-nginx-app

A action helps you initialize a Nginx server on a Virtual Machine by simply setting up environment variables.

## Inputs

#### `ip-address`
*Required*: The IP Address of the Virtual Machine hosting your server

#### `env`
*Not Required*: Environmental variables used to run your Node app. NB: variables will be populated in .env file on your virtual machine

#### `startup-script`
*Required*: The npm script to run your app. e.g npm start

#### `api-port`
*Required*: The configured PORT. e.g 5000

#### `user`
*Required*: The virtual machine user e.g root

#### `key`
*Required*: Authorized SSH key

## Usage
- uses: gentcod/do-nginx-app@v2
- with:
   - ip-address: secrets.IP_ADDR
   - env: secrets.ENV
   - startup-script: secrets.STARTUP_SCRIPT
   - api-port: secrets.API_PORT
   - user: secrets.USER
   - key: secrets.KEY`
