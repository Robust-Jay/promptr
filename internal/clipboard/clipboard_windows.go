package clipboard

import (
	"os"
	"os/exec"
	"path/filepath"
)

func copyToClipboard(text string) error {
	dir := os.TempDir()
	tmpFile := filepath.Join(dir, "promptr_clip.txt")

	if err := os.WriteFile(tmpFile, []byte(text), 0644); err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	psCmd := "Add-Type -AssemblyName System.Windows.Forms;" +
		"$c = Get-Content -LiteralPath '" + tmpFile + `' -Encoding UTF8 -Raw;` +
		"[System.Windows.Forms.Clipboard]::SetText($c)"

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	return cmd.Run()
}
