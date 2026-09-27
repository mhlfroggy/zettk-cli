package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// NewApp constructs the command line application.
func NewApp() *cli.App {
	return &cli.App{
		Name:  "zettk-cli",
		Usage: "A KISS CLI interface for a Zettelkasten note taking system",
		Description: `Zettk-CLI is a CLI for a Zettelkasten note taking system
utilizing the markdown format. Zettk is designed to be easily
customizable and modifiable, serving as a core for your
note taking system.`,
		Commands: []*cli.Command{
			{Name: "init", Usage: "Initialize a Zettelkasten", Description: "Initialize a new Zettelkasten with subdirectories. Existing directories require confirmation before replacement.", Before: exactArgs(0), Action: initAction},
			{Name: "new", Usage: "Create a new note and update the Zettelkasten", ArgsUsage: "<title>", Description: "Create a markdown note, link it from today's daily note, and open it in $EDITOR.", Flags: []cli.Flag{&cli.StringFlag{Name: "template", Aliases: []string{"t"}, Value: "note", Usage: "Specify custom note template"}}, Before: exactArgs(1), Action: newAction},
			{Name: "open", Usage: "Open a file in the default editor", ArgsUsage: "<search>", Description: "Search for a matching note and open it in your editor.", Before: exactArgs(1), Action: openAction},
			{Name: "find", Usage: "Search through the Zettelkasten", ArgsUsage: "<search>", Description: "Find a note and display its ID, creation time, location, and modification time.", Flags: []cli.Flag{
				&cli.BoolFlag{Name: "input", Usage: "Search only the input folder"},
				&cli.BoolFlag{Name: "archive", Usage: "Search only the archive folder"},
				&cli.BoolFlag{Name: "inbox", Usage: "Search only the inbox folder"},
			}, Before: func(c *cli.Context) error {
				if err := exactArgs(1)(c); err != nil {
					return err
				}
				selected := 0
				for _, name := range []string{"input", "archive", "inbox"} {
					if c.Bool(name) {
						selected++
					}
				}
				if selected > 1 {
					return errors.New("--input, --archive, and --inbox are mutually exclusive")
				}
				return nil
			}, Action: findAction},
			{Name: "sp", Usage: "Open the scratchpad note for quick notes", Before: exactArgs(0), Action: scratchpadAction},
			{Name: "daily", Usage: "Open the daily note", Before: exactArgs(0), Action: dailyAction},
		},
	}
}

func exactArgs(want int) cli.BeforeFunc {
	return func(c *cli.Context) error {
		if c.NArg() != want {
			return fmt.Errorf("expected %d argument(s), got %d", want, c.NArg())
		}
		return nil
	}
}

// Execute runs the CLI using the process arguments.
func Execute() error { return NewApp().Run(os.Args) }

// Main prints application errors and returns a nonzero exit code.
func Main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
