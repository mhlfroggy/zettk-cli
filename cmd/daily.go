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
	"time"
)

func dailyAction(cmd *cli.Context) error {
	// Find the Zettlekasten direcory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to find user home directory")
	}
	zettkDir := filepath.Join(homeDir, "zettelkasten")

	// Create the daily note
	dNote := filepath.Join(zettkDir, "ARCHIVE", "daily-notes", fmt.Sprintf("%s.md", time.Now().Format("2006-01-02")))
	dFile, err := os.OpenFile(dNote, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create daily note")
	}
	defer dFile.Close()
	dFile.Close()

	// Run neovim to open the note
	editor := exec.Command(os.Getenv("EDITOR"), dNote)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	editor.Run()

	return nil
}
