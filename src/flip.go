// Copyright 2021 The eth-bit-flip Authors
// This file is part of the eth-bit-flip library.
//
// The eth-bit-flip libary is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// The eth-bit-flip libary is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
// details.
//
// You should have received a copy of the GNU General Public License along with
// the eth-bit-flip library. If not, see <https://www.gnu.org/licenses/>.

package BitFlip

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
	"path/filepath"
	"time"

	mathEth "github.com/ethereum/go-ethereum/common/math"
	"github.com/spf13/viper"
)

// TODO: Populate with remaining functions
type IBitFlip interface {
	Initalize(pstrTestType string, pITestCount interface{}, parrErrRates []float64, pOutput Output)
	BitFlip(pbigNum *big.Int) *big.Int
}

type Config struct {
	Initialized bool `json:"initialized"`

	State struct {
		TestType string `json:"test_type"`
		TestCounter int `json:"test_counter"`
		Iterations int `json:"iterations"`
		VariablesChanged int `json:"variables_changed"`
		Duration time.Duration `json:"duration"`
		StartTime time.Time `json:"start_time"`
		RateIndex int `json:"rate_index"`
		ErrorRates []int `json:"error_rates"`
	} `json:"state_variables"`

	Server struct {
		Post bool `json:"post"`
		Host string `json:"host"`
		Port int `json:"port"`
		Endpoint string `json:"endpoint"`
	} `json:"server"`
}

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
	mstrTestType    string
	mlngCounter     int
	mlngIterations  int
	mlngVarsChanged int
	mdurDurationNs  time.Duration
	mtimStartTime   time.Time
	mintRateIndex   int = 0
)

// Set up the testing environment with the test type, number of
// changes/iterations or duration in seconds, and error rates. This is
// PER error rate. i.e. 5 minutes and ten error rates will be 50 minutes.
// Test types:
// 'iteration' - increments for each bit flipped
// 'variable' - increments for each variable, regardless of bits flipped
// 'time' - checks against passage of time since started
func Initalize(pstrTestType string, pITestCount interface{}, parrErrRates []float64, pOutput *Output) {
	mstrTestType = pstrTestType
	switch mstrTestType {
	case "iteration":
		mlngIterations = pITestCount.(int)
	case "variable":
		mlngVarsChanged = pITestCount.(int)
	case "time":
		mtimStartTime = time.Now()
		mdurDurationNs = time.Duration(pITestCount.(float64) * math.Pow(10, 9))
	default:
		log.Fatal("Must use a valid test type: 'iteration', 'variable', 'time'")
	}
	var flipData []Iteration
	for _, errRate := range parrErrRates {
		Rate := ErrorRate{errRate, flipData}
		(*pOutput).Data = append(pOutput.Data, Rate)
	}
}

// BitFlip will run the odds of flipping a bit within pbigNum based on error
// rate pdecRate. The iteration count will increment and both the new number
// and the iteration error data will be returned.
func (jsonOut *Output) BitFlip(pbigNum *big.Int) *big.Int {
	rand.Seed(time.Now().UnixNano())

	// Check for out of bounds or end of error rate
	switch mstrTestType {
	case "iteration":
		if mlngCounter >= mlngIterations {
			if mintRateIndex == len((*jsonOut).Data)-1 {
				return pbigNum
			}
			mintRateIndex++
			mlngCounter = 0
		}
	case "variable":
		if mlngCounter >= mlngVarsChanged {
			if mintRateIndex == len((*jsonOut).Data)-1 {
				return pbigNum
			}
			mintRateIndex++
			mlngCounter = 0
		}
	default:
		if time.Since(mtimStartTime) >= mdurDurationNs {
			if mintRateIndex == len((*jsonOut).Data)-1 {
				return pbigNum
			}
			mintRateIndex++
			mtimStartTime = time.Now()
		}
	}

	decRate := (*jsonOut).Data[mintRateIndex].Rate
	var arrBits []int

	// Store previous states
	lngPrevCounter := mlngCounter
	bigPrevNum, _ := new(big.Int).SetString(pbigNum.String(), 10)
	bigPrevNum = mathEth.U256(bigPrevNum)
	bytPrevNum := bigPrevNum.Bytes()
	bytNum := pbigNum.Bytes()
	intLastByte := len(bytNum) - 1

	// Run chance of flipping a bit in byte representation
	for i := range bytNum {
		for j := 0; j < 8; j++ {
			if math.Floor(rand.Float64()/decRate) == math.Floor(rand.Float64()/decRate) {
				if mstrTestType == "iteration" {
					mlngCounter++
				}
				arrBits = append(arrBits, (i*8)+j)
				bytNum[intLastByte-i] ^= (1 << j)
			}
		}
	}

	// Ensure there was a change
	if !bytes.Equal(bytNum, bytPrevNum) {
		if mstrTestType == "variable" {
			mlngCounter++
		}
		// Recreate number from byte code
		pbigNum.SetBytes(bytNum)
		// Build error data
		iteration := Iteration{
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
		(*jsonOut).Data[mintRateIndex].FlipData = append((*jsonOut).Data[mintRateIndex].FlipData, iteration)
		(*jsonOut).PostAPI("http://localhost:5000/express")
	}

	return pbigNum
}

func (jsonOut Output) Marshal() string {
	byt, err := json.Marshal(jsonOut)
	if err != nil {
		return "err"
	}
	return string(byt)
}

func (jsonOut Output) MarshalIndent() string {
	byt, err := json.MarshalIndent(jsonOut, "", "\t")
	if err != nil {
		return "err"
	}
	return string(byt)
}

func (jsonOut Output) PostAPI(url string) int {
	client := http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	query := req.URL.Query()
	query.Add("params", jsonOut.Marshal())
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return res.StatusCode
}

// TODO: Implement or Change to encoding/json
func getState(cfg *Config, cfgPath string) {
	path := filepath.Dir(cfgPath)
	fileType := filepath.Ext(cfgPath)
	fileName := filepath.Base(cfgPath)

	viper.SetConfigName(fileName[0:len(fileName)-len(fileType)])
	viper.SetConfigType(fileType[1:])
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("The file could not be found... try again?")
			//TODO: Add a recursive call for filepath?
		} else {
			log.Fatal(err)
		}
	}
	viper.Unmarshal(cfg)
}
