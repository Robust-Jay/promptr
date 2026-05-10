package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPromptSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")

	p := &Prompt{
		ID:       "test-id",
		Title:    "Test Title",
		Category: []string{"cat1", "cat2"},
		Content:  "Hello {{name}}",
	}
	if err := p.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ID != "test-id" {
		t.Errorf("ID = %q, want %q", loaded.ID, "test-id")
	}
	if loaded.Title != "Test Title" {
		t.Errorf("Title = %q, want %q", loaded.Title, "Test Title")
	}
	if len(loaded.Category) != 2 {
		t.Errorf("Category len = %d, want 2", len(loaded.Category))
	}
	if !strings.Contains(loaded.Content, "{{name}}") {
		t.Errorf("Content doesn't contain placeholder")
	}
}

func TestExtractVars(t *testing.T) {
	tests := []struct {
		content string
		want    []string
	}{
		{"Hello {{name}}", []string{"name"}},
		{"{{a}} and {{b}}", []string{"a", "b"}},
		{"No variables here", []string{}},
		{"{{duplicate}} {{duplicate}}", []string{"duplicate"}},
		{"{{multi_word_var}}", []string{"multi_word_var"}},
	}

	for _, tt := range tests {
		got := ExtractVars(tt.content)
		if len(got) != len(tt.want) {
			t.Errorf("ExtractVars(%q) = %v, want %v", tt.content, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("ExtractVars(%q)[%d] = %q, want %q", tt.content, i, got[i], tt.want[i])
			}
		}
	}
}

func TestSearch(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", dir)
	defer os.Setenv("USERPROFILE", oldHome)

	base := filepath.Join(dir, ".promptr")
	os.MkdirAll(filepath.Join(base, "builtin", "code"), 0755)
	os.MkdirAll(filepath.Join(base, "user"), 0755)

	p1 := &Prompt{ID: "search-test", Title: "SearchTitle", Category: []string{"search"}, Content: "content with keyword"}
	p1.Save(filepath.Join(base, "builtin", "code", "search-test.yaml"))

	p2 := &Prompt{ID: "other", Title: "Other", Category: []string{"other"}, Content: "other content"}
	p2.Save(filepath.Join(base, "builtin", "code", "other.yaml"))

	results, err := Search("keyword")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Search results = %d, want 1", len(results))
	}
	if results[0].ID != "search-test" {
		t.Errorf("Found ID = %q, want %q", results[0].ID, "search-test")
	}

	results, err = Search("SearchTitle")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Search by title results = %d, want 1", len(results))
	}

	results, err = Search("nonexistent")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Search for nonexistent returned %d results, want 0", len(results))
	}
}

func TestUserOverridesBuiltin(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", dir)
	defer os.Setenv("USERPROFILE", oldHome)

	base := filepath.Join(dir, ".promptr")
	os.MkdirAll(filepath.Join(base, "builtin", "code"), 0755)
	os.MkdirAll(filepath.Join(base, "user"), 0755)

	builtin := &Prompt{ID: "override", Title: "Builtin", Category: []string{"code"}, Content: "builtin content"}
	builtin.Save(filepath.Join(base, "builtin", "code", "override.yaml"))

	user := &Prompt{ID: "override", Title: "User", Category: []string{"custom"}, Content: "user content"}
	user.Save(filepath.Join(base, "user", "override.yaml"))

	all, err := LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	found := false
	for _, p := range all {
		if p.ID == "override" {
			found = true
			if p.Title != "User" {
				t.Errorf("Override title = %q, want %q", p.Title, "User")
			}
			break
		}
	}
	if !found {
		t.Error("Override prompt not found")
	}
}

func TestListByCategory(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", dir)
	defer os.Setenv("USERPROFILE", oldHome)

	base := filepath.Join(dir, ".promptr")
	os.MkdirAll(filepath.Join(base, "builtin", "code"), 0755)
	os.MkdirAll(filepath.Join(base, "user"), 0755)

	p1 := &Prompt{ID: "a", Title: "A", Category: []string{"code", "testing"}, Content: "a"}
	p1.Save(filepath.Join(base, "builtin", "code", "a.yaml"))

	p2 := &Prompt{ID: "b", Title: "B", Category: []string{"code", "review"}, Content: "b"}
	p2.Save(filepath.Join(base, "builtin", "code", "b.yaml"))

	p3 := &Prompt{ID: "c", Title: "C", Category: []string{"writing"}, Content: "c"}
	p3.Save(filepath.Join(base, "builtin", "code", "c.yaml"))

	results, err := List("code/testing")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 1 || results[0].ID != "a" {
		t.Errorf("code/testing list = %v, want [a]", results)
	}

	results, err = List("code")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("code list = %d, want 2", len(results))
	}
}

func TestAllCategories(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", dir)
	defer os.Setenv("USERPROFILE", oldHome)

	base := filepath.Join(dir, ".promptr")
	os.MkdirAll(filepath.Join(base, "builtin", "code"), 0755)

	p1 := &Prompt{ID: "a", Title: "A", Category: []string{"code", "testing"}, Content: "a"}
	p1.Save(filepath.Join(base, "builtin", "code", "a.yaml"))

	cats, err := AllCategories()
	if err != nil {
		t.Fatalf("AllCategories failed: %v", err)
	}
	if len(cats) < 2 {
		t.Errorf("Categories = %v, want at least 2", cats)
	}
}
