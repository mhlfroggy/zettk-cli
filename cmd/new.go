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

func newAction(cmd *cli.Context) error {
	args := cmd.Args().Slice()
	// Find the Zettelkasten direcory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to find user home directory")
	}
	zettkDir := filepath.Join(homeDir, "zettelkasten")

	// Get the template
	tVal := cmd.String("template")
	template, err := os.ReadFile(filepath.Join(zettkDir, "templates", tVal+".md"))

	// Create the markdown file
	currTime := time.Now().Format("200601021504")
	fName := filepath.Join(zettkDir, "INBOX", fmt.Sprintf("%s-%s.md", currTime, filepath.Clean(args[0])))
	if tVal == "input" {
		fName = filepath.Join(zettkDir, "INPUT", fmt.Sprintf("%s-%s.md", currTime, filepath.Clean(args[0])))
	}
	// Expected value is "path/to/zettk/my-name.md"
	file, err := os.Create(fName)
	if err != nil {
		fmt.Println("Failed to create markdown file", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf(string(template), filepath.Clean(args[0]), time.Now().Format("2006-01-02")))
	if err != nil {
		fmt.Println("Failed to write file")
	}
	file.Close()

	// Add the new note to the daily note
	dNote := filepath.Join(zettkDir, "ARCHIVE", "daily-notes", fmt.Sprintf("%s.md", time.Now().Format("2006-01-02")))
	dFile, err := os.OpenFile(dNote, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to create daily note")
	}
	defer dFile.Close()
	link := "\n[[" + fmt.Sprintf("%s-%s.md", time.Now().Format("2006-01-02"), filepath.Clean(args[0])) + "]]"
	_, err = dFile.WriteString(link)
	if err != nil {
		fmt.Println(err)
	}
	dFile.Close()

	// Run neovim to open the note
	editor := exec.Command(os.Getenv("EDITOR"), fName)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	editor.Run()

	return nil
}
