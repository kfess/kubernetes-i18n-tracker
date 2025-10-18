// A Hugo's front matter parser

package url

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

type FrontMatterParser interface {
	Parse(filePath string) (*FrontMatter, error)
}

// Implements FrontMatterParser for YAML front matter
type YamlFrontMatterParser struct {
	rootDir string
}

func NewYAMLFrontMatterParser(rootDir string) *YamlFrontMatterParser {
	return &YamlFrontMatterParser{rootDir: rootDir}
}

func (p *YamlFrontMatterParser) Parse(path string) (*FrontMatter, error) {
	fullPath := filepath.Join(p.rootDir, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	content, err := extractFrontMatterString(file)
	if err != nil {
		return nil, err
	}

	if content == "" {
		return &FrontMatter{}, nil
	}

	fm := &FrontMatter{}
	if err := yaml.UnmarshalWithOptions([]byte(content), fm, yaml.AllowDuplicateMapKey()); err != nil {
		return nil, err
	}

	return fm, nil
}

func extractFrontMatterString(file *os.File) (string, error) {
	scanner := bufio.NewScanner(file)
	var frontMatter strings.Builder

	inFrontMatter := false
	delimiterCount := 0
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()

		if firstLine {
			// Some files has leading blank line
			// e.g., content/en/blog/_posts/2019-08-30-announcing-etcd-3.4.,md
			// line 1: <blank line>
			// line 2: ---
			// line 3: layout: blog
			line = strings.TrimLeft(line, "\n\r\t ")
			firstLine = false
		}

		if strings.TrimSpace(line) == "---" {
			delimiterCount++
			if delimiterCount == 1 {
				inFrontMatter = true
				continue
			} else if delimiterCount == 2 {
				// End of front matter
				break
			}
		}

		if inFrontMatter {
			frontMatter.WriteString(line)
			frontMatter.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return frontMatter.String(), nil
}
