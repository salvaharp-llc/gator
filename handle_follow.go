package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/salvaharp-llc/gator/internal/database"
)

func handlerAddFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not find feed: %w", err)
	}

	followInfo, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed; %w", err)
	}

	fmt.Println("Feed follow created:")
	printFollow(followInfo.UserName, followInfo.FeedName)
	return nil
}

func handlerListFollows(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	followsInfo, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get followed feeds; %w", err)
	}

	if len(followsInfo) == 0 {
		fmt.Println("User does not follow any feeds")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, feedInfo := range followsInfo {
		fmt.Printf("- %s\n", feedInfo.FeedName)
	}
	return nil
}

func printFollow(username, feedname string) {
	fmt.Printf("- User:          %s\n", username)
	fmt.Printf("- Feed:          %s\n", feedname)
}
