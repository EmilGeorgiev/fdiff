// Package rollinghash implements Rabin hashing (fingerprinting).
//
// The Rabin fingerprinting scheme is a method for implementing
// fingerprints using polynomials over a finite field. It was
// proposed by Michael O. Rabin.
//
// The schema of the Rabin hash is defined by a polynomial over
// GF(2):
//
//	p(x) = ... + p₂x² + p₁x + p₀   where pₙ ∈ GF(2)
//
// where p₁,p₂,..  are coefficients that represents the bits of
// the message in left-to-right most-significant-bit-first order
// (the values of p₁,p₂,.. can be '0' or '1']). 'x' is multiplier
// that is powered with the number of the bit in the message.
// All coefficients are multiplied with this multiplier, because
// this avoids collisions. For example the messages "abcd" and
// "cdab" will have different fingerprint hash.
//
// After the above polynomial is calculated we pick a random
// irreducible polynomial f(x) of degree k over GF(2), and we
// define the fingerprint of the message m to be the remainder
// r(x) after division of p(x) mod f(x) over GF(2).
//
// The message to be hashed is likewise interpreted as a polynomial
// over GF(2), where the coefficients are the bits of the message
// in left-to-right most-significant-bit-first order. Given a
// message polynomial m(x) and a hashing polynomial p(x), the
// Rabin hash is simply the coefficients of m(x) mod p(x).
//
// Rabin hashing efficiently compute a "rolling hash" of data,
// where the hash value is calculated base on the bytes in the
// window of the data. This makes the algorithm ideal for splitting
// the data in chunks with boundaries that are robust when bytes
// ar shifter to the left of right.
package rollinghash

const (
	// defaultModules is irreducible polynomial over GF(2).
	// 8191 = 2^13.
	modulus = 8191

	// multiplier is used to multiply the coefficients of the polynomial. Also,
	// this number is powered with the number of the bit of the coefficient.
	multiplier = 127

	// numberOfBitsPerByte is the number of the bits in one byte.
	numberOfBitsPerByte = 8
)

// polyGF2 represent a polynomial over GF(2). it contains the value of the polynomial.
type polyGF2 struct {
	value uint64
}

// table contains all polynomials in the window. The window is a number of bytes that is rolling.
// For example if the window size is 4 for thr string "abcdefgh" we will have:
//
//	|abcd|efgh - first the window has a value "abcd"
//	a|bcde|fgh - than we're rolling th window with one byte and the next value of the window will be "bcde"
//	ab|cdef|gh - than the window is "cdef"
//	...
//
// Every one byte of the window is represented as a polynomial over GF(2). The table
// contains all these polynomials at given time. For example table will contain:
//
//	|abcd|efgh - representation of the polynomials of 'a', 'b' 'c' and 'd'
//	a|bcde|fgh - representation of the polynomials of 'b' 'c', 'd' and 'e'
//	ab|cdef|gh - representation of the polynomials of 'c', 'd', 'e' and 'f'
//	...
//
// On every step, when the window is shifted/rolling, the first polynomial (the older one)
// is removed and the new one is added to the table.
type table struct {
	polynomials []polyGF2
}

// add a new polynomial to the table.
func (t *table) add(gf2 polyGF2) {
	t.polynomials = append(t.polynomials, gf2)
}

// removeFirstPolyGF2 remove first polynomial of the table.
func (t *table) removeFirstPolyGF2() polyGF2 {
	p := t.polynomials[0]
	t.polynomials = t.polynomials[1:]
	return p
}

type rabinFingerprintHash struct {
	value uint64
	table *table
}

// NewRabinFingerprint created a new Rabin fingerprint hash.
func NewRabinFingerprint(bytes []byte) Hash {
	rfh := &rabinFingerprintHash{
		table: &table{},
	}
	var h uint64
	degreeOfPolynomial := len(bytes) * numberOfBitsPerByte
	for _, b := range bytes {
		var p polyGF2
		for j := 7; j >= 0; j-- {
			// polynomials over GF(2)
			mask := int32(1 << uint(j))
			degreeOfPolynomial -= 1
			term := pow(multiplier, degreeOfPolynomial)
			bit := int32(b) & mask
			if bit == 0 {
				continue
			}
			p.value += term
		}
		rfh.table.add(p)
		h += p.value
	}
	rfh.value = h % modulus
	return rfh
}

// Value return the value of the hash.
func (rfh *rabinFingerprintHash) Value() uint64 {
	return rfh.value
}

// Next calculate the hash of the next rolling window. The window is shifted with one byte.
func (rfh *rabinFingerprintHash) Next(b byte) uint64 {
	var polynomialOverGF2 polyGF2
	for j := 7; j >= 0; j-- {
		mask := int32(1 << uint(j))
		term := pow(multiplier, j)
		bit := int32(b) & mask
		if bit == 0 {
			continue
		}
		polynomialOverGF2.value += term
	}
	// remove the oldest element of the table because it
	// will remain outside the window after this shifting
	firstPolyGF2 := rfh.table.removeFirstPolyGF2()

	// We must power all polynomials with multiplier*8 because all polynomial will be moved with one
	// byte to the left. This operation is needed if we want to calculate properly the hash next time
	// when this method is called.
	for i, p := range rfh.table.polynomials {
		rfh.table.polynomials[i].value = (p.value * pow(multiplier, 8)) % modulus
	}

	// To get the value of the new hash we need to do the these steps:

	// 1. Remove the value the oldest polynomial, which will remain outside the window after shifting,
	// from the current hash. NOTE we use + modulus because it is possible the result of this operation
	// to be negative number. It's OK to take a MOD of a negative number. Negative numbers map to
	// positive ones in modular arithmetic, but I prefer to convert it to positive number
	v := rfh.value + modulus - firstPolyGF2.value%modulus

	// 2. Multiply the value with multiplier powered with 8, because all polynomials over GF(2) are
	// moved with one byte to the left and one byte has 8 bits, so the degree of each multiplier
	// will be greater with 8. For example:
	//
	//  Let's say that thw window size is 2 bytes and the current window contains this binary message
	//  0100000101100011, then the polynomial over GF(2) for that binary message will be:
	//      P(x) = (x^14 + x^8) + (x^6 + x^5 + x^1 + x^0)
	//
	//  When the window is rolling/shifted to the left with one byte the first byte will remain outside
	// the window and the new byte will come in place. Let's say that the new byte is 00010011. Then the
	// window will contain this binary message: 0110001100010011. Polynomial of this message is:
	//
	//     P(x) = (x^14 + x^13 + x^9 + x^8) + (x^4 + x^1 + x^0)
	//
	// as you can see the polynomial of the byte 01100011 changed from
	// (x^6 + x^5 + x^1 + x^0) => (x^14 + x^13 + x^9 + x^8). The power of all 'x' are increased with 8.
	v = (v * pow(multiplier, 8)) % modulus

	// 3. Add the new polynomial of the byte to the value and MOD the result
	v = (v + polynomialOverGF2.value) % modulus
	//rfh.value = (((rfh.value+modulus-firstPolyGF2.value%modulus)*pow(multiplier, 8))%modulus + polynomialOverGF2.value) % modulus
	rfh.value = v
	rfh.table.add(polynomialOverGF2)
	return rfh.value
}

// pow power 'a' with 'b'.
func pow(a uint64, b int) uint64 {
	if b == 0 {
		return 1
	}
	n := uint64(1)
	for i := 0; i < b; i++ {
		// we use module because we don't want t\o work with big integers. Also, it is much faster.
		n = (n * a) % modulus
	}
	return n
}
