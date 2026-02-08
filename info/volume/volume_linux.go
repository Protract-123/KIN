//go:build linux

package volume

import (
	"bytes"
	"context"
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

	// Output string is in form "Volume: 0.37"
	volumeDecimal := strings.TrimSpace(strings.Split(out.String(), " ")[1])
	volumePercent, err := strconv.ParseFloat(volumeDecimal, 64)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(volumePercent * 100)), nil
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

	// Output string is in form (probably) Volume: front-left: 24248 /  37% / -25.91 dB,   front-right: 24248 /  37% / -25.91 dB
	regex := regexp.MustCompile("front-left:\\s+\\d+\\s*/\\s*(\\d+)%.*?front-right:\\s+\\d+\\s*/\\s*(\\d+)%")
	matches := regex.FindStringSubmatch(out.String())

	volumeLeft, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}

	volumeRight, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", err
	}

	return strconv.Itoa((volumeLeft + volumeRight) / 2), nil
}
