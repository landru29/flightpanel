package detect

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"sync"

	"errors"

	"github.com/jacobsa/go-serial/serial"
	"github.com/landru29/flightpanel/internal/device"
	"github.com/sirupsen/logrus"
)

func identify(ctx context.Context, log logrus.FieldLogger, portName string) (*device.Panel, error) {
	options := serial.OpenOptions{
		PortName: portName,
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
	}

	ret := &device.Panel{
		Name: portName,
		Lock: sync.Mutex{},
	}

	lg := log.WithField("name", portName)

	lg.Infof("opening serial")

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		return ret, fmt.Errorf("serial.Open: %v", err)
	}

	ret.Reader = port

	result := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(port)
		lg.Infof("waiting for device")
		if scanner.Scan() {
			result <- scanner.Text()
		}
	}()

	select {
	case <-ctx.Done():
		return ret, errors.New("timeout")
	case r := <-result:
		lg.WithField("response", r).Infof("device responded")

		s := strings.Split(r, "/")

		if len(s) != 2 {
			log.Errorf("wrong handshake: %s", r)
			return ret, fmt.Errorf("wrong handshake: %s", r)
		}

		ret.ID = s[1]
		ret.Type = strings.ToUpper(s[0])

		log.WithFields(logrus.Fields{
			"id":   ret.ID,
			"type": ret.Type,
		}).Info("found device")

		switch device.PanelType(s[0]) {
		case device.PanelGyro, device.PanelAltitude, device.PanelEngine, device.PanelAirSpeed:
			ret.Value = []float32{0}
			return ret, nil
		default:
			log.WithFields(logrus.Fields{
				"id":   ret.ID,
				"type": ret.Type,
			}).Error("unknown type")
			return ret, errors.New("unknown type")
		}
	}
}
