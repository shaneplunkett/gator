package main

import (
	"github.com/charmbracelet/log"
	"github.com/shaneplunkett/gator/internal/config"
	"os"
)

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("Error Reading Config: %v", err)
	}
	s := &state{config: config}
	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]")
	}

	cmdarg := args[0]
	argList := args[1:]
	comm := command{name: cmdarg, arguements: argList}
	err = cmds.run(s, comm)
	if err != nil {
		log.Fatal(err)
	}
}
