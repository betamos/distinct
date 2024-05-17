package distinct

import (
	"math"
	"math/rand/v2"
)

// Counter of distinct elements, using https://arxiv.org/pdf/2301.10191
type Counter[T comparable] struct {
	thresh int
	rand   rand.Source

	// Proxy for p. Since p = 1, 1/2, 1/4... we instead store the generation count
	// and make use of the fact the
	gen uint8
	els map[T]struct{}
}

// Create a new counter with O(threshold) memory complexity.
// If src is nil, a source with a random seed will be used.
func NewCounter[T comparable](threshold int, src rand.Source) *Counter[T] {
	if src == nil {
		src = rand.NewPCG(rand.Uint64(), rand.Uint64())
	}
	return &Counter[T]{
		threshold, rand.New(src), 0, make(map[T]struct{}),
	}
}

// Returns the current estimate
func (c *Counter[T]) Estimate() uint64 {
	// From paper: N / p => N * 1/p => N * (1 << gen)
	// Example: 25 / (1/8) => 25 * 8 => 25 * (1 << 4)
	return uint64(len(c.els)) * (1 << c.gen)
}

// Add an element to the counter
func (c *Counter[T]) Add(el T) {
	delete(c.els, el)
	if c.rand.Uint64() <= math.MaxUint64>>c.gen { // probability of 1, 1/2, 1/4 etc
		c.els[el] = struct{}{}
	}
	if len(c.els) < c.thresh {
		return
	}
	for el := range c.els {
		if c.rand.Uint64() < math.MaxUint64/2 { // coin flip: probability 1/2
			delete(c.els, el)
		}
	}
	c.gen++
	if len(c.els) == c.thresh {
		panic("distinct: N = thresh")
	}
	if c.gen == 64 {
		panic("distinct: too many generations")
	}
}

// Returns a suitable threshold value given (see Chernoff Bounds):
//
//   - epsilon: relative error of the estimate, lower is more accurate
//   - delta: confidence level, lower is more accurate
//   - m: total expected elements in the stream
func Threshold(epsilon, delta float64, m int) int {
	return int(math.Ceil((12 / (epsilon * epsilon)) * math.Log(8*float64(m)/delta)))
}
