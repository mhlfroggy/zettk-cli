/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
)

func initAction(cmd *cli.Context) error {

	// Grab the home directory and build the directory path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to find user home directory")
	}
	baseDir := filepath.Join(homeDir, "zettelkasten")

	// Check if a Zettelkasten already exists at the given location
	// If it does exist, prompt the user
	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		var ans string = "N"
		fmt.Println("=== WARNING ===")
		fmt.Println("Directory exists at", baseDir)
		fmt.Print("Overwrite existing directory? (Y/N) ")
		fmt.Scan(&ans)
		if err != nil {
			fmt.Println(err)
		}
		if strings.ToUpper(ans) != "Y" {
			fmt.Println("Aborting")
			return nil
		} else {
			os.RemoveAll(baseDir)
		}
	}

	zettDirs := []string{"INBOX", "ARCHIVE", "INPUT", "ARCHIVE/daily-notes", "templates"}

	// Create the folder structure
	for _, zettDir := range zettDirs {
		err := os.MkdirAll(filepath.Join(baseDir, zettDir), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Create note.md default template and write to templates
	noteContents := `---
type: note
title: %s
created: %s
tags:
    - 
---`

	defaultTemp := filepath.Join(baseDir, "templates", "note.md")
	file, err := os.Create(defaultTemp)
	if err != nil {
		fmt.Println("Failed to create markdown file", err)
	}
	defer file.Close()
	_, err = file.WriteString(noteContents)
	if err != nil {
		fmt.Println("Failed to write file")
	}
	file.Close()

	return nil
}
