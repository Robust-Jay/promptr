package prompt

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Prompt struct {
	ID       string   `yaml:"id"`
	Title    string   `yaml:"title"`
	Category []string `yaml:"category"`
	Content  string   `yaml:"content"`
}

func Load(path string) (*Prompt, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var p Prompt
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	p.ID = strings.TrimSpace(p.ID)
	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)
	for i, c := range p.Category {
		p.Category[i] = strings.TrimSpace(c)
	}
	return &p, nil
}

func (p *Prompt) Save(path string) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal prompt: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
