package rabin

const (
	defaultModulus = 8191
	multiplier     = 257
)

type FingerprintHash interface {
	Value() uint64
}

type rabinFingerprintHash struct {
	value           uint64
	characterValues []uint64
}

func (r rabinFingerprintHash) Value() uint64 {
	return r.value
}

func Pow(a uint64, b int) uint64 {
	if b == 0 {
		return 1
	}
	n := uint64(1)
	for i := 0; i < b; i++ {
		n = (n * a) % defaultModulus
	}
	return n
}

func getNumberOfBytes(i int32) int {
	switch true {
	case i <= 127:
		return 1
	case i <= 2047:
		return 2
	case i <= 65535:
		return 3
	case i <= 1114111:
		return 4
	}

	return 0
}

func NewFingerprintHash(str string) FingerprintHash {
	var rfh rabinFingerprintHash
	var h uint64
	degreeOfPolynomial := len(str) * 8
	for _, b := range str {
		var polynomialOverGF2 uint64
		numberOfBytes := getNumberOfBytes(b)
		for j := 8*numberOfBytes - 1; j >= 0; j-- {
			// polynomials over GF(2)
			mask := int32(1 << uint(j))
			degreeOfPolynomial -= 1
			p := Pow(multiplier, degreeOfPolynomial)
			bit := b & mask
			if bit == 0 {
				continue
			}
			polynomialOverGF2 += p
		}
		rfh.characterValues = append(rfh.characterValues, polynomialOverGF2)
		h += polynomialOverGF2
	}
	rfh.value = h % defaultModulus
	return rfh
}
