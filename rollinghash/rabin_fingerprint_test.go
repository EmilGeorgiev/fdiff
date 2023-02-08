package rollinghash_test

import (
	"github.com/EmilGeorgiev/fdiff/rollinghash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRabinFingerPrintFor4Bytes(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	p := []byte("abcd")
	// Action
	h := rollinghash.NewRabinFingerprint(p)
	actual := h.Value()

	// Assert
	expected := 3873
	assert.EqualValues(t, expected, actual)
}

func TestCalculateRabinFingerPrintFor4Bytess(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	p := []byte("bcde")
	// Action
	h := rollinghash.NewRabinFingerprint(p)
	actual := h.Value()

	// Assert
	expected := 7293
	assert.EqualValues(t, expected, actual)
}

func TestCalculateNextRabinFingerPrintFor4Bytes(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	p := []byte("abcd")
	// Action
	h := rollinghash.NewRabinFingerprint(p)
	b := []byte("e")
	actual := h.Next(b[0])

	// Assert
	expected := 7293
	assert.EqualValues(t, expected, actual)
}

//
//func TestCalculateRabinFingerPrintFor48Bytes(t *testing.T) {
//	// Set up
//	p := "This is a program that calculates the difference"
//	// Action
//	h := rabin.NewHash(p)
//	actual := h.Value()
//
//	// Assert
//	expected := 1477
//	assert.EqualValues(t, expected, actual)
//}
//
//func TestCalculateRabinFingerPrintOnNonASCIICharacters(t *testing.T) {
//	// Set up
//	p := "абвгྠྡ ྡ"
//
//	// Action
//	h := rabin.NewHash(p)
//	actual := h.Value()
//
//	// Assert
//	expected := 5378
//	assert.EqualValues(t, expected, actual)
//}
//
//func TestNextRollingHash(t *testing.T) {
//	// Set up
//	rh := rabin.NewHash("abcd")
//
//	// Action
//	actual := rh.Next('e')
//
//	// Assert
//	expected := 15
//	assert.EqualValues(t, expected, actual)
//}
