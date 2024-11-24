package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	// "os/exec"
)

func main() {
	config, err := LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	if len(os.Args) < 3 {
		fmt.Printf("Expected at least 5 arguments, got %d\n", len(os.Args))
		os.Exit(1)
	}

	script, err := os.ReadFile(config.ScriptPath)
	if err != nil {
		log.Fatal("error reading script file", err)
	}

	err = ParseArgToScript(script, os.Args, config.UpdateScriptPath)
	if err != nil {
		log.Fatal("error parsing args to script", err)
	}

	// DON'T DELETE
	// pKey, err := os.ReadFile("../secrets/id_rsa")
	// if err != nil {
		// 	log.Fatal("error encoutered reading env file: ", err)
		// }
		// passphrase := []byte("passphrase")
		
	// Implement SSH
	opts := sshOpts{
		HostAddr: os.Args[1],
		Cmd: "echo \"Hello there\"",
		AuthOpts: AuthOpts{
			Type: "private-key-with-passphrase",
			User: os.Args[2],
			HostKey: false,
			PrivateKey: []byte(os.Args[3]),
			Passphrase: []byte(os.Args[4]),
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

	// Delete create script
}
