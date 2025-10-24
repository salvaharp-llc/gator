package main

import (
	"context"
	"fmt"
	"log"
	"time"
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

	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
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
		fmt.Printf("Found post %s\n", item.Title)
	}
	log.Printf("Feed %s fetched, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}
