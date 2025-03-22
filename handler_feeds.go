package main

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/shaneplunkett/gator/internal/database"
)

func handlerFeed(s *state, cmd command, user database.User) error {
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
	_, err = s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		log.Fatalf("Unable to follow feed: %v", err)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Usage: %v <url>", cmd.name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		log.Fatalf("Unable to get feed: %v", err)
	}
	feed_follow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		log.Fatalf("Unable to follow feed: %v", err)
	}
	fmt.Println("Feed Followed Successfully")
	fmt.Printf("Feed Name: %v\n", feed_follow.FeedName)
	fmt.Printf("User Name: %v\n", feed_follow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatalf("Failed to get Followed Feeds for User: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* Name:        %s\n", feed.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Usage: %v <url>", cmd.name)
	}
	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Url:    cmd.arguments[0],
		UserID: user.ID,
	})
	if err != nil {
		log.Fatalf("Failed to unfollow Feed: %v", err)
	}
	log.Info("Feed unfollowed successfully!\n")

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Usage: %v <request interval seconds>", cmd.name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		log.Fatalf("Unable to Parse Duration: %v", err)
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			log.Fatalf("Unable to Scrape Feeds: %v", err)
		}
	}
}
