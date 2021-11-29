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
)

func (cfg *Config) manageWizard() {

	for {
		fmt.Println()
		fmt.Println("Which option would you like to change?")
		fmt.Println("For option 2, each number is per error rate")
		fmt.Printf("1 - Test Type (%s)\n", cfg.State.TestType)
		switch cfg.State.TestType {
		case "bit":
			fmt.Printf("2 - Number of bits to flip (%d)\n", cfg.State.Bits)
		case "variable":
			fmt.Printf("2 - Number of variables to change (%d)\n", cfg.State.VariablesChanged)
		case "time":
			fmt.Printf("2 - Amount of time to pass (%v)\n", cfg.State.Duration)
		}
		fmt.Printf("3 - Error Rates (%v)\n", cfg.State.ErrorRates)
		fmt.Print("4 - Server options for POST requests (")
		if cfg.Server.Post {
			fmt.Printf("posting to '%s')\n", cfg.Server.Host)
		} else {
			fmt.Println("not posting)")
		}
		fmt.Println("5 - Save & Quit")
		fmt.Println()

		choice := readInt()
		switch choice {
		case 1:
			cfg.promptTestType()
			continue
		case 2:
			cfg.promptTestCount()
			continue
		case 3:
			cfg.promptErrorRates()
			continue
		case 4:
			cfg.promptServer()
			continue
		}

		if choice == 5 {
			break
		}

		log.Printf("WARNING: choice must be within range 1-5. entered choice: \"%d\"", choice)
	}

	cfg.State.TestCounter = 0
	cfg.State.RateIndex = 0

	if err := cfg.WriteConfig(); err != nil {
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
}
