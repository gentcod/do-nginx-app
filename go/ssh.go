package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type sshOpts struct {
	HostAddr string
	Protocol string
	Cmd      string
	AuthOpts AuthOpts
}

func CreateSSHClient(opts sshOpts) (*ssh.Client, error) {
	config, err := initClientConfig(opts.AuthOpts)
	if err != nil || config == nil {
		return nil, fmt.Errorf("error initializing ssh client config: %v", err)
	}

	client, err := ssh.Dial(opts.Protocol, opts.HostAddr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return client, nil
}

type AuthOpts struct {
	// Type specifies the AuthMethod type. It has the following options:
	// password | private-key-only | private-key-with-passphrase
	Type string

	// User specifies the ssh user as in root@127.0.0.1:22. With user -> root
	User string

	// HostKey specifies if host key check is to be implemented.
	HostKey bool

	// Password repreents the password is AuthOpts.Type -> password
	Password string

	// PrivateKey represents the ssh private key is AuthOpts.Type -> private-key-only or private-key-with-passphrase
	PrivateKey []byte

	// Passphrase represents the associated passphrase for a ssh private key
	// it is used only when AuthOpts.Type -> private-key-with-passphrase
	Passphrase []byte
}

func initClientConfig(opts AuthOpts) (*ssh.ClientConfig, error) {
	var callBack ssh.HostKeyCallback
	var hostkey ssh.PublicKey

	if opts.HostKey {
		callBack = ssh.FixedHostKey(hostkey)
	} else {
		callBack = ssh.InsecureIgnoreHostKey()
	}

	if opts.Type == "password" {
		return &ssh.ClientConfig{
			User: opts.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(opts.Password),
			},
			HostKeyCallback: callBack,
		}, nil
	}

	if opts.Type == "private-key-only" {
		signer, err := ssh.ParsePrivateKey(opts.PrivateKey)
		if err != nil {
			return nil, err
		}
		return &ssh.ClientConfig{
			User: opts.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: callBack,
		}, nil
	}

	if opts.Type == "private-key-with-passphrase" {
		signer, err := ssh.ParsePrivateKeyWithPassphrase(opts.PrivateKey, opts.Passphrase)
		if err != nil {
			return nil, err
		}
		return &ssh.ClientConfig{
			User: opts.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: callBack,
		}, nil
	}

	return nil, nil
}

func sshCopy(client *ssh.Client, scriptFilePath, remoteFilePath string) error {
	scpSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SCP session: %v", err)
	}
	defer scpSession.Close()

	script, err := os.ReadFile(scriptFilePath)
	if err != nil {
		return fmt.Errorf("failed to read script file: %v", err)
	}

	scpStdin, err := scpSession.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get SCP stdin: %v", err)
	}

	go func() {
		defer scpStdin.Close()
		fmt.Fprintf(scpStdin, "C0755 %d script.sh\n", len(script))
		scpStdin.Write(script)
		fmt.Fprint(scpStdin, "\x00")
	}()

	if err := scpSession.Run(fmt.Sprintf("scp -tr %s", remoteFilePath)); err != nil {
		return fmt.Errorf("failed to upload script: %v", err)
	}

	fmt.Println("Script uploaded successfully.")

	return nil
}

func sshCleanup(client *ssh.Client, remoteFilePath string) error {
	cleanupSession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create cleanup session: %v", err)
	}
	defer cleanupSession.Close()
	if err := cleanupSession.Run(fmt.Sprintf("rm %s", remoteFilePath)); err != nil {
		return fmt.Errorf("failed to clean up remote script: %v", err)
	} else {
		fmt.Println("remote script cleaned up successfully.")
	}

	return nil
}
