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
	"strings"
	"time"

	"github.com/griffindavis02/eth-bit-flip/config"
)

type errorData struct {
	PreviousValue interface{}
	PreviousByte  string
	IntBits       []int
	ErrorValue    interface{}
	ErrorByte     string
	DeltaValue    interface{}
	When          string
}

type iteration struct {
	Rate         float64
	IterationNum int
	ErrorData    errorData
}

// BitFlip will run the odds of flipping a bit within pbigNum based on error
// rate pdecRate. The iteration count will increment and both the new number
// and the iteration error data will be returned.
func BitFlip(pIFlipee interface{}) interface{} {
	cfg, err := config.ReadConfig()
	if err != nil {
		if strings.Contains(err.Error(), "error reading in config file") {
			return pIFlipee
		}
	}
	if !cfg.Start {
		return pIFlipee
	}
	if cfg.Restart {
		restart(&cfg)
	}

	// Check for out of bounds or end of error rate
	rand.Seed(time.Now().UnixNano())
	switch cfg.State.TestType {
	case "bit":
		if cfg.State.TestCounter >= cfg.State.Bits {
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
		if cfg.State.TestCounter == 0 && cfg.State.RateIndex == 0 {
			cfg.State.StartTime = time.Now().Unix()
		}
	}

	var iteration iteration
	switch pIFlipee.(type) {
	case []byte:
		iteration = flipBytes(pIFlipee.([]byte), &cfg)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}

		iteration.ErrorData.PreviousValue = iteration.ErrorData.PreviousValue.([]byte)
		iteration.ErrorData.ErrorValue = iteration.ErrorData.ErrorValue.([]byte)
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	case string:
		iteration = flipBytes([]byte(pIFlipee.(string)), &cfg)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}

		iteration.ErrorData.PreviousValue = string(iteration.ErrorData.PreviousValue.([]byte))
		iteration.ErrorData.ErrorValue = string(iteration.ErrorData.ErrorValue.([]byte))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	case int:
		switch binary.Size(pIFlipee.(int)) {
		case 32:
			bytInt := make([]byte, 4)
			binary.BigEndian.PutUint32(bytInt, uint32(pIFlipee.(int)))
			iteration = flipBytes(bytInt, &cfg)
			if iteration.ErrorData.ErrorValue == nil {
				return pIFlipee
			}

			iteration.ErrorData.PreviousValue = int(binary.BigEndian.Uint32(iteration.ErrorData.PreviousValue.([]byte)))
			iteration.ErrorData.ErrorValue = int(binary.BigEndian.Uint32(iteration.ErrorData.ErrorValue.([]byte)))
			iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
		default:
			bytInt := make([]byte, 8)
			binary.BigEndian.PutUint64(bytInt, uint64(pIFlipee.(int)))
			iteration = flipBytes(bytInt, &cfg)
			if iteration.ErrorData.ErrorValue == nil {
				return pIFlipee
			}

			iteration.ErrorData.PreviousValue = int(binary.BigEndian.Uint64(iteration.ErrorData.PreviousValue.([]byte)))
			iteration.ErrorData.ErrorValue = int(binary.BigEndian.Uint64(iteration.ErrorData.ErrorValue.([]byte)))
			iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
		}
	case int64:
		bytInt := make([]byte, 8)
		binary.BigEndian.PutUint64(bytInt, uint64(pIFlipee.(int64)))
		iteration = flipBytes(bytInt, &cfg)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}

		iteration.ErrorData.PreviousValue = int64(binary.BigEndian.Uint64(iteration.ErrorData.PreviousValue.([]byte)))
		iteration.ErrorData.ErrorValue = int64(binary.BigEndian.Uint64(iteration.ErrorData.ErrorValue.([]byte)))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	case int32:
		bytInt := make([]byte, 4)
		binary.BigEndian.PutUint32(bytInt, uint32(pIFlipee.(int32)))
		iteration = flipBytes(bytInt, &cfg)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}

		iteration.ErrorData.PreviousValue = int32(binary.BigEndian.Uint32(iteration.ErrorData.PreviousValue.([]byte)))
		iteration.ErrorData.ErrorValue = int32(binary.BigEndian.Uint32(iteration.ErrorData.ErrorValue.([]byte)))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	case uint:
		switch binary.Size(pIFlipee.(uint)) {
		case 32:
			bytInt := make([]byte, 4)
			binary.BigEndian.PutUint32(bytInt, uint32(pIFlipee.(uint)))
			iteration = flipBytes(bytInt, &cfg)
			if iteration.ErrorData.ErrorValue == nil {
				return pIFlipee
			}

			iteration.ErrorData.PreviousValue = uint(binary.BigEndian.Uint32(iteration.ErrorData.PreviousValue.([]byte)))
			iteration.ErrorData.ErrorValue = uint(binary.BigEndian.Uint32(iteration.ErrorData.ErrorValue.([]byte)))
			iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
		default:
			bytInt := make([]byte, 8)
			binary.BigEndian.PutUint64(bytInt, uint64(pIFlipee.(uint)))
			iteration = flipBytes(bytInt, &cfg)
			if iteration.ErrorData.ErrorValue == nil {
				return pIFlipee
			}

			iteration.ErrorData.PreviousValue = uint(binary.BigEndian.Uint64(iteration.ErrorData.PreviousValue.([]byte)))
			iteration.ErrorData.ErrorValue = uint(binary.BigEndian.Uint64(iteration.ErrorData.ErrorValue.([]byte)))
			iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
		}
	case uint32:
		bytInt := make([]byte, 4)
		binary.BigEndian.PutUint32(bytInt, uint32(pIFlipee.(uint32)))
		iteration = flipBytes(bytInt, &cfg)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}

		iteration.ErrorData.PreviousValue = uint32(binary.BigEndian.Uint32(iteration.ErrorData.PreviousValue.([]byte)))
		iteration.ErrorData.ErrorValue = uint32(binary.BigEndian.Uint32(iteration.ErrorData.ErrorValue.([]byte)))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	case uint64:
		bytInt := make([]byte, 8)
		binary.BigEndian.PutUint64(bytInt, uint64(pIFlipee.(uint64)))
		iteration = flipBytes(bytInt, &cfg)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}

		iteration.ErrorData.PreviousValue = uint64(binary.BigEndian.Uint64(iteration.ErrorData.PreviousValue.([]byte)))
		iteration.ErrorData.ErrorValue = uint64(binary.BigEndian.Uint64(iteration.ErrorData.ErrorValue.([]byte)))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	case *big.Int:
		iteration = flipBytes(pIFlipee.(*big.Int).Bytes(), &cfg)
		if iteration.ErrorData.ErrorValue == nil {
			return pIFlipee
		}
		iteration.ErrorData.PreviousValue = new(big.Int).SetBytes(iteration.ErrorData.PreviousValue.([]byte))
		iteration.ErrorData.ErrorValue = new(big.Int).SetBytes(iteration.ErrorData.ErrorValue.([]byte))
		iteration.ErrorData.DeltaValue = iteration.ErrorData.DeltaValue.(*big.Int)
	}

	printOut(iteration, &cfg)
	return iteration.ErrorData.ErrorValue
}

func flipBytes(pbytFlipee []byte, cfg *config.Config) iteration {
	decRate := cfg.State.ErrorRates[cfg.State.RateIndex]
	var arrBits []int
	var iter iteration

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
		fmt.Println("Decrate:", decRate)
		if cfg.State.TestType == "variable" {
			cfg.State.TestCounter++
		}
		// Build error data
		iter = iteration{
			cfg.State.ErrorRates[cfg.State.RateIndex],
			int(lngPrevCounter),
			errorData{
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

		cfg.WriteConfig()
	}

	return iter
}

func restart(cfg *config.Config) {
	cfg.State.TestCounter = 0
	cfg.State.RateIndex = 0
	cfg.Start = true
	cfg.Restart = false
}

func printOut(pIteration iteration, cfg *config.Config) {
	if pIteration.ErrorData.PreviousValue == pIteration.ErrorData.ErrorValue {
		return
	}
	// TODO: Look for logging boolean before printing?
	bytJSON, _ := json.MarshalIndent(pIteration, "", "    ")
	fmt.Println(string(bytJSON) + ",")
	if cfg.Server.Post {
		postAPI(cfg.Server.Host, pIteration)
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
