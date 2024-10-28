package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/solver"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile(time.Now().Format(time.RFC3339)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic("failed to open log file")
	}
	log.SetOutput(file)
}

const (
	sensorsIP = "localhost:8080"
	motorsIP  = sensorsIP
	robotID   = "F535AF9628574A53"
)

func main() {
	dummy := flag.Bool("dummy", false, "")
	backward := flag.Bool("bw", false, "")
	sip := flag.String("sip", sensorsIP, "")
	mip := flag.String("bip", motorsIP, "")
	id := flag.String("id", robotID, "")
	flag.Parse()

	fmt.Println("flags", *sip, *mip, *id, dummy, backward)

	var mover mo.Mover
	if *dummy {
		mover = mo.NewDummyMover(log.WithField("entity", "dummy-mover"), *sip, *mip, *id)
	} else {
		mover = mo.NewSmartMover(*sip, *mip, *id)
		//mover = mo.NewSmartMover(log.WithField("entity", "smart-mover"), *sip, *mip, *id)
	}

	config := solver.FloodFillConfig{
		StartDirection:  maze.Up,
		StartPosition:   solver.NewPosition(0, 0),
		MoveForwardOnly: !*backward,
		Mover:           mover,
		Logger:          log.WithField("entity", "flood-fill"),
	}

	ff := solver.NewFloodFill(config)
	ff.Solve()
}
