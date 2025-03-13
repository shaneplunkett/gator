package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"

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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		log.Errorf("Error generating Request: %v", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	client := http.Client{}
	res, err := client.Do(req)
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

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return &feed, nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalf("Failed to Fetch Feed: %v", err)
	}
	fmt.Printf("Channel Title: %s\n", feed.Channel.Title)
	fmt.Printf("Channel Link: %s\n", feed.Channel.Link)
	fmt.Printf("Channel Description: %s\n", feed.Channel.Description)
	for i := range feed.Channel.Item {
		fmt.Printf("Item Title: %s\n", feed.Channel.Item[i].Title)
		fmt.Printf("Item Link: %s\n", feed.Channel.Item[i].Link)
		fmt.Printf("Item Description: %s\n", feed.Channel.Item[i].Description)
		fmt.Printf("Item PubDate: %s\n", feed.Channel.Item[i].PubDate)
	}

	return nil
}
