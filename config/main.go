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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
)

type state struct {
	TestType         string        `json:"test_type"`
	TestCounter      int           `json:"test_counter"`
	Bits             int           `json:"bits"`
	VariablesChanged int           `json:"variables_changed"`
	Duration         time.Duration `json:"duration"`
	StartTime        int64         `json:"start_time"`
	RateIndex        int           `json:"rate_index"`
	ErrorRates       []float64     `json:"error_rates"`
}

type server struct {
	Post bool   `json:"post"`
	Host string `json:"host"`
}

type Config struct {
	Initialized bool   `json:"initialized"`
	Start       bool   `json:"start"`
	Restart     bool   `json:"restart"`
	State       state  `json:"state_variables"`
	Server      server `json:"server"`
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
	file   = filepath.Join(os.Getenv("HOME"), ".flipconfig", "flipconfig.json")

	DefaultConfig = Config{
		Initialized: false,
		Start:       false,
		Restart:     false,
		State: state{
			TestType:         "bit",
			TestCounter:      0,
			Bits:             0,
			VariablesChanged: 0,
			Duration:         time.Duration(0),
			StartTime:        time.Now().Unix(),
			RateIndex:        0,
			ErrorRates:       []float64{0.1},
		},
		Server: server{
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
	fmt.Println("| This is flipconfig, the soft-error simulation config tool |")
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

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		cfg.newWizard()
	} else {
		for {
			fmt.Println("What would you like to do?")
			fmt.Println("1 - Manage existing configuration")
			fmt.Println("2 - Create a new configuration")
			choice := readInt()

			if choice == 1 {
				cfg, err := ReadConfig()
				if err != nil {
					log.Fatalf("ERROR: Corrupted configuration\n%v", err)
				}
				cfg.manageWizard()
				break
			} else if choice == 2 {
				cfg.newWizard()
				break
			} else {
				log.Println("WARNING: must enter '1' or '2'")
			}
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

func (cfg *Config) WriteConfig() error {
	bytCfg, err := json.MarshalIndent(cfg, "", "\t")
	if err == nil {
		if dErr := os.MkdirAll(filepath.Dir(file), os.ModePerm); dErr != nil {
			return fmt.Errorf("error creating directory \"%s\"", filepath.Dir(file))
		}
		if fErr := os.WriteFile(file, bytCfg, 0644); fErr != nil {
			return fmt.Errorf("error writing to file \"%s\"", file)
		}
		return nil
	}
	return fmt.Errorf("error marshaling config")
}

func ReadConfig() (Config, error) {
	if bytes, fErr := os.ReadFile(file); fErr == nil {
		var cfg Config
		if err := json.Unmarshal(bytes, &cfg); err != nil {
			return Config{}, fmt.Errorf("error unmarshaling file data into config")
		}
		return cfg, nil
	}

	return Config{}, fmt.Errorf("error reading in config file from %s", file)
}
