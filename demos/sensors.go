// How to run:
// 	./sensors <MQTT SERVER>
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
	"gobot.io/x/gobot/platforms/mqtt"
)

var light int

func main() {
	a := edison.NewAdaptor()
	lux := aio.NewAnalogSensorDriver(a, "0")
	screen := i2c.NewGroveLcdDriver(a)

	server := mqtt.NewAdaptor(os.Args[1], "sensors")

	work := func() {
		lux.On(aio.Data, func(data interface{}) {
			light = data.(int)
		})

		gobot.Every(500*time.Millisecond, func() {
			data := make([]byte, 1)
			data[0] = byte(light)
			server.Publish("miniluminado/sensors/light", data)
		})

		gobot.Every(250*time.Millisecond, func() {
			screen.Clear()
			screen.Home()

			msg := fmt.Sprintf("%d", light)
			screen.Write(msg)

			if light > 200 {
				screen.SetRGB(0, 255, 0) // green
			} else if light > 100 {
				screen.SetRGB(0, 0, 255) // blue
			} else {
				screen.SetRGB(255, 0, 0) // red
			}
		})
	}

	robot := gobot.NewRobot("sensors",
		[]gobot.Connection{a, server},
		[]gobot.Device{lux, screen},
		work,
	)

	robot.Start()
}
