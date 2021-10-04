package BitFlip

import (
	"encoding/hex"
	"math"
	"math/big"
	"math/rand"
	"time"

	mathEth "github.com/ethereum/go-ethereum/common/math"
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

// BitFlip will run the odds of flipping a bit within pbigNum based on error
// rate pdecRate. The iteration count will increment and both the new number
// and the iteration error data will be returned.
func BitFlip(pbigNum *big.Int, pdecRate float64, plngFlipCount int) (*big.Int, Iteration) {
	rand.Seed(time.Now().UnixNano())

	var arrBits []int

	bigPrevNum, _ := new(big.Int).SetString(pbigNum.String(), 10)
	bigPrevNum = mathEth.U256(bigPrevNum)
	bytPrevNum := bigPrevNum.Bytes()
	bytNum := pbigNum.Bytes()

	for i, byt := range bytNum {
		for j := 0; j < 8; j++ {
			if math.Floor(rand.Float64()/pdecRate) == math.Floor(rand.Float64()/pdecRate) {
				plngFlipCount++
				arrBits = append(arrBits, (i*8)+j)
				bytNum[i] = byt ^ (1 << j)
			}
		}
	}

	pbigNum.SetBytes(bytNum)

	iteration := Iteration{
		int(plngFlipCount),
		ErrorData{
			bigPrevNum,
			hex.EncodeToString(bytPrevNum),
			arrBits,
			pbigNum,
			hex.EncodeToString(bytNum),
			big.NewInt(0).Sub(pbigNum, bigPrevNum),
			time.Now().Format("01-02-2006-15:04:05.000000000"),
		},
	}

	return pbigNum, iteration
}
