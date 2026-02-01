//go:build linux

package active_window

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

var activeWindowFetchers = []func() (string, error){
	fetchActiveWindowHyprCtl,
}

func FetchActiveWindowName() string {
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
		return "", fmt.Errorf("class not found")
	}

	return formatActiveString(strings.TrimSpace(match[1])), nil
}

func formatActiveString(s string) string {
	// Replace common separators with spaces
	replacer := strings.NewReplacer(
		"-", " ",
		"_", " ",
	)
	s = replacer.Replace(s)

	// Split on whitespace
	words := strings.Fields(s)

	// Capitalize each word
	for i, w := range words {
		runes := []rune(w)
		if len(runes) == 0 {
			continue
		}

		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}

		words[i] = string(runes)
	}

	return strings.Join(words, " ")
}
