// How to run:
//  sudo ./drone <NAME OF DRONE> <MQTT SERVER>
package main

import (
	"os"
	"sync/atomic"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/joystick"
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot/platforms/parrot/minidrone"
)

const (
	Landing = 1
	Takeoff = 2
)

type pair struct {
	x float64
	y float64
}

var leftX, leftY, rightX, rightY atomic.Value

const offset = 32767.0

func main() {
	pwd, _ := os.Getwd()
	joystickConfig := pwd + "/demos/dualshock3.json"

	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, joystickConfig)

	droneAdaptor := ble.NewClientAdaptor(os.Args[1])
	drone := minidrone.NewDriver(droneAdaptor)

	server := mqtt.NewAdaptor(os.Args[2], "drone")

	work := func() {
		leftX.Store(float64(0.0))
		leftY.Store(float64(0.0))
		rightX.Store(float64(0.0))
		rightY.Store(float64(0.0))

		stick.On(joystick.CirclePress, func(data interface{}) {
			drone.FrontFlip()
		})

		stick.On(joystick.TrianglePress, func(data interface{}) {
			drone.HullProtection(true)
			drone.TakeOff()

			msg := make([]byte, 1)
			msg[0] = byte(Takeoff)
			server.Publish("miniluminado/drones/flights", msg)
		})

		stick.On(joystick.SquarePress, func(data interface{}) {
			drone.Stop()
		})

		stick.On(joystick.XPress, func(data interface{}) {
			drone.Land()

			msg := make([]byte, 1)
			msg[0] = byte(Landing)
			server.Publish("miniluminado/drones/flights", msg)
		})

		stick.On(joystick.LeftX, func(data interface{}) {
			val := float64(data.(int16))
			leftX.Store(val)
		})

		stick.On(joystick.LeftY, func(data interface{}) {
			val := float64(data.(int16))
			leftY.Store(val)
		})

		stick.On(joystick.RightX, func(data interface{}) {
			val := float64(data.(int16))
			rightX.Store(val)
		})

		stick.On(joystick.RightY, func(data interface{}) {
			val := float64(data.(int16))
			rightY.Store(val)
		})

		gobot.Every(10*time.Millisecond, func() {
			rightStick := getRightStick()

			switch {
			case rightStick.y < -10:
				drone.Forward(minidrone.ValidatePitch(rightStick.y, offset))
			case rightStick.y > 10:
				drone.Backward(minidrone.ValidatePitch(rightStick.y, offset))
			default:
				drone.Forward(0)
			}

			switch {
			case rightStick.x > 10:
				drone.Right(minidrone.ValidatePitch(rightStick.x, offset))
			case rightStick.x < -10:
				drone.Left(minidrone.ValidatePitch(rightStick.x, offset))
			default:
				drone.Right(0)
			}
		})

		gobot.Every(10*time.Millisecond, func() {
			leftStick := getLeftStick()
			switch {
			case leftStick.y < -10:
				drone.Up(minidrone.ValidatePitch(leftStick.y, offset))
			case leftStick.y > 10:
				drone.Down(minidrone.ValidatePitch(leftStick.y, offset))
			default:
				drone.Up(0)
			}

			switch {
			case leftStick.x > 20:
				drone.Clockwise(minidrone.ValidatePitch(leftStick.x, offset))
			case leftStick.x < -20:
				drone.CounterClockwise(minidrone.ValidatePitch(leftStick.x, offset))
			default:
				drone.Clockwise(0)
			}
		})
	}

	robot := gobot.NewRobot("minidrone",
		[]gobot.Connection{joystickAdaptor, droneAdaptor, server},
		[]gobot.Device{stick, drone},
		work,
	)

	robot.Start()
}

func getLeftStick() pair {
	s := pair{x: 0, y: 0}
	s.x = leftX.Load().(float64)
	s.y = leftY.Load().(float64)
	return s
}

func getRightStick() pair {
	s := pair{x: 0, y: 0}
	s.x = rightX.Load().(float64)
	s.y = rightY.Load().(float64)
	return s
}
