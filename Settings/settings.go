package settings

import "fmt"

type AppConfig struct {
	Name       string
	Debug      bool
	AppVersion string
	PcapSize   int //Bytes
}

var Config = AppConfig{
	Name:       "Heimdall",
	Debug:      false,
	AppVersion: "v0.2",
	PcapSize:   1 * (1000 * 1024), //MB
}

func (c AppConfig) PrintCLi() {
	fmt.Println("- Name:", c.Name)
	fmt.Println("- Debug:", c.Debug)
	fmt.Println("- App version:", c.AppVersion)
	fmt.Println("- PcapSize:", (float32(c.PcapSize) / 1000 / 1024), "MB")
}
