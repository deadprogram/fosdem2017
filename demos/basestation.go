// How to run:
//    ./basestation <MQTT SERVER>
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync/atomic"
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

var lightLevel, blinking atomic.Value
var blinker *time.Ticker

const (
	Landing = 1
	Takeoff = 2
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

	lightLevel.Store(uint16(0))
	blinking.Store(false)

	work := func() {
		light.On(mqtt.Data, func(data interface{}) {
			var lvl uint16
			buf := bytes.NewReader(data.([]byte))
			binary.Read(buf, binary.LittleEndian, &lvl)
			lightLevel.Store(lvl)
			fmt.Println("Light:", lvl)
			b := blinking.Load().(bool)
			if !b {
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
	lvl := lightLevel.Load().(uint16)
	switch {
	case lvl > 200:
		leds.SetRGB(0, 255, 0) // green
	case lvl > 100:
		leds.SetRGB(0, 0, 255) // blue
	default:
		leds.SetRGB(255, 0, 0) // red
	}
}

func startBlinking() {
	b := blinking.Load().(bool)
	if b {
		return
	}

	on := false
	blinking.Store(true)
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
	blinking.Store(false)
}
