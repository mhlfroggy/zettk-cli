/*
Copyright © 2026 Matthew Labrecque <mlabrecque2002@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path/filepath"
)

func scratchpadAction(cmd *cli.Context) error {
	// Find the Zettelkasten direcory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to find user home directory")
	}
	zettkDir := filepath.Join(homeDir, "zettelkasten")

	// Add the new note to the daily note
	spPath := filepath.Join(zettkDir, "PRIVATE", "scratchpad.md")
	spFile, err := os.OpenFile(spPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to find or create the scratchpad")
	}
	spFile.Close()

	// Run neovim to open the note
	editor := exec.Command(os.Getenv("EDITOR"), spPath)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	editor.Run()

	return nil
}
