package device

import (
	"encoding/binary"
	"io"
	"math"
	"sync"

	"github.com/sirupsen/logrus"
)

type PanelType string

type Panel struct {
	Name   string
	ID     string
	Type   string
	Reader io.WriteCloser
	Value  []float32
	Lock   sync.Mutex
}

const (
	PanelGyro     = PanelType("GYRO")
	PanelAltitude = PanelType("ALT")
	PanelEngine   = PanelType("RPM")
	PanelAirSpeed = PanelType("AIR")
)

func (p *Panel) Close(log logrus.FieldLogger) error {
	log.WithFields(logrus.Fields{
		"name": p.Name,
		"id":   p.ID,
		"type": p.Type,
	}).Infof("closing device")
	return p.Reader.Close()
}

func (p *Panel) SendValue(log logrus.FieldLogger, index int) error {
	buf := []byte{0, 0, 0, 0}
	p.Lock.Lock()

	val := float32(0)
	if index < len(p.Value) {
		val = p.Value[index]
	}

	p.Lock.Unlock()

	binary.LittleEndian.PutUint32(buf, math.Float32bits(val))

	_, err := p.Reader.Write(buf)
	if err != nil {
		return err
	}
	log.WithField("value", val).Debugf(p.Name)

	return nil
}

func (p *Panel) SetValue(log logrus.FieldLogger, index int, value float32) {
	p.Lock.Lock()
	toAdd := index + 1 - len(p.Value)
	if toAdd > 0 {
		for i := 0; i < toAdd; i++ {
			p.Value = append(p.Value, 0)
		}
	}
	p.Value[index] = value
	p.Lock.Unlock()
	log.WithField("value", value).Debugf(string(PanelAltitude))
}
