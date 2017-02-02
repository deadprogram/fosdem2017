// How to run:
// 	./base <MQTT SERVER>
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/mqtt"
)

var board *firmata.Adaptor
var leds *gpio.RgbLedDriver

var server *mqtt.Adaptor
var light *mqtt.Driver
var drone *mqtt.Driver

var lightLevel int
var blinker *time.Ticker
var blinking bool

const (
	Landing = iota
	Takeoff = iota
)

func main() {
	master := gobot.NewMaster()
	a := api.NewAPI(master)
	a.Start()

	board = firmata.NewAdaptor("/dev/ttyACM0")
	leds = gpio.NewRgbLedDriver(board, "3", "6", "5")

	server = mqtt.NewAdaptor(os.Args[1], "basestation")
	light = mqtt.NewDriver(server, "miniluminado/sensors/light")
	drone = mqtt.NewDriver(server, "miniluminado/drones/flights")

	work := func() {
		light.On(mqtt.Data, func(data interface{}) {
			lightLevel = int(data.([]byte)[0])
			fmt.Println("Light:", lightLevel)
			if !blinking {
				displayLightLevel()
			}
		})

		drone.On(mqtt.Data, func(data interface{}) {
			flightStatus := int(data.([]byte)[0])
			if flightStatus == Takeoff {
				startBlinking()
			} else {
				stopBlinking()
			}
		})
	}

	robot := gobot.NewRobot("basestation",
		[]gobot.Connection{board, server},
		[]gobot.Device{leds, light, drone},
		work,
	)

	master.AddRobot(robot)

	master.Start()
}

func displayLightLevel() {
	if lightLevel > 200 {
		leds.SetRGB(0, 255, 0) // green
		return
	}

	if lightLevel > 100 {
		leds.SetRGB(0, 0, 255) // blue
		return
	}

	leds.SetRGB(255, 0, 0) // red
	return
}

func startBlinking() {
	if blinking {
		return
	}

	on := false
	blinking = true
	blinker = gobot.Every(500*time.Millisecond, func() {
		if on {
			leds.Off()
			on = false
		} else {
			displayLightLevel()
			on = true
		}
	})
}

func stopBlinking() {
	blinker.Stop()
	blinking = false
}
