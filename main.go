package main

import (
	"fmt"
	"github.com/shaneplunkett/gator/internal/config"
	"os"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Printf("Error Reading Config: %s", err)
		os.Exit(1)
	}
	s := &state{config: config}
	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Printf("Not Enough Arguments provided: %s\n", args)
		os.Exit(1)
	}

	cmdarg := args[0]
	argList := args[1:]
	comm := command{name: cmdarg, arguements: argList}
	err = cmds.run(s, comm)
	if err != nil {
		fmt.Printf("Command Failed: %s\n", err)
	}
	os.Exit(0)
}
