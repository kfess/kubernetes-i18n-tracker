package workflow

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kfess/kubernetes-i18n-tracker/internal/git"
)

// loadAllPaths reads all non-empty lines from a file.
func loadAllPaths(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var allPaths []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			allPaths = append(allPaths, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return allPaths, nil
}

// loadEvents reads git events from a JSONL file.
func loadEvents(path string) ([]*git.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var events []*git.Event
	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var event git.Event
		if err := json.Unmarshal(line, &event); err != nil {
			// ログに記録するが処理は継続
			continue
		}

		events = append(events, &event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return events, nil
}

// loadExistingPathsMap reads paths and returns them as a map for quick lookup.
func loadExistingPathsMap(path string) (map[string]bool, error) {
	paths, err := loadAllPaths(path)
	if err != nil {
		return nil, err
	}

	pathsMap := make(map[string]bool, len(paths))
	for _, p := range paths {
		pathsMap[p] = true
	}
	return pathsMap, nil
}
