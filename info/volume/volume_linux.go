//go:build linux

package volume

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var volumeFetchers = []func() (string, error){
	getWirePlumberVolume,
	getPulseAudioVolume,
}

func fetchVolume() string {
	var volume string

	for _, volumeFetcher := range volumeFetchers {
		output, err := volumeFetcher()
		if err == nil {
			volume = strings.TrimSpace(output)
			return volume
		}
		log.Printf("Unable to fetch volume: %v", err)
	}

	return ""
}

func getWirePlumberVolume() (string, error) {
	cmd := exec.CommandContext(
		context.Background(),
		"wpctl",
		"get-volume", "@DEFAULT_SINK@",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	fields := strings.Fields(out.String())
	if len(fields) < 2 {
		return "", errors.New("unexpected wpctl output")
	}

	volume, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(volume * 100)), nil
}

func getPulseAudioVolume() (string, error) {
	cmd := exec.CommandContext(
		context.Background(),
		"pactl",
		"get-sink-volume", "@DEFAULT_SINK@",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(\d+)%`)
	matches := re.FindAllStringSubmatch(out.String(), -1)

	if len(matches) == 0 {
		return "", errors.New("no volume percentages found")
	}

	sum := 0
	for _, m := range matches {
		v, _ := strconv.Atoi(m[1])
		sum += v
	}

	return strconv.Itoa(sum / len(matches)), nil
}
