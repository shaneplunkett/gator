package main

import (
	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
	"github.com/shaneplunkett/gator/internal/config"
	"github.com/shaneplunkett/gator/internal/database"
)

type command struct {
	name       string
	arguements []string
}

type state struct {
	db     *database.Queries
	config *config.Config
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.list[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exist := c.list[cmd.name]
	if !exist {
		log.Fatalf("Unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if cmd.arguements == nil {
		log.Fatalf("Username arguement required")
	}
	if len(cmd.arguements) > 1 {
		log.Fatalf("Login only accepts one arguement")
	}
	err := s.config.SetUser(cmd.arguements[0])
	if err != nil {
		return err
	}
	log.Infof("User has been set to: %s\n", cmd.arguements[0])

	return nil
}

func handlerRegister(s *state, cmd command) error
