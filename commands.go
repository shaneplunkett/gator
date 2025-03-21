package main

import (
	"context"
	"database/sql"
	"fmt"
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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Usage: %v <name>", cmd.name)
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
	err := s.db.DeleteUser(context.Background())
	if err != nil {
		log.Fatalf("Error deleting user table: %v", err)
	}
	log.Info("User table cleared successfully")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("Error getting users: %v", err)
	}
	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalf("Failed to Fetch Feed: %v", err)
	}
	fmt.Printf("Feed: %+v\n", feed)

	return nil
}

func handlerFeed(s *state, cmd command) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("Usage: %v <name> <url>", cmd.name)
	}
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		log.Fatalf("Unable to get user: %v", err)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	})
	if err != nil {
		log.Fatalf("Unable to generate feed: %v", err)
	}
	fmt.Printf("Feed: %+v\n", feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("Unable to get feeds: %v", err)
	}
	if len(feeds) == 0 {
		log.Fatal("No Feeds Found")
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			log.Fatalf("Failed to get user by ID: %v", err)
		}
		printFeed(feed, user)
		fmt.Println("======================================")
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
