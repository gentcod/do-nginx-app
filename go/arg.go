package main

import (
	"fmt"
	"os"
	"strings"
)

type ProgArgs struct {
	HostAddr      string
	StartupScript []byte
}

func GetProgArgs() (ProgArgs, error) {
	var progArgs ProgArgs
	if len(os.Args) > 3 {
		hostAddr := os.Args[1]
		env := os.Args[2]
		startupScript := os.Args[3]

		fmt.Printf("IP Address: %s\n", hostAddr)
		fmt.Printf("Environment: %s\n", env)
		fmt.Printf("Startup Script: %s\n", startupScript)
	} else {
		return progArgs, fmt.Errorf("insufficient arguments passed")
	}

	return progArgs, nil
}

func ParseArgToScript(script []byte, args []string, updatedScriptPath string) error {
	updatedScript := string(script)
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i) // Start replacing from $3
		updatedScript = strings.ReplaceAll(updatedScript, placeholder, arg)
	}

	// Write the updated script back to a file (or a new file)
	err := os.WriteFile(updatedScriptPath, []byte(updatedScript), 0755)
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("Error writing updated script: %v\n", err)
	}

	return nil
}