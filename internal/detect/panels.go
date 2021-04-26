package detect

import (
	"context"
	"errors"
	"time"

	"github.com/landru29/flightpanel/internal/device"
	"github.com/sirupsen/logrus"
)

func Panels(ctx context.Context, log logrus.FieldLogger) (map[device.PanelType]*device.Panel, error) {
	res := map[device.PanelType]*device.Panel{}

	k, err := listSerial(log, SerialPattern)
	if err != nil {
		return res, err
	}

	for _, name := range k {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		dev, err := identify(ctx, log, name)
		cancel()
		if err != nil {
			log.WithError(err).Error("fail to identify device")
			if dev.Reader != nil {
				_ = dev.Close(log)
			}
			continue
		}

		res[device.PanelType(dev.Type)] = dev
	}

	if len(res) == 0 {
		return res, errors.New("no device found")
	}

	return res, nil
}
