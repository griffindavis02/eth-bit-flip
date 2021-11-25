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
	"encoding/binary"
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
	Rate         float64
	IterationNum int
	ErrorData    ErrorData
}

var (
	cfg config.Config
)

// BitFlip will run the odds of flipping a bit within pbigNum based on error
// rate pdecRate. The iteration count will increment and both the new number
// and the iteration error data will be returned.
// TODO : use interface to accept multiple data types
func BitFlip(pIFlipee interface{}, cfgPath string) interface{} {
	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Config initialization error: %v", err)
	}

	// Check for out of bounds or end of error rate
	switch cfg.State.TestType {
	case "bit":
		if cfg.State.TestCounter >= cfg.State.Iterations {
			if cfg.State.RateIndex == len(cfg.State.ErrorRates)-1 {
				return pIFlipee
			}
			cfg.State.RateIndex++
			cfg.State.TestCounter = 0
		}
	case "variable":
		if cfg.State.TestCounter >= cfg.State.VariablesChanged {
			if cfg.State.RateIndex == len(cfg.State.ErrorRates)-1 {
				return pIFlipee
			}
			cfg.State.RateIndex++
			cfg.State.TestCounter = 0
		}
	default:
		if time.Since(time.Unix(cfg.State.StartTime, 0)) >= cfg.State.Duration {
			if cfg.State.RateIndex == len(cfg.State.ErrorRates)-1 {
				return pIFlipee
			}
			cfg.State.RateIndex++
			cfg.State.StartTime = time.Now().Unix()
		}
	}
	rand.Seed(time.Now().UnixNano())

	var iteration Iteration
	switch pIFlipee.(type) {
	case string:
		iteration = flipBytes([]byte(pIFlipee.(string)), cfgPath)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}

		iteration.ErrorData.PreviousValue = string(iteration.ErrorData.PreviousValue.([]byte))
		iteration.ErrorData.ErrorValue = string(iteration.ErrorData.ErrorValue.([]byte))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	case int:
		var iteration Iteration
		switch binary.Size(pIFlipee.(int)) {
		case 32:
			bytInt := make([]byte, 32)
			binary.BigEndian.PutUint32(bytInt, uint32(pIFlipee.(int)))
			iteration = flipBytes(bytInt, cfgPath)
			if iteration.ErrorData.ErrorValue == nil {
				return pIFlipee
			}

			iteration.ErrorData.PreviousValue = int(binary.BigEndian.Uint32(iteration.ErrorData.PreviousValue.([]byte)))
			iteration.ErrorData.ErrorValue = int(binary.BigEndian.Uint32(iteration.ErrorData.ErrorValue.([]byte)))
			iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
		default:
			bytInt := make([]byte, 64)
			binary.BigEndian.PutUint64(bytInt, uint64(pIFlipee.(int)))
			iteration = flipBytes(bytInt, cfgPath)
			if iteration.ErrorData.ErrorValue == nil {
				return pIFlipee
			}

			iteration.ErrorData.PreviousValue = int(binary.BigEndian.Uint64(iteration.ErrorData.PreviousValue.([]byte)))
			iteration.ErrorData.ErrorValue = int(binary.BigEndian.Uint64(iteration.ErrorData.ErrorValue.([]byte)))
			iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
		}
	case *big.Int:
		iteration = flipBytes(pIFlipee.(*big.Int).Bytes(), cfgPath)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}
		iteration.ErrorData.PreviousValue = new(big.Int).SetBytes(iteration.ErrorData.PreviousValue.([]byte))
		iteration.ErrorData.ErrorValue = new(big.Int).SetBytes(iteration.ErrorData.ErrorValue.([]byte))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	}

	printOut(iteration)
	return iteration.ErrorData.ErrorValue
}

func flipBytes(pbytFlipee []byte, cfgPath string) Iteration {
	decRate := cfg.State.ErrorRates[cfg.State.RateIndex]
	var arrBits []int
	var iteration Iteration

	// Store previous states
	lngPrevCounter := cfg.State.TestCounter
	var bytPrevFlipee []byte
	bytPrevFlipee = append(bytPrevFlipee, pbytFlipee...)
	intLastByte := len(pbytFlipee) - 1

	// Run chance of flipping a bit in byte representation
	for i := range bytPrevFlipee {
		for j := 0; j < 8; j++ {
			if math.Floor(rand.Float64()/decRate) == math.Floor(rand.Float64()/decRate) {
				if cfg.State.TestType == "bit" {
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
			cfg.State.ErrorRates[cfg.State.RateIndex],
			int(lngPrevCounter),
			ErrorData{
				bytPrevFlipee,
				"0x" + hex.EncodeToString(bytPrevFlipee),
				arrBits,
				pbytFlipee,
				"0x" + hex.EncodeToString(pbytFlipee),
				big.NewInt(0).Sub(big.NewInt(0).SetBytes(pbytFlipee),
					big.NewInt(0).SetBytes(bytPrevFlipee)),
				time.Now().Format("01-02-2006 15:04:05.000000000"),
			},
		}

		config.WriteConfig(cfgPath, cfg)
	}

	return iteration
}

func printOut(iteration Iteration) {
	if iteration.ErrorData.PreviousValue == iteration.ErrorData.ErrorValue {
		return
	}
	// TODO: Look for logging boolean before printing?
	bytJSON, _ := json.MarshalIndent(iteration, "", "    ")
	fmt.Println(string(bytJSON) + ",")
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
