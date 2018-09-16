package weather

import (
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type PeriphLed struct {
	pin     string
	gpioPin gpio.PinIO
	state   bool
}

func NewPeriphLed(pin string) Led {
	gpioPin := gpioreg.ByName(pin)
	return &PeriphLed{
		pin:     pin,
		gpioPin: gpioPin,
		state:   false,
	}
}

func (l *PeriphLed) On() {
	l.gpioPin.Out(gpio.High)
	l.state = true
}

func (l *PeriphLed) Off() {
	l.gpioPin.Out(gpio.Low)
	l.state = false
}

func (l *PeriphLed) GetCurrentState() bool {
	return l.state
}

func (l *PeriphLed) Close() error {
	l.Off()
	return l.gpioPin.Halt()
}
