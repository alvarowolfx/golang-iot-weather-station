package weather

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/devices/tm1637"
)

type TM1637Display struct {
	clkPin  string
	dataPin string
	dev     *tm1637.Dev
}

func NewTM1637Display(clkPin, dataPin string) (Display, error) {
	d := &TM1637Display{
		clkPin:  clkPin,
		dataPin: dataPin,
	}
	err := d.init()
	return d, err
}

func (d *TM1637Display) init() error {
	clk := gpioreg.ByName(d.clkPin)
	data := gpioreg.ByName(d.dataPin)
	if clk == nil || data == nil {
		return errors.New("Failed to find tm1637 pins")
	}

	dev, err := tm1637.New(clk, data)
	if err != nil {
		return err
	}
	dev.SetBrightness(tm1637.Brightness10)

	d.dev = dev

	return nil
}

func (d *TM1637Display) Display(value float64) error {
	valueStr := fmt.Sprintf("%2.2f", value)

	numbers := strings.Split(valueStr, ".")
	firstPart, _ := strconv.Atoi(numbers[0])
	lastPart := 0
	if len(numbers) > 1 {
		lastPart, _ = strconv.Atoi(numbers[1])
	}

	_, err := d.dev.Write(tm1637.Clock(firstPart, lastPart, true))

	return err
}

func (d *TM1637Display) Close() error {
	return d.dev.Halt()
}
