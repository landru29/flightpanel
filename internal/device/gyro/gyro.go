package giro

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"time"

	"github.com/landru29/flightpanel/internal/device"
	sim "github.com/micmonay/simconnect"
)

func openGyro(dev device.Panel) {
	for {
		process(dev.Reader)
	}
}

func process(port io.Writer) error {
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
				out := ""
				switch simVar.Name {
				case "GYRO DRIFT ERROR":
					currentGyroDrift = f * 180 / math.Pi
				case "MAGNETIC COMPASS":
					out = fmt.Sprintln("Gyro", fmt.Sprintf("=> %f\n", f+currentGyroDrift))

				case "INDICATED ALTITUDE":
					out = fmt.Sprintln("Alt", fmt.Sprintf("=> %f\n", f))

				case "GENERAL ENG RPM:index":
					out = fmt.Sprintln("Rpm", fmt.Sprintf("=> %f\n", f))

				default:
					out = fmt.Sprintln(simVar.Name, fmt.Sprintf("=> %f", f))
				}

				_, err = port.Write([]byte(out))
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}
				fmt.Printf("%s", out)
				time.Sleep(400 * time.Millisecond)
			}

		case <-crashed:
			<-sc.Close()                 // Wait close confirmation
			return errors.New("crashed") // This example close after crash in the sim
		}
	}
}
