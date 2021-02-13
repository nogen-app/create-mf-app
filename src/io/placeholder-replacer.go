package io

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//Placeholder is the struct used to create a placeholder with a key and the value it should be replaced with
type Placeholder struct {
	Key   string
	Value string
}

//comeup with better name
type PlaceholderReplacer struct {
	Placeholders []Placeholder
}

//ReplacePlaceholders walks through a given path, and replaces all placeholder keys with their respective value
//in all files and folders
func (pr *PlaceholderReplacer) ReplacePlaceholders(path string) error {
	if err := filepath.Walk(path, pr.replacePlaceholder); err != nil {
		return err
	}

	return nil
}

func (pr *PlaceholderReplacer) readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (pr *PlaceholderReplacer) writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func (pr *PlaceholderReplacer) replacePlaceholder(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		for _, placeholder := range pr.Placeholders {
			if info.Name() == placeholder.Key {
				newPath := strings.ReplaceAll(path, info.Name(), placeholder.Value)
				os.Rename(path, newPath)
			}
		}

	} else {
		lines, err := pr.readLines(path)
		if err != nil {
			return err
		}

		var newLines []string

		for _, line := range lines {
			newLine := line
			for _, placeholder := range pr.Placeholders {
				newLine = strings.ReplaceAll(newLine, placeholder.Key, placeholder.Value)
			}
			newLines = append(newLines, newLine)
		}

		if err := pr.writeLines(newLines, path); err != nil {
			return err
		}
	}

	return nil
}
