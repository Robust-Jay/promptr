//go:build darwin

package clipboard

import (
	"os/exec"
)

func copyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(text))
	}()
	return cmd.Run()
}
