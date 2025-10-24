package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/salvaharp-llc/gator/internal/database"
)

func handlerAggregation(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Could not find a feed to fetch: %v", err)
		return
	}

	log.Println("Found a feed to fetch!")
	scrapeFeed(s, feed)
}

func scrapeFeed(s *state, feed database.Feed) {
	_, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Could not mark feed %s fetched: %v", feed.Name, err)
		return
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Could not fetch feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		pubDate, err := parseDate(item.PubDate)
		if err != nil {
			log.Printf("Could not parse pub date from %s: %v", item.Title, err)
			continue
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s fetched, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}

func parseDate(strDate string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		"02 Jan 2006 15:04:05 MST",
		"02 Jan 2006 15:04:05 -0700",
	}
	var timeDate time.Time
	var err error

	for _, layout := range layouts {
		timeDate, err = time.Parse(layout, strDate)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, err
	}
	return timeDate, nil
}
