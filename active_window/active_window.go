package active_window

import (
	"KIN/kbhid"
	"log"
	"time"

	"github.com/sstallion/go-hid"
)

func SendActiveWindowData(device *hid.Device, cfg kbhid.KeyboardHIDInfo) {
	window := FetchActiveWindowName()

	if window != "" {
		data := kbhid.PrepareCStringPayload(
			window,
			cfg.ReportLength-1,
		)

		payload := kbhid.BuildPayload(
			kbhid.PayloadActiveWindow,
			data,
			cfg.ReportLength,
		)

		if err := kbhid.SendRawReport(device, cfg, payload); err != nil {
			log.Printf("write failed: %v", err)
		}
	}

	time.Sleep(time.Second)

}
