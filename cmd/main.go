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
}

const (
	sensorsIP = "localhost:8080"
	motorsIP  = sensorsIP
	robotID   = "F535AF9628574A53"
)

func main() {
	dummy := flag.Bool("dummy", false, "")
	backward := flag.Bool("bw", false, "")
	stdLogs := flag.Bool("std", false, "")
	sip := flag.String("sip", sensorsIP, "")
	mip := flag.String("bip", motorsIP, "")
	id := flag.String("id", robotID, "")
	flag.Parse()

	fmt.Printf(
		"is_dummy=\t%v\n"+
			"backward=\t%v\n"+
			"logs_to_std=\t%v\n"+
			"sensors_ip=\t%s\n"+
			"motors_ip=\t%s\n"+
			"robot_id=\t%s\n",
		*dummy, *backward, *stdLogs, *sip, *mip, *id,
	)
	if *stdLogs {
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(time.Now().Format(time.RFC3339)+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Panic("failed to open log file")
		}
		log.SetOutput(file)
	}

	var mover mo.Mover
	if *dummy {
		mover = mo.NewDummyMover(log.WithField("entity", "dummy-mover"), *sip, *mip, *id)
	} else {
		mover = mo.NewSmartMover(log.WithField("entity", "smart-mover"), *sip, *mip, *id)
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
