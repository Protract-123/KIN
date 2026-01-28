package app

type PayloadType uint8

const (
	PayloadUnknown PayloadType = iota
	PayloadActiveWindow
	PayloadVolume
)

// BuildPayload creates a payload based on the payload type, the data and the
// max payload size which can be sent
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

// PrepareCString converts a string into a null-terminated CString
func PrepareCString(s string, maxLen int) []byte {
	data := []byte(s + "\x00")

	if len(data) > maxLen {
		data = data[:maxLen]
		data[maxLen-1] = 0
	}

	return data
}
