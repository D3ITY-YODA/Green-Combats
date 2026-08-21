package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"green-compass-backend/internal/config"
	"green-compass-backend/internal/migrate"
	"green-compass-backend/migrations"
	"green-compass-backend/pkg/database"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

	flags := flag.NewFlagSet("migrate", flag.ContinueOnError)
	dsnFlag := flags.String("dsn", "", "database url (overrides GC_DATABASE_URL and config file)")
	steps := flags.Int("steps", 0, "number of migrations to apply or revert (0 = all)")
	flags.Usage = printUsage
	if err := flags.Parse(os.Args[1:]); err != nil {
		return 2
	}
	if flags.NArg() != 1 {
		printUsage()
		return 2
	}

	url, err := resolveDSN(*dsnFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		return 2
	}

	pool, err := database.Connect(ctx, database.Options{URL: url})
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		return 1
	}
	defer pool.Close()

	switch flags.Arg(0) {
	case "up":
		applied, err := migrate.Up(ctx, pool, migrations.FS, *steps)
		report("applied", applied, err)
		return exitCode(err)
	case "down":
		reverted, err := migrate.Down(ctx, pool, migrations.FS, *steps)
		report("reverted", reverted, err)
		return exitCode(err)
	case "status":
		st, err := migrate.Current(ctx, pool, migrations.FS)
		if err != nil {
			fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
			return 1
		}
		fmt.Printf("current version: %d\n", st.Current)
		if len(st.Pending) == 0 {
			fmt.Println("database is up to date")
		} else {
			fmt.Printf("pending versions: %v\n", st.Pending)
		}
		return 0
	default:
		fmt.Fprintf(os.Stderr, "migrate: unknown command %q\n", flags.Arg(0))
		printUsage()
		return 2
	}
}

func resolveDSN(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if env := os.Getenv("GC_DATABASE_URL"); env != "" {
		return env, nil
	}
	cfg, err := config.Load(config.LoadOptions{})
	if err != nil {
		return "", fmt.Errorf("load config for database url: %w", err)
	}
	if cfg.Database.URL == "" {
		return "", fmt.Errorf("no database url: pass -dsn, set GC_DATABASE_URL, or set database.url in config")
	}
	return cfg.Database.URL, nil
}

func report(verb string, versions []int64, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
	}
	if len(versions) > 0 {
		fmt.Printf("%s migrations: %v\n", verb, versions)
	} else if err == nil {
		fmt.Println("nothing to do")
	}
}

func exitCode(err error) int {
	if err != nil {
		return 1
	}
	return 0
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `usage: migrate [-dsn url] [-steps n] <command>

commands:
  up      apply pending migrations
  down    revert applied migrations (newest first)
  status  show current and pending migration versions`)
}
