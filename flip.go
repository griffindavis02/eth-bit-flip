package BitFlip

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	mathEth "github.com/ethereum/go-ethereum/common/math"
)

type ErrorData struct {
	PreviousValue int64
	IntBit        int16
	ErrorValue    int64
	DeltaValue    int64
	When          string
}

type Iteration struct {
	IterationNum int64
	Data         ErrorData
}

type ErrorRate struct {
	Rate     float64
	FlipData []Iteration
}

type Output struct {
	Data []ErrorRate
}

// BitFlip will run the odds of flipping a bit within plngNum based on error
// rate pdecRate. The iteration count will increment and both the new number
// and the iteration error data will be returned.
func BitFlip(plngNum int64, pintWordSize int16, pdecRate float64, plngFlipCount int64) (int64, []Iteration) {
	rand.Seed(time.Now().UnixNano())

	var arrErr []Iteration
	for i := int16(0); i < pintWordSize*8; i++ {
		if math.Floor(rand.Float64()/pdecRate) == math.Floor(rand.Float64()/pdecRate) {
			plngFlipCount++

			lngPrevNum := plngNum

			plngNum ^= (1 << i)

			errData := ErrorData{int64(lngPrevNum), int16(i), int64(plngNum), int64(plngNum - lngPrevNum), time.Now().Format("01-02-2006 15:04:06.000000000")}
			iterData := Iteration{int64(plngFlipCount), ErrorData(errData)}
			arrErr = append(arrErr, iterData)
		}
	}
	// fmt.Println(arrErr)
	return plngNum, arrErr
}

func Testing() {
	test := mathEth.U256(big.NewInt(mathEth.MaxBig256.Int64()))
	bTest := test.Bytes()

	fmt.Println(test)
	fmt.Println(*test)
	fmt.Println(bTest)

	for i, byt := range bTest {
		bit := int64(math.Floor(rand.Float64() * 8))
		bTest[i] = byt ^ (1 << bit)
	}

	test.SetBytes(bTest)

	fmt.Println(test)
	fmt.Println(*test)
	fmt.Println(test.Bytes())
}
