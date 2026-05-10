package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func BaseDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".promptr"), nil
}

func EnsureDirs() (string, error) {
	base, err := BaseDir()
	if err != nil {
		return "", err
	}
	for _, sub := range []string{"builtin", "user"} {
		if err := os.MkdirAll(filepath.Join(base, sub), 0755); err != nil {
			return "", fmt.Errorf("create %s: %w", sub, err)
		}
	}
	return base, nil
}

func isYAML(name string) bool {
	return strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")
}

func walkPrompts(root string) ([]*Prompt, error) {
	var prompts []*Prompt
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !isYAML(info.Name()) {
			return nil
		}
		p, err := Load(path)
		if err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}
		prompts = append(prompts, p)
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return prompts, nil
}

func LoadAll() ([]*Prompt, error) {
	base, err := BaseDir()
	if err != nil {
		return nil, err
	}

	builtinPrompts, _ := walkPrompts(filepath.Join(base, "builtin"))
	userPrompts, err := walkPrompts(filepath.Join(base, "user"))
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var result []*Prompt
	for _, p := range userPrompts {
		seen[p.ID] = true
		result = append(result, p)
	}
	for _, p := range builtinPrompts {
		if seen[p.ID] {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

func FindByID(id string) (*Prompt, string, error) {
	base, err := BaseDir()
	if err != nil {
		return nil, "", err
	}

	userPath := filepath.Join(base, "user", id+".yaml")
	if _, err := os.Stat(userPath); err == nil {
		p, err := Load(userPath)
		return p, userPath, err
	}

	builtinPath := filepath.Join(base, "builtin", id+".yaml")
	entries, err := os.ReadDir(filepath.Join(base, "builtin"))
	if err == nil {
		for _, entry := range entries {
			subDir := filepath.Join(base, "builtin", entry.Name())
			if info, err := os.Stat(subDir); err == nil && info.IsDir() {
				candidate := filepath.Join(subDir, id+".yaml")
				if _, err := os.Stat(candidate); err == nil {
					p, err := Load(candidate)
					return p, candidate, err
				}
			}
		}
	}

	// fallback: try flat builtin
	if _, err := os.Stat(builtinPath); err == nil {
		p, err := Load(builtinPath)
		return p, builtinPath, err
	}

	return nil, "", fmt.Errorf("prompt not found: %s", id)
}

func Search(query string) ([]*Prompt, error) {
	all, err := LoadAll()
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var results []*Prompt
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.ID), q) ||
			strings.Contains(strings.ToLower(p.Title), q) ||
			containsCategory(p.Category, q) ||
			strings.Contains(strings.ToLower(p.Content), q) {
			results = append(results, p)
		}
	}
	return results, nil
}

func List(categoryFilter string) ([]*Prompt, error) {
	all, err := LoadAll()
	if err != nil {
		return nil, err
	}
	if categoryFilter == "" {
		return all, nil
	}
	filters := strings.Split(categoryFilter, "/")
	var results []*Prompt
	for _, p := range all {
		if hasAllCategories(p.Category, filters) {
			results = append(results, p)
		}
	}
	return results, nil
}

func AllCategories() ([]string, error) {
	all, err := LoadAll()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	var cats []string
	for _, p := range all {
		for _, c := range p.Category {
			c = strings.TrimSpace(c)
			if c != "" && !seen[c] {
				seen[c] = true
				cats = append(cats, c)
			}
		}
	}
	return cats, nil
}

func UserFilePath(id string) (string, error) {
	base, err := BaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "user", id+".yaml"), nil
}

func IsBuiltin(id string) bool {
	_, path, err := FindByID(id)
	if err != nil {
		return false
	}
	base, _ := BaseDir()
	return !strings.HasPrefix(path, filepath.Join(base, "user"))
}

func containsCategory(categories []string, query string) bool {
	for _, c := range categories {
		if strings.Contains(strings.ToLower(c), query) {
			return true
		}
	}
	return false
}

func hasAllCategories(promptCats, filters []string) bool {
	for _, f := range filters {
		found := false
		for _, c := range promptCats {
			if strings.EqualFold(strings.TrimSpace(c), strings.TrimSpace(f)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
