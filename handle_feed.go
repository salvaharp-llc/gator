package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/salvaharp-llc/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	userName := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return err
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	fmt.Printf("Id: %v\nName: %s\nURL: %s\nUserID: %v\n",
		feed.ID,
		feed.Name,
		feed.Url,
		feed.UserID,
	)
	return nil
}
