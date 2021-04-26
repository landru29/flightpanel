package device

import (
	"context"
	"errors"
	"fmt"
	"math"

	sim "github.com/micmonay/simconnect"
	"github.com/sirupsen/logrus"
)

func Listen(ctx context.Context, lg logrus.FieldLogger, p map[PanelType]*Panel) error {
	sc, err := sim.NewEasySimConnect()
	if err != nil {
		return err
	}
	sc.SetLoggerLevel(sim.LogInfo) // It is better if you the set before connect
	c, err := sc.Connect("MyApp")
	if err != nil {
		return err
	}
	<-c // Wait connection confirmation
	cSimVar, err := sc.ConnectToSimVar(
		sim.SimVarIndicatedAltitude(),
		sim.SimVarGeneralEngRpm(1),
		sim.SimVarGyroDriftError(),
		sim.SimVarMagneticCompass(),
		sim.SimVarAirspeedIndicated(),
	)
	if err != nil {
		return err
	}
	cSimStatus := sc.ConnectSysEventSim()
	//wait sim start
	for {
		if <-cSimStatus {
			break
		}
	}
	crashed := sc.ConnectSysEventCrashed()

	currentGyroDrift := float64(0)

	for {
		select {
		case result := <-cSimVar:
			for _, simVar := range result {
				f, err := simVar.GetFloat64()
				if err != nil {
					fmt.Println("return error :", err)
				}
				switch simVar.Name {
				case "GYRO DRIFT ERROR":
					currentGyroDrift = f * 180 / math.Pi
				case "MAGNETIC COMPASS":
					if dev, ok := p[PanelGyro]; ok {
						dev.SetValue(lg.WithField("type", PanelGyro), 0, float32(f+currentGyroDrift))
					}
				case "INDICATED ALTITUDE":
					if dev, ok := p[PanelAltitude]; ok {
						dev.SetValue(lg.WithField("type", PanelAltitude), 0, float32(f))
					}
				case "GENERAL ENG RPM:index":
					if dev, ok := p[PanelEngine]; ok {
						dev.SetValue(lg.WithField("type", PanelEngine), 0, float32(f))
					}
				case "AIRSPEED INDICATED":
					if dev, ok := p[PanelAirSpeed]; ok {
						dev.SetValue(lg.WithField("type", PanelAirSpeed), 0, float32(f))
					}
				}
			}

		case <-crashed:
			<-sc.Close()                 // Wait close confirmation
			return errors.New("crashed") // This example close after crash in the sim

		case <-ctx.Done():
			<-sc.Close()
			return nil
		}

	}
}
