package rabin

const (
	// 2^13 irreducible polynomial
	defaultModulus      = 8191
	multiplier          = 257
	numberOfBitsPerByte = 8
)

type FingerprintHash interface {
	Value() uint64
	Next(byte) uint64
}

type rabinFingerprintHash struct {
	value                    uint64
	polynomialsOverGF2Values []uint64
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

func NewFingerprintHash(bytes []byte) FingerprintHash {
	var rfh rabinFingerprintHash
	var h uint64
	degreeOfPolynomial := len(bytes) * numberOfBitsPerByte
	for _, b := range bytes {
		var polynomialOverGF2 uint64
		for j := 7; j >= 0; j-- {
			// polynomials over GF(2)
			mask := int32(1 << uint(j))
			degreeOfPolynomial -= 1
			term := Pow(multiplier, degreeOfPolynomial)
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

func (rfh rabinFingerprintHash) Next(b byte) uint64 {
	var polynomialOverGF2 uint64
	for j := 7; j >= 0; j-- {
		mask := int32(1 << uint(j))
		term := Pow(multiplier, j)
		bit := int32(b) & mask
		if bit == 0 {
			continue
		}
		polynomialOverGF2 += term
	}
	rfh.value = (((rfh.value+defaultModulus-rfh.polynomialsOverGF2Values[0]%defaultModulus)*Pow(multiplier, 8))%defaultModulus + polynomialOverGF2) % defaultModulus
	rfh.polynomialsOverGF2Values = append(rfh.polynomialsOverGF2Values[1:], polynomialOverGF2)
	return rfh.value
}
