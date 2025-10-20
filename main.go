package main

import (
	"log"
	"os"

	"github.com/salvaharp-llc/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := state{
		cfg: &cfg,
	}
	cmds := commands{
		registry: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
}
