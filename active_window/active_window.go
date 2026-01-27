package active_window

import (
	"KIN/config"
	"KIN/kbhid"
	"log"
	"time"
)

func SendActiveWindowData() {
	window := FetchActiveWindowName()

	if window != "" {

		keyboards := config.PayloadToKeyboards["volume"]

		for _, cfg := range keyboards {

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
}
