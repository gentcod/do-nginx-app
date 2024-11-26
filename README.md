# do-nginx-app

A action helps you initialize a Nginx server on a Virtual Machine by simply setting up environment variables.

## Inputs

#### `host`
> *Required*: The IP Address of the Virtual Machine hosting your server.

#### `protocol`
> *Not Required*: SSH connection protocol. Default is set to "tcp".

#### `port`
> *Not Required*: The port for host connection. Default is set to -> 22.

#### `user`
> *Required*: The virtual machine user e.g root.

#### `password`
> *Not Required*: Password to authenticate ssh connection. NB: It is not required if authentication is done using SSH private key.

#### `key`
> *Not Required*: Authorized SSH key to authenticate ssh connection. NB: It is not required if authentication is not done with password.

#### `passphrase`
> *Not Required*: Associated passphrase if any, to the provided authorized SSH key. NB: It is not required if authentication is not done with password.

#### `github-repo`
> *Required*: The github repository with the server code you're trying to run.

#### `startup-script`
> *Required*: The application script to run your app. e.g npm start.

#### `api-port`
> *Required*: The configured PORT. e.g 5000

#### `env`
> *Not Required*: Environmental variables used to run your Node app. NB: variables will be populated in .env file on your virtual machine. Key-value pairs -> KEY=VALUE


## Usage
```yaml
- uses: gentcod/do-nginx-app@v2
- with:
   - host: ${{ secrets.HOST }}
   - user: ${{ secrets.USER }}
   - key: ${{ secrets.KEY }}
   - passphrase: ${{ secrets.PASSPHRASE }}
   - github-repo: ${{ secrets.GITHUB_REPO }}
   - startup-script: ${{ secrets.STARTUP_SCRIPT }}
   - api-port: ${{ secrets.API_PORT }}
```
