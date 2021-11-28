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

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func (cfg *Config) newWizard() error {
	cfg.promptTestType()
	cfg.promptTestCount()
	cfg.promptErrorRates()
	cfg.promptServer()

	cfg.Initialized = true
	if err := WriteConfig(*cfg); err != nil {
		log.Fatalf("ERROR: %v", err)
	} else {
		if cfgByt, marshErr := json.MarshalIndent(cfg, "", "\t"); marshErr == nil {
			fmt.Println(`
Congrats! You've configured your next soft error injection test!
Here is your configuration:
	`)
			fmt.Println(string(cfgByt))
		}
	}

	return nil
}

func (cfg *Config) promptTestType() {
	var testType string
	for {
		fmt.Println("What type of test will this be? Bit, variable, or time based?")
		testType = strings.ToLower(promptInput("Bit counts per bit flipped.\nVariable counts per variable enacted upon.\nTime counts... well based on time."))

		if strings.Compare(testType, "bit") == 0 ||
			strings.Compare(testType, "variable") == 0 ||
			strings.Compare(testType, "time") == 0 {
			cfg.State.TestType = testType
			break
		}
		log.Println("WARNING: ", fmt.Sprintf("Test type \"%s\" not accepted", testType))
	}
}

func (cfg *Config) promptTestCount() {
	switch cfg.State.TestType {
	case "bit":
		cfg.State.Bits = promptIntCB("How many bits per error rate?",
			func(input int) (int, error) {
				if input >= 0 {
					return input, nil
				} else {
					return -1, fmt.Errorf("received invalid iteration count of %d", input)
				}
			})
	case "variable":
		cfg.State.VariablesChanged = promptIntCB("How many variables changed per error rate?",
			func(input int) (int, error) {
				if input >= 0 {
					return input, nil
				} else {
					return -1, fmt.Errorf("received invalid variable count of %d", input)
				}
			})
	case "time":
		// FIXME: Need to update cfg.State.StartTime when geth is actually instantiated
		cfg.State.Duration = time.Duration(
			float64(promptIntCB("How long for the test in seconds?",
				func(input int) (int, error) {
					if input >= 0 {
						return input, nil
					} else {
						return -1, fmt.Errorf("received invalid duration of %d seconds", input)
					}
				})) * math.Pow(10, 9), // nanoseconds
		)
		cfg.State.StartTime = time.Now().Unix()
	}
}

func (cfg *Config) promptErrorRates() {
	for {
		strRates := promptStringCB("Please enter the error rates you would like to test for as a comma\nseparated list, no spaces.",
			func(input string) (string, error) {
				if strings.Contains(input, " ") {
					return "", fmt.Errorf("cannot use spaces")
				}
				if strings.Compare(input, "") == 0 {
					return "", fmt.Errorf("must list error rates")
				}
				return input, nil
			})
		var status int = 0
		var tmpArray []float64
		errRates := strings.Split(strRates, ",")
		for i, rate := range errRates {
			if decRate, err := strconv.ParseFloat(rate, 64); err == nil {
				if decRate <= 0 {
					log.Println("WARNING: error rate cannot be negative or 0")
					status = -1
					break
				} else if decRate > 1 {
					log.Println("WARNING: error rate cannot be more than 1")
					status = -1
					break
				} else {
					tmpArray = append(tmpArray, decRate)
				}
			} else {
				log.Printf("WARNING: invalid error rate \"%s\" in array", errRates[i])
				status = -1
				break
			}
		}
		if status == 0 {
			cfg.State.RateIndex = 0
			cfg.State.ErrorRates = tmpArray
			break
		}
	}
}

func (cfg *Config) promptServer() {
	for {
		post := promptInput("Will you be posting results to an API? [y/n]")
		post = strings.ToLower(post)
		if post == "y" || post == "yes" {
			cfg.Server.Post = true
			break
		} else if post == "n" || post == "no" {
			cfg.Server.Post = false
			break
		}
		log.Println("WARNING:", fmt.Sprintf("invalid response \"%s\". Require \"y/yes\" or \"n/no\"", post))
	}

	if cfg.Server.Post {
		cfg.Server.Host = promptStringCB("What is the API hostname?",
			func(input string) (string, error) {
				if strings.Compare(input, "") == 0 {
					return "", fmt.Errorf("must include a hostname if posting to API")
				}
				return input, nil
			})
	}
}
