package distinct

import (
	"testing"
)

func TestDistinct(t *testing.T) {
	distinct := 300
	c := NewCounter[int](200, nil)
	for i := 0; i < 100000; i++ {
		n := i % distinct
		c.Add(n)
	}
	res := c.Estimate()
	t.Logf("estimate = %v", res)
	if res < 200 || res > 400 {
		t.Errorf("expected between [200, 400], got %v", res)
	}
}
