# gator

gator is a small command-line blog aggregator that fetches posts from RSS/Atom feeds, stores them in Postgres, and provides a simple CLI to manage feeds and users.

## Prerequisites

- Go (1.20+ recommended) installed and on your PATH. See https://golang.org/doc/install
- PostgreSQL server available and a database created for gator
- See the "Creating the database for gator" section below for migration instructions (goose) and example commands.

## Install

Install the gator CLI with `go install`:

```bash
go install github.com/salvaharp-llc/gator@latest
```

## Configuration

gator expects a JSON config file in your home directory named `.gatorconfig.json` with the following structure:

```json
{"db_url":"postgres://username@localhost:5432/gator?sslmode=disable"}
```

Place this file at `~/.gatorconfig.json`. Replace `username`, `localhost`, `5432`, and `gator` with your database connection details.

Notes:
- Use `sslmode=disable` for local development when Postgres does not use TLS.
- The program reads this file at runtime to obtain the `db_url`.

## Usage

Run the CLI like:

```bash
gator <command> [args...]
```

Available commands
- `login` — Log in a user (used by commands that require authentication).
- `register` — Register a new user account and log in to said account.
- `reset` — Reset the database.
- `users` — List users in the system.
- `agg` — Trigger aggregation/fetch for feeds (pulls the latest posts).
- `addfeed` — Add a feed to the system (requires login).
- `feeds` — List all known feeds.
- `follow` — Follow a feed for the current user (requires login).
- `following` — List feeds the current user is following (requires login).
- `unfollow` — Unfollow a feed (requires login).
- `browse` — Browse posts (requires login).

Examples

Create a new user:

```bash
gator register <name>
```

Add a feed:

```bash
gator addfeed <url>
```

Start the aggregator:

```bash
gator agg 30s
```

View the posts:

```bash
gator browse [limit]
```

## Creating the database for gator

The repository contains SQL migrations in `sql/schema`. We use goose to apply these migrations.

Important: a binary installed with `go install` does not include the repository files (including the `sql/schema` migration files). To run migrations you should clone the repo and run goose from the project root, or use another migration runner that has access to the SQL files.

Install the goose CLI:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Create the database (example):

```sql
CREATE DATABASE gator;
```

Run migrations (from the repository root):

```bash
export DB_URL="postgres://postgres:your_password@localhost:5432/gator?sslmode=disable"
goose -dir sql/schema postgres "$DB_URL" up
```

Check migration status:

```bash
goose -dir sql/schema postgres "$DB_URL" status
```
