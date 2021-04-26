package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/landru29/flightpanel/internal/detect"
	"github.com/landru29/flightpanel/internal/device"
	"github.com/sirupsen/logrus"
)

func main() {
	lg := logrus.New()

	ctx, cancelFunc := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		lg.Infof("exiting program")
		cancelFunc()
	}()

	devices, err := detect.Panels(ctx, lg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		device.Listen(ctx, lg, devices)
	}()

	device.Send(ctx, lg, devices)
}
