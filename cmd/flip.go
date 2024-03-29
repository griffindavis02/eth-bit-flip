package cmd

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/griffindavis02/eth-bit-flip/config"
	"gopkg.in/urfave/cli.v1"
)

var (
	flipCommand = cli.Command{
		Action:   utils.MigrateFlags(setStatus),
		Name:     "flip",
		Usage:    "Handle the bit flipping environment set up with the 'flipconfig' CLI",
		Flags:    []cli.Flag{utils.FlipStart, utils.FlipStop, utils.FlipRestart}, // change to flip flags
		Category: "CONSOLE COMMANDS",
		Description: `
This command allows you to start soft error simulation on a running node. See
https://github.com/griffindavis02/eth-bit-flip
		`,
	}
)

func setStatus(ctx *cli.Context) error {
	var cfg  config.Config

	cfg, err := config.ReadConfig()
	if err != nil {
		utils.Fatalf("Failed to read in the error injection confguration file:", err)
	}

	if ctx.GlobalIsSet(utils.FlipStart.Name) {
		cfg.Start = true
		cfg.Restart = false
	}

	if ctx.GlobalIsSet(utils.FlipRestart.Name) {
		cfg.Restart = true
		cfg.Start = true
	}

	if ctx.GlobalIsSet(utils.FlipStop.Name) {
		cfg.Start = false
		cfg.Restart = false
	}

	(&cfg).WriteConfig()

	return nil
}
