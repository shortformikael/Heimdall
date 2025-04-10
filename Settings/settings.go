package settings

import "fmt"

type AppConfig struct {
	Name       string
	AppVersion string
	Debug      bool
	PcapSize   int //Bytes
}

var Config = AppConfig{
	Name:       "Heimdall",
	AppVersion: "v0.2",
	Debug:      false,
	PcapSize:   1 * (1000 * 1024), //MB
}

func (c AppConfig) PrintCLi() {
	fmt.Println("- Name:", c.Name, c.AppVersion)
	if c.Debug {
		fmt.Println("- Debug:", "ON")
	} else {
		fmt.Println("- Debug:", "OFF")

	}
	fmt.Println("- Size of Pcaps:", (float32(c.PcapSize) / 1000 / 1024), "MB")
}
