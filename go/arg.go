package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ProgArgs struct {
	HostAddr      string
	Protocol      string
	AuthType      string
	Port          int
	User          string
	Password      string
	PKey          string
	Passphrase    string
	GitHubRepo    string
	StartupScript string
	ApiPort       string
	Env           string
}

func GetProgArgs() (*ProgArgs, error) {
	port := 22
	if portStr := os.Getenv("INPUT_PORT"); portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %w", err)
		}
		port = p
	}

	protocol := "tcp"
	if p := os.Getenv("INPUT_PROTOCOL"); p != "" {
		protocol = p
	}

	config := &ProgArgs{
		HostAddr:      os.Getenv("INPUT_HOST"),
		Protocol:      protocol,
		Port:          port,
		User:          os.Getenv("INPUT_USERNAME"),
		Password:      os.Getenv("INPUT_PASSWORD"),
		PKey:          os.Getenv("INPUT_PKEY"),
		Passphrase:    os.Getenv("INPUT_PASSPHRASE"),
		GitHubRepo:    os.Getenv("INPUT_GITHUB_REPO"),
		StartupScript: os.Getenv("INPUT_STARTUP_SCRIPT"),
		ApiPort:       os.Getenv("INPUT_API_PORT"),
		Env:           os.Getenv("INPUT_ENV"),
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, err
	}

	config.getAuthType()

	return config, nil
}

func (args *ProgArgs) validate() error {
	if args.HostAddr == "" {
		return fmt.Errorf("host is required")
	}
	if args.User == "" {
		return fmt.Errorf("username is required")
	}
	if args.Port < 1 || args.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", args.Port)
	}
	if args.Password == "" && args.PKey == "" {
		return fmt.Errorf("one of password or ssh key is required")
	}
	if args.PKey == "" && args.Passphrase != "" {
		return fmt.Errorf("ssh key is required if passhrase is provided")
	}
	return nil
}

func (arg *ProgArgs) getAuthType() {
	if arg.PKey != "" && arg.Passphrase == "" {
		arg.AuthType = "private-key-only"
	}
	if arg.PKey != "" && arg.Passphrase != "" {
		arg.AuthType = "private-key-with-passphrase"
	}
	arg.AuthType = "password"
}

func ParseArgToScript(script []byte, args []string, updatedScriptPath string) error {
	updatedScript := string(script)
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i)
		updatedScript = strings.ReplaceAll(updatedScript, placeholder, arg)
	}

	err := os.WriteFile(updatedScriptPath, []byte(updatedScript), 0755)
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("Error writing updated script: %v\n", err)
	}

	return nil
}
