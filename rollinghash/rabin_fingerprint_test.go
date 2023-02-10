package rollinghash_test

import (
	"testing"

	"github.com/EmilGeorgiev/fdiff/rollinghash"
	"github.com/stretchr/testify/assert"
)

var (
	// a = 97 = 01100001 = (127^6 + 127^5 +127^0) % 8191 = 3002
	fingerprintOfA = (pow(127, 6) + pow(127, 5) + 1) % 8191

	// b = 98 =  01100010 = (127^6 + 127^5 +127^1) % 8191 = 3128
	fingerprintOfB = (pow(127, 6) + pow(127, 5) + 127) % 8191

	// c = 99 =  01100011 = (127^6 + 127^5 +127^1 + 127^0) % 8191 = 3129
	fingerprintOfC = (pow(127, 6) + pow(127, 5) + 127 + 1) % 8191

	// d = 100 = 01100100 = (127^6 + 127^5 +127^2) % 8191 = 2748
	fingerprintOfD = (pow(127, 6) + pow(127, 5) + (127 * 127 % 8191)) % 8191

	// d = 101 = 01100101 = (127^6 + 127^5 + 127^2 + 1) % 8191 = 2749
	fingerprintOfE = (pow(127, 6) + pow(127, 5) + (127 * 127 % 8191) + 1) % 8191

	// f = 102 = 01100110 = (127^6 + 127^5 + 127^2 + 127 ) % 8191 = 2875
	fingerprintOfF = (pow(127, 6) + pow(127, 5) + (127 * 127 % 8191) + 127) % 8191

	// g = 103 = 01100111 = (127^6 + 127^5 + 127^2 + 127 + 1) % 8191 = 2876
	fingerprintOfG = (pow(127, 6) + pow(127, 5) + (127 * 127 % 8191) + 127 + 1) % 8191
)

func TestNewRabinFingerprint(t *testing.T) {
	// SetUp

	// a = 97 = 01100001 = (127^6 + 127^5 +127^0) % 8191 = 3002
	polyOfA := pow(127, 6) + pow(127, 5) + 1
	fingerprintOfA := polyOfA % 8191

	// b = 98 =  01100010 = (127^6 + 127^5 +127^1) % 8191 = 3128
	polyOfB := pow(127, 6) + pow(127, 5) + 127
	fingerprintOfB := polyOfB % 8191

	// c = 99 =  01100011 = (127^6 + 127^5 +127^1 + 127^0) % 8191 = 3129
	polyOfC := pow(127, 6) + pow(127, 5) + 127 + 1
	fingerprintOfC := polyOfC % 8191

	// d = 100 = 01100100 = (127^6 + 127^5 +127^2) % 8191 = 2748
	polyOfD := pow(127, 6) + pow(127, 5) + ((127 * 127) % 8191)
	fingerprintOfD := polyOfD % 8191

	cases := []struct {
		name                string
		text                string
		expectedFingerPrint uint64
	}{
		{name: "fingerprint for a", text: "a", expectedFingerPrint: fingerprintOfA},
		{name: "fingerprint for b", text: "b", expectedFingerPrint: fingerprintOfB},
		{name: "fingerprint for c", text: "c", expectedFingerPrint: fingerprintOfC},
		{name: "fingerprint for d", text: "d", expectedFingerPrint: fingerprintOfD},
		{
			name:                "fingerprint for ab",
			text:                "ab",
			expectedFingerPrint: ((fingerprintOfA * pow(127, 8)) + fingerprintOfB) % 8191,
		},
		{
			name:                "fingerprint for bc",
			text:                "bc",
			expectedFingerPrint: ((fingerprintOfB * pow(127, 8)) + fingerprintOfC) % 8191,
		},
		{
			name:                "fingerprint for cd",
			text:                "cd",
			expectedFingerPrint: ((fingerprintOfC * pow(127, 8)) + fingerprintOfD) % 8191,
		},
		{
			name:                "fingerprint for da",
			text:                "da",
			expectedFingerPrint: ((fingerprintOfD * pow(127, 8)) + fingerprintOfA) % 8191,
		},
		{
			name: "fingerprint for abcd",
			text: "abcd",
			expectedFingerPrint: ((fingerprintOfA * pow(127, 24)) +
				(fingerprintOfB * pow(127, 16)) +
				(fingerprintOfC * pow(127, 8)) +
				fingerprintOfD) % 8191,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Action
			actual := rollinghash.NewRabinFingerprint([]byte(c.text))

			// Assert
			assert.EqualValues(t, c.expectedFingerPrint, actual.Value())
		})
	}
}

func TestNext(t *testing.T) {
	// SetUp
	h := rollinghash.NewRabinFingerprint([]byte("abcd"))

	cases := []struct {
		name                string
		nextByte            byte
		expectedFingerPrint uint64
	}{
		{
			name:     "next character is 'e'",
			nextByte: 101,
			expectedFingerPrint: ((fingerprintOfB * pow(127, 24)) +
				(fingerprintOfC * pow(127, 16)) +
				(fingerprintOfD * pow(127, 8)) +
				fingerprintOfE) % 8191,
		},
		{
			name:     "next character is 'f'",
			nextByte: 102,
			expectedFingerPrint: ((fingerprintOfC * pow(127, 24)) +
				(fingerprintOfD * pow(127, 16)) +
				(fingerprintOfE * pow(127, 8)) +
				fingerprintOfF) % 8191,
		},
		{
			name:     "next character is 'g'",
			nextByte: 103,
			expectedFingerPrint: ((fingerprintOfD * pow(127, 24)) +
				(fingerprintOfE * pow(127, 16)) +
				(fingerprintOfF * pow(127, 8)) +
				fingerprintOfG) % 8191,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Action
			actual := h.Next(c.nextByte)

			// Assert
			assert.EqualValues(t, c.expectedFingerPrint, actual)
		})
	}
}

func pow(a uint64, b int) uint64 {
	if b == 0 {
		return 1
	}
	n := uint64(1)
	for i := 0; i < b; i++ {
		// we use module because we don't want t\o work with big integers. Also, it is much faster.
		n = (n * a) % 8191
	}
	return n
}
