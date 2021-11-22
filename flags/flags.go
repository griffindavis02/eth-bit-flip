package flags

import (
	"github.com/griffindavis02/eth-bit-flip/config"
	"gopkg.in/urfave/cli.v1"
)

var (
	// Flags for simulating soft errors in the blockchain
	// TODO: Add flag declarations to main.go
	FlipInitialized = cli.BoolFlag{
		Name: "flip.initialized",
		Usage: "Whether the test configuration was set up with the CLI",
	}

	FlipPath = cli.StringFlag{
		Name:  "flip.path",
		Usage: "Path to soft error configuration file",
		Value: config.Path,
	}

	FlipStart = cli.BoolFlag{
		Name:  "flip.start",
		Usage: "Start the soft error simulation",
	}

	FlipStop = cli.BoolFlag{
		Name:  "flip.stop",
		Usage: "Stop the soft error simulation",
	}

	FlipRestart = cli.BoolFlag{
		Name:  "flip.restart",
		Usage: "Restart the soft error simulation, discarding current results",
	}

	FlipType = cli.StringFlag{
		Name:  "flip.test_type",
		Usage: "Determines how to increment test counter",
		Value: config.DefaultConfig.State.TestType,
	}

	FlipCounter = cli.IntFlag{
		Name:  "flip.test_counter",
		Usage: "Counter for iteration and variable based tests",
		Value: config.DefaultConfig.State.TestCounter,
	}

	FlipIterations = cli.IntFlag{
		Name:  "flip.iterations",
		Usage: "Number of iterations to run through per error rate",
		Value: config.DefaultConfig.State.Iterations,
	}

	FlipVariables = cli.IntFlag{
		Name:  "flip.variables_changed",
		Usage: "Number of variables to change per error rate",
		Value: config.DefaultConfig.State.VariablesChanged,
	}

	FlipDuration = cli.DurationFlag{
		Name:  "flip.duration",
		Usage: "How long per error rate to run a test for",
		Value: config.DefaultConfig.State.Duration,
	}

	FlipTime = cli.Int64Flag{
		Name:  "flip.start_time",
		Usage: "Start time for each error rate's test",
		Value: config.DefaultConfig.State.StartTime,
	}

	FlipRate = cli.IntFlag{
		Name:  "flip.rate_index",
		Usage: "Index of rate within the flip.error_rates flag",
		Value: config.DefaultConfig.State.RateIndex,
	}

	FlipRates = cli.StringFlag{
		Name:  "flip.error_rates",
		Usage: "String of error rates to iterate through",
		Value: config.DefaultConfig.State.ErrorRates,
	}

	FlipPost = cli.BoolFlag{
		Name:  "flip.post",
		Usage: "Whether or not to post results to host in flip.host flag",
	}

	FlipHost = cli.StringFlag{
		Name:  "flip.host",
		Usage: "Host/API to push flip results to",
		Value: config.DefaultConfig.Server.Host,
	}
)

func FlagtoConfig(ctx *cli.Context) config.Config {
	var cfg config.Config
	if initialized := ctx.GlobalBool(FlipInitialized.Name); initialized {
		cfg.Initialized = initialized

		cfg.State.TestType = ctx.GlobalString(FlipType.Name)
		cfg.State.TestCounter = ctx.GlobalInt(FlipCounter.Name)
		cfg.State.Iterations = ctx.GlobalInt(FlipIterations.Name)
		cfg.State.VariablesChanged = ctx.GlobalInt(FlipIterations.Name)
		cfg.State.Duration = ctx.GlobalDuration(FlipDuration.Name)
		cfg.State.StartTime = ctx.GlobalInt64(FlipTime.Name)
		cfg.State.RateIndex = ctx.GlobalInt(FlipRate.Name)
		cfg.State.ErrorRates = ctx.GlobalString(FlipRates.Name)

		cfg.Server.Post = ctx.GlobalBool(FlipPost.Name)
		cfg.Server.Host = ctx.GlobalString(FlipHost.Name)

		return cfg
	}
	return config.Config{}
}