package git

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

// ParseGitLog parses the output of git log with --numstat format.
// The output format is:
// <hash>\x1F<author>\x1F<date>\x1F<message>
// <insertions>\t<deletions>\t<filepath>
//
// \n\n (empty line separates commits)
func ParseGitLog(output []byte) ([]*Event, error) {
	var events []*Event

	// Split by double newlines to get commit blocks
	blocks := bytes.Split(output, []byte("\n\n"))

	for _, block := range blocks {
		if len(block) == 0 {
			continue
		}

		lines := bytes.Split(block, []byte("\n"))
		if len(lines) < 2 {
			continue
		}

		// Parse metadata line
		metaParts := bytes.Split(lines[0], []byte("\x1F"))
		if len(metaParts) < 4 {
			continue
		}

		hash := string(metaParts[0])
		author := string(metaParts[1])
		date := string(metaParts[2])
		message := string(metaParts[3])

		// Parse numstat lines
		for i := 1; i < len(lines); i++ {
			line := lines[i]
			if len(line) == 0 {
				continue
			}

			parts := bytes.Split(line, []byte("\t"))
			if len(parts) < 3 {
				continue
			}

			insertionsStr := string(parts[0])
			deletionsStr := string(parts[1])
			filepath := string(parts[2])

			// Only process files under content/
			if !strings.Contains(filepath, "content/") {
				continue
			}

			// Parse insertions and deletions
			var insertions, deletions *int
			if insertionsStr != "-" {
				if val, err := strconv.Atoi(insertionsStr); err == nil {
					insertions = &val
				}
			}
			if deletionsStr != "-" {
				if val, err := strconv.Atoi(deletionsStr); err == nil {
					deletions = &val
				}
			}

			// Parse rename path
			oldPath, newPath := parseRenamePath(filepath)

			events = append(events, &Event{
				Hash:    hash,
				Author:  author,
				Date:    date,
				Message: message,
				File: FileInfo{
					Path:       newPath,
					Insertions: insertions,
					Deletions:  deletions,
					OldPath:    oldPath,
				},
			})
		}
	}

	return events, nil
}

// renameRegex matches git rename patterns like "old/path/{old => new}/suffix"
var renameRegex = regexp.MustCompile(`^(.*)\{(.*) => (.*)\}(.*)$`)

// parseRenamePath parses a git rename path and returns (oldPath, newPath).
// If the path is not a rename, returns ("", path).
func parseRenamePath(path string) (string, string) {
	matches := renameRegex.FindStringSubmatch(path)
	if len(matches) != 5 {
		return "", cleanPath(path)
	}

	prefix := matches[1]
	oldPart := strings.TrimSpace(matches[2])
	newPart := strings.TrimSpace(matches[3])
	suffix := matches[4]

	oldPath := cleanPath(prefix + oldPart + suffix)
	newPath := cleanPath(prefix + newPart + suffix)

	return oldPath, newPath
}

// cleanPath removes quotes and other unnecessary characters from a path.
func cleanPath(path string) string {
	path = strings.Trim(path, "\"")
	return path
}
