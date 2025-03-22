package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/shaneplunkett/gator/internal/database"

	"github.com/charmbracelet/log"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		log.Errorf("Error generating Request: %v", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("Error Making Request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	var feed RSSFeed
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	next, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Log(log.WarnLevel, "Failed to get Next Feed:", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: next.ID,
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Log(log.WarnLevel, "Failed Mark Feed Fetched:", err)
	}
	feed, err := fetchFeed(context.Background(), next.Url)
	if err != nil {
		log.Log(log.WarnLevel, "Failed to Fetch Feed:", err)
	}
	for _, item := range feed.Channel.Item {
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: parsePubDate(item.PubDate),
			FeedID:      next.ID,
		})
		if err != nil {
			log.Log(log.DebugLevel, "Unable to Add Post: %v", err)
		}
		fmt.Printf("* Title: %s\n", item.Title)
	}
	return nil
}

func parsePubDate(PubDate string) sql.NullTime {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822,
		time.RFC822Z,
		"2006-01-02T15:04:05Z",
	}
	for _, format := range formats {
		parsedTime, err := time.Parse(format, PubDate)
		if err == nil {
			return sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}
	}
	log.Logf(log.WarnLevel, "Unable to Convert Date: %v", PubDate)
	return sql.NullTime{
		Time:  time.Time{},
		Valid: false,
	}
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		log.Fatalf("Failed to Fetch Posts for User: %v", err)
	}
	if len(posts) == 0 {
		log.Fatal("No Posts Found for User")
	}
	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("* Publish:            %v\n", post.PublishedAt.Time)
	fmt.Printf("* Title:            %s\n", post.Title)
	fmt.Printf("* Url:            %s\n", post.Url)
	fmt.Printf("* Descrption:            %v\n", post.Description.String)
}
