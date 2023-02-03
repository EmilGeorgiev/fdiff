package rabin_test

import (
	"testing"

	"github.com/EmilGeorgiev/fdiff/rabin"
	"github.com/stretchr/testify/assert"
)

func TestCalculateRabinFingerPrintFor4Bytes(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	p := "abcd"
	// Action
	h := rabin.NewFingerprintHash(p)
	actual := h.Value()

	// Assert
	expected := 3873
	assert.EqualValues(t, expected, actual)
}

func TestCalculateRabinFingerPrintFor48Bytes(t *testing.T) {
	// Set up
	p := "This is a program that calculates the difference"
	// Action
	h := rabin.NewFingerprintHash(p)
	actual := h.Value()

	// Assert
	expected := 1477
	assert.EqualValues(t, expected, actual)
}

func TestCalculateRabinFingerPrintOnNonASCIICharacters(t *testing.T) {
	// Set up
	p := "абвгྠྡ ྡ"

	// Action
	h := rabin.NewFingerprintHash(p)
	actual := h.Value()

	// Assert
	expected := 5378
	assert.EqualValues(t, expected, actual)
}
