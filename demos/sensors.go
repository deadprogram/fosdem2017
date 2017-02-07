// How to run:
//    ./sensors <MQTT SERVER>
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
	"gobot.io/x/gobot/platforms/mqtt"
)

var light atomic.Value

//var light int = 0

func main() {
	a := edison.NewAdaptor()
	lux := aio.NewAnalogSensorDriver(a, "0")
	screen := i2c.NewGroveLcdDriver(a)

	server := mqtt.NewAdaptor(os.Args[1], "sensors")

	light.Store(int(0))

	work := func() {
		lux.On(aio.Data, func(data interface{}) {
			lvl := data.(int)
			light.Store(lvl)
		})

		gobot.Every(500*time.Millisecond, func() {
			buf := new(bytes.Buffer)
			val := light.Load().(int)
			binary.Write(buf, binary.LittleEndian, uint16(val))
			server.Publish("miniluminado/sensors/light", buf.Bytes())
		})

		gobot.Every(250*time.Millisecond, func() {
			lvl := light.Load().(int)

			screen.Clear()
			screen.Home()
			msg := fmt.Sprintf("%d", lvl)
			screen.Write(msg)

			switch {
			case lvl > 200:
				screen.SetRGB(0, 255, 0) // green
			case lvl > 100:
				screen.SetRGB(0, 0, 255) // blue
			default:
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
