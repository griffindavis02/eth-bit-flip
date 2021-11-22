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

	mathEth "github.com/ethereum/go-ethereum/common/math"
	"github.com/griffindavis02/eth-bit-flip/config"
	"github.com/griffindavis02/eth-bit-flip/flags"
	"gopkg.in/urfave/cli.v1"
)

type ErrorData struct {
	PreviousValue *big.Int
	PreviousByte  string
	IntBits       []int
	ErrorValue    *big.Int
	ErrorByte     string
	DeltaValue    *big.Int
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
		postAPI(cfg.Server.Host, boiler)
	}
}

// BitFlip will run the odds of flipping a bit within pbigNum based on error
// rate pdecRate. The iteration count will increment and both the new number
// and the iteration error data will be returned.
func BitFlip(pbigNum *big.Int, ctx *cli.Context) *big.Int {
	if !ctx.GlobalBool(flags.FlipStart.Name) || ctx.GlobalBool(flags.FlipStop.Name) {
		return pbigNum
	}
	
	initalize(ctx)
	rand.Seed(time.Now().UnixNano())

	// Check for out of bounds or end of error rate
	switch cfg.State.TestType {
	case "iteration":
		if cfg.State.TestCounter >= cfg.State.Iterations {
			if cfg.State.RateIndex == len(marrErrRates)-1 {
				return pbigNum
			}
			cfg.State.RateIndex++
			cfg.State.TestCounter = 0
		}
	case "variable":
		if cfg.State.TestCounter >= cfg.State.VariablesChanged {
			if cfg.State.RateIndex == len(marrErrRates)-1 {
				return pbigNum
			}
			cfg.State.RateIndex++
			cfg.State.TestCounter = 0
		}
	default:
		if time.Since(time.Unix(cfg.State.StartTime, 0)) >= cfg.State.Duration {
			if cfg.State.RateIndex == len(marrErrRates)-1 {
				return pbigNum
			}
			cfg.State.RateIndex++
			cfg.State.StartTime = time.Now().Unix()
		}
	}

	decRate := marrErrRates[cfg.State.RateIndex]
	var arrBits []int

	// Store previous states
	lngPrevCounter := cfg.State.TestCounter
	bigPrevNum, _ := new(big.Int).SetString(pbigNum.String(), 10)
	bigPrevNum = mathEth.U256(bigPrevNum)
	bytPrevNum := bigPrevNum.Bytes()
	bytNum := pbigNum.Bytes()
	intLastByte := len(bytNum) - 1

	// Run chance of flipping a bit in byte representation
	for i := range bytNum {
		for j := 0; j < 8; j++ {
			if math.Floor(rand.Float64()/decRate) == math.Floor(rand.Float64()/decRate) {
				if cfg.State.TestType == "iteration" {
					cfg.State.TestCounter++
				}
				arrBits = append(arrBits, (i*8)+j)
				bytNum[intLastByte-i] ^= (1 << j)
			}
		}
	}

	// Ensure there was a change
	if !bytes.Equal(bytNum, bytPrevNum) {
		if cfg.State.TestType == "variable" {
			cfg.State.TestCounter++
		}
		// Recreate number from byte code
		pbigNum.SetBytes(bytNum)
		// Build error data
		iteration := Iteration{
			marrErrRates[cfg.State.RateIndex],
			int(lngPrevCounter),
			ErrorData{
				bigPrevNum,
				"0x" + hex.EncodeToString(bytPrevNum),
				arrBits,
				pbigNum,
				"0x" + hex.EncodeToString(bytNum),
				big.NewInt(0).Sub(pbigNum, bigPrevNum),
				time.Now().Format("01-02-2006 15:04:05.000000000"),
			},
		}

		// Pretty print JSON in console and append to error rate data
		bytJSON, _ := json.MarshalIndent(iteration, "", "    ")
		fmt.Println(string(bytJSON))
		if cfg.Server.Post {
			postAPI(cfg.Server.Host, iteration)
		}
	}

	return pbigNum
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
