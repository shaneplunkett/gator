package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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
	if cmd.arguments == nil {
		log.Fatalf("Username arguement required")
	}
	if len(cmd.arguments) > 1 {
		log.Fatalf("Login only accepts one arguement")
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

func handlerRegister(s *state, cmd command) error {
	if cmd.arguments == nil {
		log.Fatalf("Username arguement required")
	}
	if len(cmd.arguments) > 1 {
		log.Fatalf("Register only accepts one arguement")
	}
	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == nil {
		log.Fatalf("User already exists")
	} else if err != sql.ErrNoRows {
		log.Fatalf("Error checking for user: %v", err)
	}
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.arguments[0],
		},
	)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	err = s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	log.Infof("User Successfully Created: %s", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if cmd.arguments != nil {
		log.Fatalf("Reset does not accept any arguments")
	}
	err := s.db.DeleteUser(context.Background())
	if err != nil {
		log.Fatalf("Error deleting user table: %v", err)
	}
	log.Info("User table deleted successfully")

	return nil
}
