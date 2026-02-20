package icon

import (
	_ "embed"
)

//go:embed "export/cross_icon.png"
var CrossIcon []byte

//go:embed "export/tick_icon.png"
var TickIcon []byte

//go:embed "export/quit_icon.png"
var QuitIcon []byte

//go:embed "export/config_icon.png"
var ConfigIcon []byte
