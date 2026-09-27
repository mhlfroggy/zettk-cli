/*
Copyright © 2026 Matthew Labrecque <mlabrecque2002@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func findAction(cmd *cli.Context) error {
	args := cmd.Args().Slice()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to find user home directory")
	}

	searchQuerry := filepath.Clean(args[0])

	directories := []string{"INBOX", "ARCHIVE", "INPUT"}
	inbox := cmd.Bool("inbox")
	input := cmd.Bool("input")
	archive := cmd.Bool("archive")

	// Note that the slice indexes are being updated dynamically
	if inbox {
		directories = slices.Delete(directories, 1, 3)
	} else if input {
		directories = slices.Delete(directories, 0, 1)
		directories = slices.Delete(directories, 1, 2)
	} else if archive {
		directories = slices.Delete(directories, 0, 1)
		directories = slices.Delete(directories, 0, 1)
	}

	// Walk through all the directories
	var foundFiles []string = make([]string, 0, 10)
	targetRegexp := regexp.MustCompile(".*" + searchQuerry + ".*")
	for _, dir := range directories {
		files, err := os.ReadDir(filepath.Join(homeDir, "zettelkasten", dir))
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range files {
			if !file.IsDir() && targetRegexp.MatchString(file.Name()) {
				// Note we have to add the dir for file path construction
				// later on
				foundFiles = append(foundFiles, dir+"/"+file.Name())
			}
		}
	}

	// Determine if additional querry is needed
	var foundFile string = ""
	if len(foundFiles) == 0 {
		fmt.Println("Found 0 notes matching search querry:", searchQuerry)
		return nil
	} else if len(foundFiles) > 1 {
		var ans int
		fmt.Println(fmt.Sprintf("Multiple files matching search term \"%s\" found.", searchQuerry))
		for i := 0; i < len(foundFiles); i++ {
			fmt.Println(fmt.Sprintf("  %d) %s", i, foundFiles[i]))
		}
		fmt.Println()
		fmt.Print(fmt.Sprintf("Which file do you want? [0-%d] ", len(foundFiles)-1))
		fmt.Scan(&ans)
		for ans < 0 || ans >= len(foundFiles) {
			fmt.Println("Please select from within the valid options.")
			fmt.Println()
			fmt.Print(fmt.Sprintf("Which file do you want? [0-%d] ", len(foundFiles)-1))
			fmt.Scanf("%d", &ans)
		}
		foundFile = foundFiles[ans]
	} else {
		foundFile = foundFiles[0]
	}

	// Reconstruct the file path
	fullPath, err := os.Stat(filepath.Join(homeDir, "zettelkasten", foundFile))
	if err != nil {
		fmt.Println(err)
	}

	// Get all info about the file (note we have to do some additional
	// processing here)
	split := strings.SplitN(foundFile, "/", 2)
	split2 := strings.SplitN(split[1], "-", 2)
	fileName := split2[1]
	zettkID := split2[0]
	location := split[0]
	cTime := cTime(split2[0])
	mTime := fullPath.ModTime().Format("2006-01-02 15:04")

	// Display information about the file
	fmt.Println(fmt.Sprintf("Info about %s\n ZettkID: %s\n Location: %s\n Creation Time: %s\n Last Modified: %s", fileName, zettkID, location, cTime, mTime))

	return nil
}

// cTime formats a timestamp identifier as a human-readable creation time.
func cTime(ctime string) string {
	s := ctime
	return s[:4] + "-" + s[4:6] + "-" + s[6:8] + " " + s[8:10] + ":" + s[10:12]
}
