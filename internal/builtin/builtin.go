package builtin

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"promptr/internal/prompt"
)

//go:embed code/*.yaml writing/*.yaml general/*.yaml
var embedded embed.FS

func Extract() error {
	base, err := prompt.BaseDir()
	if err != nil {
		return err
	}
	builtinDir := filepath.Join(base, "builtin")

	if info, err := os.Stat(builtinDir); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(builtinDir)
		if len(entries) > 0 {
			return nil
		}
	}

	if err := os.MkdirAll(builtinDir, 0755); err != nil {
		return fmt.Errorf("create builtin dir: %w", err)
	}

	err = fs.WalkDir(embedded, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := embedded.ReadFile(path)
		if err != nil {
			return err
		}
		targetDir := filepath.Join(builtinDir, filepath.Dir(path))
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}
		target := filepath.Join(builtinDir, path)
		if err := os.WriteFile(target, data, 0644); err != nil {
			return fmt.Errorf("write %s: %w", target, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("extract builtin prompts: %w", err)
	}
	return nil
}

func Init() error {
	if _, err := prompt.EnsureDirs(); err != nil {
		return err
	}
	return Extract()
}
