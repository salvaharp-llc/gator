package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not reset users: %w", err)
	}

	fmt.Println("users reset successfully")
	return nil
}
