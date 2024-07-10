# do-nginx-app

A action helps you initialize a Nginx server on a Virtual Machine by simply setting up environment variables.

## Inputs

## `ip-address`

**Required** The IP Address of the Virtual Machine hosting your server

## `env`

**Not Required** Environmental variables used to run your Node app

## `startup-script`

**Required** The npm script to run your app. e.g npm start

## `ip-address`

**Required** The configured PORT. e.g 5000

uses: actions/do-nginx-app@v2
with:
   ip-address: ${{ secrets.IP_ADDR }}
   env=${{ secrets.ENV }}
   startup-script=${{ secrets.STARTUP_SCRIPT }}
   api-port=${{ secrets.API_PORT }}
