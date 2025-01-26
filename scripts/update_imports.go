package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "c:/Programmering/Go/horse_tracking_go"
	oldImport := "github.com/polyfant/hulta_pregnancy_app"
	newImport := "github.com/polyfant/horse_tracking_go"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if strings.Contains(path, ".git") {
			return nil
		}

		// Only process .go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Read file
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace imports
		newContent := strings.ReplaceAll(string(content), oldImport, newImport)

		// Write back if changed
		if newContent != string(content) {
			fmt.Printf("Updating imports in: %s\n", path)
			err = ioutil.WriteFile(path, []byte(newContent), 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
