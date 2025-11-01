package url

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-yaml"
)

type FrontMatterParser interface {
	Parse(content string) (*FrontMatter, error)
}

type YamlFrontMatterParser struct{}

func NewYAMLFrontMatterParser() *YamlFrontMatterParser {
	return &YamlFrontMatterParser{}
}

func (p *YamlFrontMatterParser) Parse(content string) (*FrontMatter, error) {
	fmString, err := extractFrontMatterString(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	if fmString == "" {
		return &FrontMatter{}, nil
	}

	fm := &FrontMatter{}
	if err := yaml.UnmarshalWithOptions([]byte(fmString), fm, yaml.AllowDuplicateMapKey()); err != nil {
		return nil, fmt.Errorf("failed to unmarshal front matter: %w", err)
	}

	return fm, nil
}

// extractFrontMatterString extracts front matter string from reader
func extractFrontMatterString(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	var frontMatter strings.Builder

	inFrontMatter := false
	delimiterCount := 0
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()

		if firstLine {
			// Some files have leading blank lines
			// e.g., content/en/blog/_posts/2019-08-30-announcing-etcd-3.4.md
			// line 1: <blank line>
			// line 2: ---
			// line 3: layout: blog
			line = strings.TrimPrefix(line, "\ufeff")
			line = strings.TrimLeft(line, "\n\r\t ")
			firstLine = false
		}

		// Some files has invalid front matter delimiters like "------"
		// https://github.com/kubernetes/website/blob/main/content/en/blog/_posts/2024-04-24-validating-admission-policy-ga/index.md
		if strings.TrimSpace(line) == "---" || strings.TrimSpace(line) == "------" {
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

	// If there are less than 2 delimiters, return empty string
	if delimiterCount < 2 {
		return "", nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return frontMatter.String(), nil
}
