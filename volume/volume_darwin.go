//go:build darwin

package volume

import "C"
import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

func FetchVolume() (string, error) {
	cmd := exec.CommandContext(
		context.Background(),
		"osascript",
		"-e", "output volume of (get volume settings)",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	v := strings.TrimSpace(out.String())
	return v, nil
}
