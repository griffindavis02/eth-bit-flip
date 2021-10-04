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
	"time"

	mathEth "github.com/ethereum/go-ethereum/common/math"
)

type IBitFlip interface {
	Initalize(pstrTestType string, pITestCount interface{}, parrErrRates []float64, pOutput Output)
	BitFlip(pbigNum *big.Int) *big.Int
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
	mdurNanoSeconds time.Duration
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
		mdurNanoSeconds = time.Duration(pITestCount.(float64) * math.Pow(10, 9))
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
func (this *Output) BitFlip(pbigNum *big.Int) *big.Int {
	rand.Seed(time.Now().UnixNano())

	// Check for out of bounds or end of error rate
	switch mstrTestType {
	case "iteration":
		if mlngCounter >= mlngIterations {
			if mintRateIndex == len((*this).Data)-1 {
				return pbigNum
			}
			mintRateIndex++
			mlngCounter = 0
		}
	case "variable":
		if mlngCounter >= mlngVarsChanged {
			if mintRateIndex == len((*this).Data)-1 {
				return pbigNum
			}
			mintRateIndex++
			mlngCounter = 0
		}
	default:
		if time.Since(mtimStartTime) >= mdurNanoSeconds {
			if mintRateIndex == len((*this).Data)-1 {
				return pbigNum
			}
			mintRateIndex++
			mtimStartTime = time.Now()
		}
	}

	decRate := (*this).Data[mintRateIndex].Rate
	var arrBits []int

	// Store previous states
	lngPrevCounter := mlngCounter
	bigPrevNum, _ := new(big.Int).SetString(pbigNum.String(), 10)
	bigPrevNum = mathEth.U256(bigPrevNum)
	bytPrevNum := bigPrevNum.Bytes()
	bytNum := pbigNum.Bytes()

	// Run chance of flipping a bit in byte representation
	for i, byt := range bytNum {
		for j := 0; j < 8; j++ {
			if math.Floor(rand.Float64()/decRate) == math.Floor(rand.Float64()/decRate) {
				mlngCounter++
				arrBits = append(arrBits, (i*8)+j)
				bytNum[i] = byt ^ (1 << j)
			}
		}
	}

	// Ensure there was a change
	if !bytes.Equal(bytNum, bytPrevNum) {
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
		(*this).Data[mintRateIndex].FlipData = append((*this).Data[mintRateIndex].FlipData, iteration)
	}

	return pbigNum
}

func (this Output) MarshalIndent() string {
	byt, err := json.MarshalIndent(this, "", "\t")
	if err != nil {
		return "err"
	}
	return string(byt)
}
