package main

import (
	"fmt"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	f := firmata.NewAdaptor("/dev/ttyACM0")
	f.Connect()

	led := gpio.NewLedDriver(f, "2")
	led.Start()
	led.Off()

	button := gpio.NewButtonDriver(f, "3")
	button.Start()

	buttonEvents := button.Subscribe()
	for event := range buttonEvents {
		fmt.Println("Event:", event.Name, event.Data)
		if event.Name == gpio.ButtonPush {
			led.Toggle()
		}
	}
}
