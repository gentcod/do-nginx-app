# do-nginx-app

An action/service helps you initialize a Nginx proxied web server on a Virtual Machine by spinning up the setups for needed dependencies.

It carries out the process in phases: 
   - File/script copy.
   - Script execution.
   - File/script cleanup.

Making sure unnecessary residuals are not persisted.

> Currently supports Debian and Ubuntu based images.

> Currently only supports NodeJS servers.

### Prequisites

For best experience: 

- Ensure you have a sudo user created with administrative auth. If you do not have that implemented; run the following scripts on your remote machine to create a user that will be used.
```bash
   # Create user and add to sudo group
   sudo useradd -m "username"
   
   # Set password (using chpasswd to avoid interactive prompt)
   echo "username:password" | sudo chpasswd
   
   # Add user to sudo group
   sudo usermod -aG sudo "username"
   
   # Set up NOPASSWD privileges
   echo "username ALL=(ALL) NOPASSWD: ALL" | sudo tee "/etc/sudoers.d/username"
   sudo chmod 0440 "/etc/sudoers.d/username"
```

- Ensure you have your ssh authorizations properly set up. If not follow the following steps on your local machine, follow the prompts when required:
```bash
   # Create ssh key
   ssh-keygen -t rsa -b 4096 -C "your_email@example.com"

   # Add key to Authorized Keys on remote/virtual machine
   ssh-copy-id -i my_key.pub username@vm-ip-address

   # Connect to virtual machine to verify and test key
   ssh -i my_key username@vm-ip-address

```

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
- name: Run Do-nginx
   uses: gentcod/do-nginx-app@v1
   with:
      host: ${{ secrets.HOST }}
      user: ${{ secrets.USER }}
      key: ${{ secrets.KEY }}
      passphrase: ${{ secrets.PASSPHRASE }}
      github-repo: ${{ secrets.GITHUB_REPO }}
      startup-script: ${{ secrets.STARTUP_SCRIPT }}
      api-port: ${{ secrets.API_PORT }}
```
