package hyperloglog

import (
	"errors"
	"hash/fnv"
	"math"
)

const twoIn32Power float64 = 4294967296

// HyperLogLog is a cardinality estimator
type HyperLogLog struct {
	m         uint32
	k         float64
	kComp     int
	alpha     float64
	registers []uint8
}

// New creates a new HLL estimator with a given precision (error).
// The `good` value for error is (0.01 - 0.0001)
// Smaller values will case more memory usage
func New(err float64) (*HyperLogLog, error) {

	if err > 0.01 || err < 0.00001 {
		return nil, errors.New("invalied value of err. err should be in (0.1, 0.0001)")
	}

	var h HyperLogLog

	m := 1.04 / err

	h.k = math.Ceil(math.Log2(m * m))
	h.kComp = int(32 - h.k)
	h.m = uint32(math.Pow(2, h.k))
	h.alpha = getOptimalAlpha(h.m)
	h.registers = make([]uint8, h.m)

	return &h, nil
}

// Add adds a new []byte value to the HLL
func (h *HyperLogLog) Add(val string) {
	hash := createHash(val)
	r := 1
	for (hash&1) == 0 && r <= h.kComp {
		r++
		hash >>= 1
	}

	j := hash >> uint(h.kComp)
	if r > int(h.registers[j]) {
		h.registers[j] = uint8(r)
	}
}

// Count estimates the cardinality
func (h *HyperLogLog) Count() uint64 {
	var c float64
	for i := uint32(0); i < h.m; i++ {
		c += (1 / math.Pow(2, float64(h.registers[i])))
	}
	E := h.alpha * float64(h.m*h.m) / c

	// Correct E - just repeat the paper
	if E <= 5/2*float64(h.m) {
		var V float64
		for i := uint32(0); i < h.m; i++ {
			if h.registers[i] == 0 {
				V++
			}
		}
		if V > 0 {
			E = float64(h.m) * math.Log(float64(h.m)/V)
		}
	} else if E > 1/30*twoIn32Power {
		E = -twoIn32Power * math.Log(1-E/twoIn32Power)
	}
	return uint64(E)
}

// creates a 32-bit hash
func createHash(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	sum := h.Sum32()
	h.Reset()
	return sum
}

// returns the optimal `m` value acording to the paper
func getOptimalAlpha(m uint32) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1 + 1.079/float64(m))
	}
}
