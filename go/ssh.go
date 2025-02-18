package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	file, err := os.Open(scriptFilePath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer file.Close()

	stdin, err := scpSession.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	defer stdin.Close()

	remoteDir := filepath.Dir(remoteFilePath)
	cmd := fmt.Sprintf("mkdir -p %s; cat > %s", remoteDir, remoteFilePath)
	if err := scpSession.Start(cmd); err != nil {
		return fmt.Errorf("failed to start remote command: %v", err)
	}

	if _, err := io.Copy(stdin, file); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	stdin.Close()

	if err := scpSession.Wait(); err != nil {
		return fmt.Errorf("remote command failed: %v", err)
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
	if err := cleanupSession.Run(fmt.Sprintf("rm -rf %s", filepath.Dir(remoteFilePath))); err != nil {
		return fmt.Errorf("failed to clean up remote script: %v", err)
	} else {
		fmt.Println("remote script cleaned up successfully.")
	}

	return nil
}

func sshExec(client *ssh.Client, remoteFilePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Create a pipe to capture stdout and stderr
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	combinedOutput := io.MultiReader(stdoutPipe, stderrPipe)

	var readErr error
	var b bytes.Buffer
	go func() {
		tee := io.TeeReader(combinedOutput, &b)
		_, err := io.Copy(os.Stdout, tee)
		if err != nil && err != io.EOF {
			readErr = fmt.Errorf("error reading from combined output: %v", err)
		}
	}()

	if readErr != nil {
		return readErr
	}

	if err := session.Run(fmt.Sprintf("bash %s", remoteFilePath)); err != nil {
		return fmt.Errorf("failed to run: %v", err)
	}

	fmt.Println("Script output: ....")
	fmt.Println(b.String())
	return nil
}
