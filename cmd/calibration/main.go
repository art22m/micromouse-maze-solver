package main

import (
	"flag"
	"jackson/internal/mover"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	sensorsIP = "192.168.68.202"
	motorsIP  = sensorsIP
	robotID   = "7536AF961D784A53"
)

func main() {
	sip := flag.String("sip", sensorsIP, "")
	mip := flag.String("bip", motorsIP, "")
	id := flag.String("id", robotID, "")

	log := logrus.New()

	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetOutput(os.Stdout)

	sm := mover.NewSmartMover(log.WithField("entity", "smart-mover-calibration"), *sip, *mip, *id)
	sm.Calibrate()
}
