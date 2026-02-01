//go:build darwin

package volume

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"strings"
)

func FetchVolume() string {
	cmd := exec.CommandContext(
		context.Background(),
		"osascript",
		"-e", "output volume of (get volume settings)",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Printf("Unable to fetch volume: %v", err)
		return ""
	}

	v := strings.TrimSpace(out.String())
	return v, nil
}
