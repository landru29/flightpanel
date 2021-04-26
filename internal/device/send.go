package device

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func Send(ctx context.Context, lg logrus.FieldLogger, p map[PanelType]*Panel) error {
	ticker := time.NewTicker(400 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			for _, dev := range p {
				_ = dev.Close(lg)
			}
			return nil
		case <-ticker.C:
			for _, dev := range p {
				for i := range dev.Value {
					dev.SendValue(lg.WithFields(logrus.Fields{
						"name": dev.Name,
						"type": dev.Type,
					}), i)
				}
			}
		}
	}
}
