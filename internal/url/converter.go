package url

import (
	"context"
	"fmt"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/path"
)

type Converter struct {
	config Config
}

type Config struct {
	BaseUrl        string
	ExistingUrls   map[string]bool
	SupportedLangs []language.Language
	SupportedExts  []string
	ValidSections  []path.Category
}

func NewConverter(config Config) *Converter {
	return &Converter{
		config: config,
	}
}

func (c *Converter) Convert(ctx context.Context, filePath string, fm *FrontMatter) (string, error) {
	if fm == nil {
		return "", fmt.Errorf("front matter is required for URL conversion of file: %s", filePath)
	}

	p, err := path.ParseWithValidation(filePath, c.config.SupportedLangs, c.config.SupportedExts, c.config.ValidSections)
	if err != nil {
		return "", err
	}

	if !fm.IsPublic() {
		return "", fmt.Errorf("file %s is not public, skipping URL generation", filePath)
	}

	url, err := generateUrl(c.config.BaseUrl, p, fm, c.config.ExistingUrls)
	if err != nil {
		return "", err
	}

	return url, nil
}
