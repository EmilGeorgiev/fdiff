package rabin_test

import (
	"fmt"
	"testing"

	"github.com/EmilGeorgiev/fdiff/rabin"
	"github.com/stretchr/testify/assert"
)

func TestP(t *testing.T) {
	// a = 97  = 01100001 31, 30,29 ,28,27,26,25, 24
	// b = 98  = 01100010 23, 22,21, 20,19,18, 17, 16
	// c = 99  = 01100011 15, 14,13, 12,11,10, 9,8
	// d = 100 = 01100100 7,6,5,4,3,2,1,0
	polyAOfGM32 := rabin.Pow(257, 30) + rabin.Pow(257, 29) + rabin.Pow(257, 24)
	polyBOfGM32 := rabin.Pow(257, 22) + rabin.Pow(257, 21) + rabin.Pow(257, 17)
	polyCOfGM32 := rabin.Pow(257, 14) + rabin.Pow(257, 13) + rabin.Pow(257, 9) + rabin.Pow(257, 8)
	polyDOfGM32 := rabin.Pow(257, 5) + rabin.Pow(257, 6) + rabin.Pow(257, 2)

	h1 := (polyAOfGM32 + polyBOfGM32 + polyCOfGM32 + polyDOfGM32) % 8191
	fmt.Println(h1)

	// b = 98  = 01100010 31,30,29,28,27,26,25,24
	// c = 99  = 01100011 23,22,21,20,19,18,17,16
	// d = 100 = 01100100 15,14,13,12,11,10,9,8
	// e = 101 = 01100101 7,6,5,4,3,2,1,0
	poly1BOfGM32 := rabin.Pow(257, 30) + rabin.Pow(257, 29) + rabin.Pow(257, 25)
	poly1COfGM32 := rabin.Pow(257, 22) + rabin.Pow(257, 21) + rabin.Pow(257, 17) + rabin.Pow(257, 16)
	poly1DOfGM32 := rabin.Pow(257, 14) + rabin.Pow(257, 13) + rabin.Pow(257, 10)
	polyEOfGM32 := rabin.Pow(257, 6) + rabin.Pow(257, 5) + rabin.Pow(257, 2) + rabin.Pow(257, 0)

	h2 := (poly1BOfGM32 + poly1COfGM32 + poly1DOfGM32 + polyEOfGM32) % 8191
	fmt.Println(h2)

	h3 := (((h1+8191-polyAOfGM32%8191)*rabin.Pow(257, 8))%8191 + polyEOfGM32) % 8191
	fmt.Println(h3)
}

func TestCalculateRabinFingerPrintFor4Bytes(t *testing.T) {
	// Set up
	//p := []byte("This is a program that calculates the difference")
	p := []byte("abcd")
	// Action
	h := rabin.NewFingerprintHash(p)
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
	h := rabin.NewFingerprintHash(p)
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
	h := rabin.NewFingerprintHash(p)
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
//	h := rabin.NewFingerprintHash(p)
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
//	h := rabin.NewFingerprintHash(p)
//	actual := h.Value()
//
//	// Assert
//	expected := 5378
//	assert.EqualValues(t, expected, actual)
//}
//
//func TestNextRollingHash(t *testing.T) {
//	// Set up
//	rh := rabin.NewFingerprintHash("abcd")
//
//	// Action
//	actual := rh.Next('e')
//
//	// Assert
//	expected := 15
//	assert.EqualValues(t, expected, actual)
//}
