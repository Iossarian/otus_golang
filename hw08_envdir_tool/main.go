package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	env, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}
	RunCmd(args[2:], env)
}
