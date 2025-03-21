package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/shaneplunkett/gator/internal/config"
	"github.com/shaneplunkett/gator/internal/database"
)

type command struct {
	name      string
	arguments []string
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
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Usage: %v <name>", cmd.name)
	}
	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		log.Fatalf("User does not exist: %v", err)
	}
	err = s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	log.Infof("User has been set to: %s\n", cmd.arguments[0])

	return nil
}

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			log.Fatalf("Unable to get user: %v", err)
		}
		return handler(s, cmd, user)
	}
}
