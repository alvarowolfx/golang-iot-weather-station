package weather

import (
	"fmt"
	"strconv"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
)

type BMP280EnvironmentSensor struct {
	i2cPort string
	bus     i2c.BusCloser
	bmp     *bmxx80.Dev
	env     physic.Env
}

func NewBMP280EnvironmentSensorAddress(i2cPort string, address uint16) (EnvironmentSensor, error) {
	es := &BMP280EnvironmentSensor{
		i2cPort: i2cPort,
	}
	err := es.init(address)
	return es, err
}

func NewBMP280EnvironmentSensor(i2cPort string) (EnvironmentSensor, error) {
	return NewBMP280EnvironmentSensorAddress(i2cPort, 0x76)
}

func (es *BMP280EnvironmentSensor) init(address uint16) error {
	bus, err := i2creg.Open(es.i2cPort)

	if err != nil {
		return err
	}

	es.bus = bus

	bmp, err := bmxx80.NewI2C(bus, address, &bmxx80.DefaultOpts)
	if err != nil {
		return err
	}

	es.bmp = bmp

	envCh, err := bmp.SenseContinuous(500 * time.Millisecond)
	if err != nil {
		return err
	}

	go func() {
		for env := range envCh {
			es.env = env
		}
	}()

	return nil
}

func (es *BMP280EnvironmentSensor) ReadTemperature() (float64, error) {
	tempString := fmt.Sprintf("%8s", es.env.Temperature)
	temp, err := strconv.ParseFloat(tempString[0:6], 64)
	return temp, err
}

func (es *BMP280EnvironmentSensor) ReadPressure() (float64, error) {
	pressureString := fmt.Sprintf("%8s", es.env.Pressure)
	pressure, err := strconv.ParseFloat(pressureString[0:6], 64)
	return pressure, err
}

func (es *BMP280EnvironmentSensor) Close() error {
	err := es.bmp.Halt()
	err = es.bus.Close()
	return err
}
