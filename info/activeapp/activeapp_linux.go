//go:build linux

package activeapp

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

var activeWindowFetchers = []func() (string, error){
	fetchActiveWindowHyprCtl,
}

func fetchActiveAppName() string {
	var activeWindow string

	for _, activeWindowFetcher := range activeWindowFetchers {
		output, err := activeWindowFetcher()
		if err == nil {
			activeWindow = strings.TrimSpace(output)
			return activeWindow
		}
		log.Printf("Unable to fetch active window: %v", err)
	}

	return ""
}

func fetchActiveWindowHyprCtl() (string, error) {
	cmd := exec.CommandContext(
		context.Background(),
		"hyprctl",
		"activewindow",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Target out string is class: jetbrains-goland
	regex := regexp.MustCompile(`(?m)^\s*class:\s*(.+)$`)
	match := regex.FindStringSubmatch(out.String())
	if len(match) < 2 {
		return "", log.Println("class not found")
	}

	return formatAppString(strings.TrimSpace(match[1])), nil
}
