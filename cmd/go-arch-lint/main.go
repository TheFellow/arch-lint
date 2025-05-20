package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/TheFellow/go-arch-lint/pkg/config"
	"github.com/TheFellow/go-arch-lint/pkg/linter"
)

func main() {
	app := &cli.Command{
		Name:  "arch-lint",
		Usage: "Enforce forbidden import rules via globs",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "config.yaml", Usage: "Path to config file"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			cfg, err := config.Load(c.String("config"))
			if err != nil {
				return err
			}

			violations, err := linter.Run(cfg)
			if err != nil {
				return err
			}
			if len(violations) > 0 {
				for _, v := range violations {
					log.Println(v)
				}
				os.Exit(1)
			}
			log.Println("✔ arch-lint: no forbidden imports found.")
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
