package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	config, err := LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	args, err := GetProgArgs()
	if err != nil {
		log.Fatal("error getting environmental variables: ", err)
	}

	script, err := os.ReadFile(config.ScriptPath)
	if err != nil {
		log.Fatal("error reading script file", err)
	}

	err = ParseArgToScript(
		script, []string{
			args.HostAddr,
			args.Env,
			args.StartupScript,
			args.ApiPort,
			args.GitHubRepo,
		},
		config.UpdateScriptPath,
	)
	if err != nil {
		log.Fatal("error parsing args to script", err)
	}

	// Implement SSH
	opts := sshOpts{
		HostAddr: fmt.Sprintf("%s:%d", args.HostAddr, args.Port),
		Protocol: args.Protocol,
		Cmd:      "echo \"Hello there\"",
		AuthOpts: AuthOpts{
			Type:       args.AuthType,
			User:       args.User,
			HostKey:    false, // TODO: determine when to use hostkey validation
			Password:   args.Password,
			PrivateKey: []byte(args.PKey),
			Passphrase: []byte(args.Passphrase),
		},
	}

	client, err := CreateSSHClient(opts)
	if err != nil {
		log.Fatal("failed to create ssh Client: ", err)
	}
	defer client.Close()

	err = sshCopy(client, config.UpdateScriptPath, config.RemoteFilePath)
	if err != nil {
		log.Fatal("failed to copy script: ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b

	if err := session.Run(fmt.Sprintf("bash %s", config.RemoteFilePath)); err != nil {
		log.Fatal("failed to run: " + err.Error())
	}
	fmt.Println("Script output: ....")
	fmt.Println(b.String())

	err = sshCleanup(client, config.RemoteFilePath)
	if err != nil {
		log.Fatal("failed to cleanup copied script: ", err)
	}

	err = os.Remove(config.UpdateScriptPath)
	if err != nil {
		log.Fatal("Error deleting file: ", err)
	}
}

// TODO: write testcases for key funcs
// TODO: Implement updating server from http to https: Param: http, https. 