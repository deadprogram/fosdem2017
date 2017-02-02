package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	a := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(a, "13")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{a},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
