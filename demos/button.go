package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")

	led := gpio.NewLedDriver(firmataAdaptor, "2")
	button := gpio.NewButtonDriver(firmataAdaptor, "3")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led, button},
		work,
	)

	robot.Start()
}
