package app

type PayloadType uint8

const (
	PayloadUnknown PayloadType = iota
	PayloadActiveWindow
	PayloadVolume
)

func BuildPayload(payloadType PayloadType, data []byte, payloadSize int) []byte {
	payload := make([]byte, payloadSize)
	payload[0] = byte(payloadType)

	maxDataLen := payloadSize - 1
	if len(data) > maxDataLen {
		data = data[:maxDataLen]
	}

	copy(payload[1:], data)
	return payload
}

func PrepareCString(s string, maxLen int) []byte {
	data := []byte(s + "\x00")

	if len(data) > maxLen {
		data = data[:maxLen]
		data[maxLen-1] = 0
	}

	return data
}
