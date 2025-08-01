package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ShowTree(root string, ignore []string) error {
	ignoreSet := make(map[string]struct{})
	for _, dir := range ignore {
		ignoreSet[dir] = struct{}{}
	}
	fmt.Println("📦 Project Structure")
	return printTree(root, "", ignoreSet)
}

func printTree(path string, prefix string, ignoreSet map[string]struct{}) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter out ignored
	var filtered []os.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		if _, ignored := ignoreSet[name]; !ignored {
			filtered = append(filtered, entry)
		}
	}
	entries = filtered

	// Sort alphabetically
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for i, entry := range entries {
		isLastEntry := i == len(entries)-1
		connector := "├───"
		if isLastEntry {
			connector = "└───"
		}

		entryName := entry.Name()
		fullPath := filepath.Join(path, entryName)
		display := fmt.Sprintf("%s%s %s", prefix, connector, entryName)

		if entry.IsDir() {
			display = fmt.Sprintf("%s 📁 %s", prefix+connector, entryName)
			fmt.Println(display)

			newPrefix := prefix
			if isLastEntry {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}

			printTree(fullPath, newPrefix, ignoreSet)
		} else {
			fmt.Println(display)
		}
	}
	return nil
}
