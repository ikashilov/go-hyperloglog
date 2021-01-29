package hyperloglog

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// get random int's as a []string slice
func getRandomData(size int) []string {
	var res []string

	rand.Seed(time.Now().Unix())
	for i := 0; i < size; i++ {
		res = append(res, strconv.Itoa(rand.Int()))
	}
	return res
}

// count the exact cardinality using hashmap
func estimateReal(stream []string) int {
	m := map[string]struct{}{}

	for _, s := range stream {
		m[s] = struct{}{}
	}
	return len(m)
}

func TestHLL(t *testing.T) {
	stream := getRandomData(100000)

	hll, _ := New(0.001)
	for _, x := range stream {
		hll.Add(x)
	}

	estimated := hll.Count()
	real := estimateReal(stream)

	fmt.Printf("classic estimate: %v\n", real)
	fmt.Printf("hyperloglog estimate: %v\n", estimated)
	fmt.Printf("percentage missed: %.2f\n", 100.-(float64(estimated)/float64(real))*100)
}

func BenchmarkHLL(b *testing.B) {
	stream := getRandomData(b.N)

	hll, _ := New(0.001)

	for n := 0; n < b.N; n++ {
		hll.Add(stream[n])
	}
}
