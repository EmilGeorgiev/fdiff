package rollinghash

import "fmt"

const (
	// 2^13 irreducible polynomial
	defaultModulus      = 8191
	multiplier          = 257
	numberOfBitsPerByte = 8
)

type rabinFingerprintHash struct {
	value                    uint64
	polynomialsOverGF2Values []uint64
}

// NewRabinFingerprint created a new Rabin fingerprint hash.
func NewRabinFingerprint(bytes []byte) Hash {
	rfh := &rabinFingerprintHash{}
	var h uint64
	degreeOfPolynomial := len(bytes) * numberOfBitsPerByte
	for _, b := range bytes {
		var polynomialOverGF2 uint64
		for j := 7; j >= 0; j-- {
			// polynomials over GF(2)
			mask := int32(1 << uint(j))
			degreeOfPolynomial -= 1
			term := pow(multiplier, degreeOfPolynomial)
			bit := int32(b) & mask
			if bit == 0 {
				continue
			}
			polynomialOverGF2 += term
		}
		rfh.polynomialsOverGF2Values = append(rfh.polynomialsOverGF2Values, polynomialOverGF2)
		h += polynomialOverGF2
	}
	rfh.value = h % defaultModulus
	return rfh
}

// Value return the value of the hash.
func (rfh *rabinFingerprintHash) Value() uint64 {
	return rfh.value
}

// Next calculate the hash of the next rolling window. The window is shifted with one byte.
func (rfh *rabinFingerprintHash) Next(b byte) uint64 {
	var polynomialOverGF2 uint64
	if rfh.value == 5643 {
		for i, ff := range rfh.polynomialsOverGF2Values {
			fmt.Printf("Index: %d value: %d\n", i, ff)
		}
	}
	for j := 7; j >= 0; j-- {
		mask := int32(1 << uint(j))
		term := pow(multiplier, j)
		bit := int32(b) & mask
		if bit == 0 {
			continue
		}
		polynomialOverGF2 += term
	}
	for i, p := range rfh.polynomialsOverGF2Values[1:] {
		rfh.polynomialsOverGF2Values[i+1] = (p * pow(multiplier, 8)) % defaultModulus
	}
	rfh.value = (((rfh.value+defaultModulus-rfh.polynomialsOverGF2Values[0]%defaultModulus)*pow(multiplier, 8))%defaultModulus + polynomialOverGF2) % defaultModulus
	rfh.polynomialsOverGF2Values = append(rfh.polynomialsOverGF2Values[1:], polynomialOverGF2)
	return rfh.value
}

func pow(a uint64, b int) uint64 {
	if b == 0 {
		return 1
	}
	n := uint64(1)
	for i := 0; i < b; i++ {
		// we use module because we don't want t\o work with big integers. Also, it is much faster.
		n = (n * a) % defaultModulus
	}
	return n
}

func Pow(a uint64, b int) uint64 {
	if b == 0 {
		return 1
	}
	n := uint64(1)
	for i := 0; i < b; i++ {
		// we use module because we don't want t\o work with big integers. Also, it is much faster.
		n = (n * a) % 113
	}
	return n
}
