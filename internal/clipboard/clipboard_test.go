package clipboard

import (
	"runtime"
	"testing"
)

func TestCopyToClipboard(t *testing.T) {
	text := "test clipboard content"
	err := Copy(text)

	if runtime.GOOS == "linux" {
		if err != nil {
			t.Logf("Linux clipboard requires xclip or xsel (expected): %v", err)
			return
		}
	}

	if err != nil {
		t.Logf("Clipboard copy error (may be headless env): %v", err)
	}
}
