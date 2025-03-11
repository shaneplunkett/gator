package main

import (
	"fmt"
	"github.com/shaneplunkett/gator/internal/config"
)

type command struct {
	name       string
	arguements []string
}

type state struct {
	config *config.Config
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.list[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	command, exist := c.list[cmd.name]
	if !exist {
		return fmt.Errorf("Unknown command: %s", cmd.name)
	}
	return command(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if cmd.arguements == nil {
		return fmt.Errorf("Username arguement required")
	}
	if len(cmd.arguements) > 1 {
		return fmt.Errorf("Login only accepts one arguement")
	}
	err := s.config.SetUser(cmd.arguements[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %s\n", cmd.arguements[0])

	return nil
}
