// TODO: Add Licensure

package main

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

var (
	enableFlag = cli.BoolFlag{
		Name: "flip.enable",
		Usage: "Enable soft error simulation",
	}
	disableFlag = cli.BoolFlag{
		Name: "flip.disable",
		Usage: "Disable soft error simulation",
	}
	resetFlag = cli.BoolFlag{
		Name: "flip.reset",
		Usage: "Reset state variables if an instance was quit during testing",
	}

	flipCommand = cli.Command {
		Action: utils.MigrateFlags(flipInit),
		Name: "flip",
		Usage: "Set up a soft error test environment for geth",
		Flags: []cli.Flag{enableFlag,disableFlag,resetFlag,},
		Category: "SOFT ERROR INJECTION",
		Description: `
		The flip command allows one to configure the parameters by which to simulate
		soft errors in the EVM.
		
		usage: geth flip [--enable] [--disable] [--reset]
		Lack of flag will begin the configuration wizard.
		`,
	}
)

// flipInit starts a wizard for defining soft error test variables, or either
// enables, disables, or resets with a pre-configured test environment.
func flipInit(ctx *cli.Context) {
	// TODO: Add error return
	// FIXME: Add these flags to the utils file for reference in other code
	if !ctx.GlobalIsSet(enableFlag.Name) && !ctx.GlobalIsSet(disableFlag.Name) && !ctx.GlobalIsSet(resetFlag.Name) {
		flipWizard(ctx)
	} else if ctx.GlobalIsSet(resetFlag.Name) && ctx.GlobalBool(resetFlag.Name) {
		// TODO: reset code
	} else if ctx.GlobalIsSet(disableFlag.Name) && ctx.GlobalBool(disableFlag.Name) {
		// TODO: d
	}
	// TODO: else we should log the command isn't being used properly
}

// TODO: populate
func flipWizard(ctx *cli.Context) {

}