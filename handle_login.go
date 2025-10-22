package main

import (
	"context"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("unknown user: %w", err)
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Printf("User set to %s\n", userName)
	return nil
}
