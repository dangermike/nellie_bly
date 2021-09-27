package die

import (
	cr "crypto/rand"
	"math/big"
	"math/rand"
)

// Die is a n-sided thing
type Die struct {
	sides int
	rnd   rand.Rand
}

// New makes a new die with the specified number of sides
func New(sides int) *Die {
	var rr = cr.Reader
	seed, _ := cr.Int(rr, big.NewInt((1<<63)-1))
	return FromSeed(sides, seed.Int64())
}

func FromSeed(sides int, seed int64) *Die {
	return &Die{
		sides: sides,
		rnd:   *rand.New(rand.NewSource(seed)),
	}
}

func (die *Die) Roll() int {
	return die.rnd.Intn(die.sides) + 1
}
