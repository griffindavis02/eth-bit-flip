package cmd

import (
	"gopkg.in/urfave/cli.v1"
)

var (
	// Flags for simulating soft errors in the blockchain
	FlipStart = cli.BoolFlag{
		Name:  "flipstart",
		Usage: "Start the soft error simulation",
	}

	FlipStop = cli.BoolFlag{
		Name:  "flipstop",
		Usage: "Stop the soft error simulation",
	}

	FlipRestart = cli.BoolFlag{
		Name:  "fliprestart",
		Usage: "Restart the soft error simulation, discarding current results",
	}
)
