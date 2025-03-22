package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"

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
		log.Fatalf("Failed to get Next Feed: %v", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: next.ID,
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Fatalf("Failed Mark Feed Fetched: %v", err)
	}
	feed, err := fetchFeed(context.Background(), next.Url)
	if err != nil {
		log.Fatalf("Failed to Fetch Feed: %v", err)
	}
	for _, item := range feed.Channel.Item {
		fmt.Printf("* Title: %s\n", item.Title)
	}

	return nil
}
