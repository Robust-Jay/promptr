//go:build linux

package clipboard

import (
	"os/exec"
)

func copyToClipboard(text string) error {
	for _, tool := range []string{"xclip", "xsel"} {
		if _, err := exec.LookPath(tool); err == nil {
			var cmd *exec.Cmd
			switch tool {
			case "xclip":
				cmd = exec.Command("xclip", "-selection", "clipboard")
			case "xsel":
				cmd = exec.Command("xsel", "--clipboard", "--input")
			}
			stdin, err := cmd.StdinPipe()
			if err != nil {
				continue
			}
			go func() {
				defer stdin.Close()
				stdin.Write([]byte(text))
			}()
			return cmd.Run()
		}
	}
	return exec.ErrNotFound
}
