// Copyright 2021 The eth-bit-flip Authors
// This file is part of the eth-bit-flip library.
//
// The eth-bit-flip libary is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// The eth-bit-flip libary is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// the eth-bit-flip library. If not, see <https://www.gnu.org/licenses/>.

package injection

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"time"

	"github.com/griffindavis02/eth-bit-flip/config"
	"github.com/griffindavis02/eth-bit-flip/flags"
	"gopkg.in/urfave/cli.v1"
)

type ErrorData struct {
	PreviousValue interface{}
	PreviousByte  string
	IntBits       []int
	ErrorValue    interface{}
	ErrorByte     string
	DeltaValue    interface{}
	When          string
}

type Iteration struct {
	Rate float64
	IterationNum int
	ErrorData    ErrorData
}

type ErrorRate struct {
	Rate     float64
	FlipData []Iteration
}

type Output struct {
	Data []ErrorRate
}

var (
	cfg config.Config
	marrErrRates []float64
)

// Set up the testing environment with the test type, number of
// changes/iterations or duration in seconds, and error rates. This is
// PER error rate. i.e. 5 minutes and ten error rates will be 50 minutes.
// Test types:
// 'iteration' - increments for each bit flipped
// 'variable' - increments for each variable, regardless of bits flipped
// 'time' - checks against passage of time since started
func initalize(ctx *cli.Context) {
	cfg = flags.FlagtoConfig(ctx)
	var (
		boiler Output
		flipData []Iteration
	)
	marrErrRates, _ = config.AtoF64Arr(ctx.GlobalString(flags.FlipRates.Name))
	for _, decErrRate := range marrErrRates {
		Rate := ErrorRate{decErrRate, flipData}
		boiler.Data = append(boiler.Data, Rate)
	}
	if cfg.Server.Post {
		postAPI(cfg.Server.Host + "/initialize", boiler)
	}
}

// BitFlip will run the odds of flipping a bit within pbigNum based on error
// rate pdecRate. The iteration count will increment and both the new number
// and the iteration error data will be returned.
// TODO : use interface to accept multiple data types
func BitFlip(pIFlipee interface{}, ctx *cli.Context) interface{} {
	if !ctx.GlobalBool(flags.FlipStart.Name) || ctx.GlobalBool(flags.FlipStop.Name) {
		return pIFlipee
	}

	// Check for out of bounds or end of error rate
	switch cfg.State.TestType {
	case "iteration":
		if cfg.State.TestCounter >= cfg.State.Iterations {
			if cfg.State.RateIndex == len(marrErrRates)-1 {
				ctx.GlobalSet(flags.FlipStop.Name, "true")
			}
			cfg.State.RateIndex++
			cfg.State.TestCounter = 0
		}
	case "variable":
		if cfg.State.TestCounter >= cfg.State.VariablesChanged {
			if cfg.State.RateIndex == len(marrErrRates)-1 {
				ctx.GlobalSet(flags.FlipStop.Name, "true")
			}
			cfg.State.RateIndex++
			cfg.State.TestCounter = 0
		}
	default:
		if time.Since(time.Unix(cfg.State.StartTime, 0)) >= cfg.State.Duration {
			if cfg.State.RateIndex == len(marrErrRates)-1 {
				ctx.GlobalSet(flags.FlipStop.Name, "true")
			}
			cfg.State.RateIndex++
			cfg.State.StartTime = time.Now().Unix()
		}
	}

	initalize(ctx)
	rand.Seed(time.Now().UnixNano())

	switch pIFlipee.(type) {
	case string:
		iteration := flipBytes([]byte(pIFlipee.(string)), ctx)
		iteration.ErrorData.PreviousValue = iteration.ErrorData.PreviousValue.(string)
		iteration.ErrorData.ErrorValue = iteration.ErrorData.ErrorValue.(string)
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(string)
		printOut(iteration)
		return iteration.ErrorData.ErrorValue
	case int:
		iteration := flipBytes([]byte(pIFlipee.(string)), ctx)
		iteration.ErrorData.PreviousValue = iteration.ErrorData.PreviousValue.(int)
		iteration.ErrorData.ErrorValue = iteration.ErrorData.ErrorValue.(int)
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(int)
		printOut(iteration)
		return iteration.ErrorData.ErrorValue
	default: // case *big.Int:
		iteration := flipBytes([]byte(pIFlipee.(string)), ctx)
		iteration.ErrorData.PreviousValue = *big.NewInt(0).SetBytes(iteration.ErrorData.PreviousValue.([]byte))
		iteration.ErrorData.ErrorValue = *big.NewInt(0).SetBytes(iteration.ErrorData.ErrorValue.([]byte))
		iteration.ErrorData.DeltaValue = *big.NewInt(0).SetBytes(iteration.ErrorData.DeltaValue.([]byte))
		printOut(iteration)
		return iteration.ErrorData.ErrorValue
	}
}

func printOut(iteration Iteration) {
	bytJSON, _ := json.MarshalIndent(iteration, "", "    ")
		fmt.Println(string(bytJSON))
		if cfg.Server.Post {
			postAPI(cfg.Server.Host, iteration)
		}
}

func postAPI(url string, jsonOut interface{}) int {
	client := http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	query := req.URL.Query()
	params, _ := json.Marshal(jsonOut)
	query.Add("params", string(params))
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return res.StatusCode
}

func flipBytes(pbytFlipee []byte, ctx *cli.Context) Iteration {
	decRate := marrErrRates[cfg.State.RateIndex]
	var arrBits []int
	var iteration Iteration

	// Store previous states
	lngPrevCounter := cfg.State.TestCounter
	var bytPrevFlipee = pbytFlipee
	intLastByte := len(pbytFlipee) - 1

	// Run chance of flipping a bit in byte representation
	for i := range bytPrevFlipee {
		for j := 0; j < 8; j++ {
			if math.Floor(rand.Float64()/decRate) == math.Floor(rand.Float64()/decRate) {
				if cfg.State.TestType == "iteration" {
					cfg.State.TestCounter++
				}
				arrBits = append(arrBits, (i*8)+j)
				pbytFlipee[intLastByte-i] ^= (1 << j)
			}
		}
	}

	// Ensure there was a change
	if !bytes.Equal(pbytFlipee, bytPrevFlipee) {
		if cfg.State.TestType == "variable" {
			cfg.State.TestCounter++
		}
		// Build error data
		iteration = Iteration{
			marrErrRates[cfg.State.RateIndex],
			int(lngPrevCounter),
			ErrorData{
				bytPrevFlipee,
				"0x" + hex.EncodeToString(bytPrevFlipee),
				arrBits,
				pbytFlipee,
				"0x" + hex.EncodeToString(bytPrevFlipee),
				big.NewInt(0).Sub(big.NewInt(0).SetBytes(pbytFlipee),
					big.NewInt(0).SetBytes(bytPrevFlipee)),
				time.Now().Format("01-02-2006 15:04:05.000000000"),
			},
		}
	}
	
	return iteration
}