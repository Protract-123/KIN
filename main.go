package main

import (
	"KIN/config"
)

func main() {
	//if err := hid.Init(); err != nil {
	//	log.Fatalf("hid init failed: %v", err)
	//}
	//defer hid.Exit()
	//
	//dev, err := kbhid.FindRawHIDDevice(DefaultKeyboard)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer dev.Close()
	//
	//qt6.NewQApplication(os.Args)
	//configWindow := ui.NewConfigWindow()
	//ui.CreateTray(configWindow)
	//
	//go func() {
	//	for {
	//		active_window.SendActiveWindowData(dev, DefaultKeyboard)
	//	}
	//}()
	//go func() {
	//	for {
	//		volume.SendVolumeData(dev, DefaultKeyboard)
	//	}
	//}()
	//
	//qt6.QApplication_Exec()

	err := config.InitializeConfig()
	if err != nil {
		return
	}

	err = config.LoadConfig()
	if err != nil {
		return
	}

	config.BuildPayloadToKeyboards(config.ActiveConfig)
}
