package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("Hello, Initializing your Nginx node app")
	fmt.Println("...................")

	cmd := exec.Command("bash", "./scripts/test.sh")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}