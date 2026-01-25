package volume

import (
	"KIN/kbhid"
	"context"
	"log"
	"time"

	"github.com/sstallion/go-hid"
)

func SendVolumeData(device *hid.Device, cfg kbhid.KeyboardHIDInfo) {
	volume, _ := FetchVolume(context.Background())

	if volume != "" {
		data := kbhid.PrepareCStringPayload(
			volume,
			cfg.ReportLength-1,
		)

		payload := kbhid.BuildPayload(
			kbhid.PayloadVolume,
			data,
			cfg.ReportLength,
		)

		if err := kbhid.SendRawReport(device, cfg, payload); err != nil {
			log.Printf("write failed: %v", err)
		}
	}

	time.Sleep(time.Millisecond * 100)
}
