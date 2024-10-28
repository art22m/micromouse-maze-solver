package main

import (
	"flag"
	"fmt"

	"jackson/internal/maze"
	mo "jackson/internal/mover"
	"jackson/internal/solver"
)

const (
	sensorsIP = "192.168.68.151"
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
	fmt.Println(*sip, *mip, *id, dummy, backward)

	var mover mo.Mover
	if *dummy {
		mover = mo.NewDummyMover(*sip, *mip, *id)
	} else {
		mover = mo.NewSmartMover(*sip, *mip, *id)
	}

	config := solver.FloodFillConfig{
		StartDirection:  maze.Up,
		StartPosition:   solver.NewPosition(0, 0),
		MoveForwardOnly: !*backward,
		Mover:           mover,
	}

	ff := solver.NewFloodFill(config)
	ff.Solve()
}
