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
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
)

type State struct {
	TestType         string        `json:"test_type"`
	TestCounter      int           `json:"test_counter"`
	Iterations       int           `json:"iterations"`
	VariablesChanged int           `json:"variables_changed"`
	Duration         time.Duration `json:"duration"`
	StartTime        int64         `json:"start_time"`
	RateIndex        int           `json:"rate_index"`
	ErrorRates       string        `json:"error_rates"`
}

type Server struct {
	Post bool   `json:"post"`
	Host string `json:"host"`
}

type Config struct {
	Initialized bool   `json:"initialized"`
	State       State  `json:"state_variables"`
	Server      Server `json:"server"`
}

var (
	// flipCommand = cli.Command {
	// 	Action: flipInit,
	// 	Name: "flip",
	// 	Usage: "Set up a soft error test environment for geth",
	// 	Flags: []cli.Flag{enableFlag,disableFlag,resetFlag,},
	// 	Category: "SOFT ERROR INJECTION",
	// 	Description: `
	// 	The flip command allows one to configure the parameters by which to simulate
	// 	soft errors in the EVM.

	// 	usage: geth flip [--enable] [--disable] [--reset]
	// 	Lack of flag will begin the configuration wizard.
	// 	`,
	// }

	reader = bufio.NewReader(os.Stdin)
	Path   = filepath.FromSlash("./flipconfig.json")

	DefaultConfig = Config{
		Initialized: false,
		State: State{
			TestType:         "iteration",
			TestCounter:      0,
			Iterations:       0,
			VariablesChanged: 0,
			Duration:         time.Duration(0),
			StartTime:        time.Now().Unix(),
			RateIndex:        0,
			ErrorRates:       "0.1",
		},
		Server: Server{
			Post: false,
			Host: "http://localhost:5000",
		},
	}
)

func RunConfig() {
	app := cli.NewApp()
	app.Name = "flipcfg"
	app.Usage = "Set up a soft error test environment for go-ethereum"
	app.Flags = []cli.Flag{
		// TODO: Add flags from utils
	}
	app.Action = flipWizard
	app.Run(os.Args)
}

func flipWizard(ctx *cli.Context) {
	var cfg = DefaultConfig

	fmt.Println(" ----------------------------------------------------------- ")
	fmt.Println("| This is flipcfg, your soft error simulation config tool   |")
	fmt.Println("|                                                           |")
	fmt.Println("| This will allow you to configure the test parameters for  |")
	fmt.Println("| simulating soft errors in the EVM. The parameters decide  |")
	fmt.Println("| the ranomness of bit flipping based on data you wish to   |")
	fmt.Println("| collect.                                                  |")
	fmt.Println("|                                                           |")
	fmt.Println("| This tool and all other items within the eth-bit-flip     |")
	fmt.Println("| repository are modifiable under the terms of GPLv3.0.     |")
	fmt.Println(" ----------------------------------------------------------- ")
	fmt.Println()

	var testType string
	for {
		fmt.Println("What type of test will this be? Iteration, variable, or time based?")
		testType = strings.ToLower(promptInput("Iteration counts per bit flipped.\nVariable counts per variable enacted upon.\nTime counts... well based on time."))

		if strings.Compare(testType, "iteration") == 0 ||
			strings.Compare(testType, "variable") == 0 ||
			strings.Compare(testType, "time") == 0 {
			cfg.State.TestType = testType
			break
		}
		log.Println("WARNING: ", fmt.Sprintf("Test type \"%s\" not accepted", testType))
	}

	switch testType {
	case "iteration":
		cfg.State.Iterations = promptIntCB("How many iterations per error rate?",
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
				}))  * math.Pow(10, 9), // nanoseconds
		)
		cfg.State.StartTime = time.Now().Unix()
	}

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
		if _, err := AtoF64Arr(strRates, "0"); err == nil {
			cfg.State.ErrorRates = strRates
			break
		} else {
			log.Println("WARNING:", err)
		}
	}

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

	cfg.Initialized = true
	if err := WriteConfig(cfg); err != nil {
		log.Println("ERROR:", err)
	} else {
		if cfgByt, marshErr := json.MarshalIndent(cfg, "", "\t"); marshErr == nil {
			fmt.Println(`Congrats! You've configured your next soft error injection test!
	Here is your configuration:
	`)
			fmt.Println(string(cfgByt))
		}
	}
}

func promptInput(prompt string) string {
	prompt = strings.TrimSpace(prompt)
	if strings.Compare(prompt, "") == 0 {
		return ""
	}
	fmt.Println(prompt)
	return readString()
}

func readString() string {
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		fmt.Println()
		return text
	}
}

func readInt() int {
	text := readString()
	if i, err := strconv.Atoi(text); err == nil {
		return i
	}
	return -1
}

func promptStringCB(prompt string, callback func(input string) (string, error)) string {
	prompt = strings.TrimSpace(prompt)
	if strings.Compare(prompt, "") == 0 {
		return ""
	}
	for {
		fmt.Println(prompt)
		if input, err := callback(readString()); err == nil {
			return input
		} else {
			log.Println("WARNING:", err)
		}
	}
}

func promptIntCB(prompt string, callback func(input int) (int, error)) int {
	prompt = strings.TrimSpace(prompt)
	if strings.Compare(prompt, "") == 0 {
		return -1
	}
	for {
		fmt.Println(prompt)
		if input, err := callback(readInt()); err == nil {
			return input
		} else {
			log.Println("WARNING:", err)
		}
	}
}

func WriteConfig(cfg Config) error {
	bytCfg, err := json.MarshalIndent(cfg, "", "\t")
	if err == nil {
		if fErr := os.WriteFile(Path, bytCfg, 0644); fErr != nil {
			return fmt.Errorf("error writing to file \"%s\"", Path)
		}
		fmt.Println("\nYou've configured your next soft error injection test! Here is your configuration:")
		fmt.Println(string(bytCfg))
		return nil
	}
	return fmt.Errorf("error marshaling config")
}

func ReadConfig() (Config, error) {
	if bytes, fErr := os.ReadFile(Path); fErr == nil {
		var cfg Config
		if err := json.Unmarshal(bytes, &cfg); err != nil {
			return Config{}, fmt.Errorf("error unmarshaling file data into config")
		}
		return cfg, nil
	}

	return Config{}, fmt.Errorf("error reading in config file from %s", Path)
}

// params[0] would be the array and params[1] a number to avoid
func AtoF64Arr(params ...string) ([]float64, error) {
	pLen := len(params)
	var avoid float64
	if pLen > 2 {
		return nil, fmt.Errorf("you can enter just a string array, or one number to avoid as a string ('0')")
	}
	if pLen == 2 {
		if ch, err := strconv.ParseFloat(params[1], 64); err != nil {
			return nil, fmt.Errorf("error parsing character to avoid :\"%s\"", params[1])
		} else {
			avoid = ch
		}
	}
	var tmpArray []float64
	errRates := strings.Split(params[0], ",")
	for i, rate := range errRates {
		if decRate, err := strconv.ParseFloat(rate, 64); err == nil {
			if pLen == 2 && decRate == avoid {
				return nil, fmt.Errorf("entry \"%f\" matches number to avoid \"%f\"", decRate, avoid)
			}
			tmpArray = append(tmpArray, decRate)
		} else {
			return nil, fmt.Errorf("invalid entry \"%s\" in array.\nyou will need to enter a new array", errRates[i])
		}
	}
	return tmpArray, nil
}
