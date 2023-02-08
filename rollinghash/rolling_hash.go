package rollinghash

// Hash is the common interface implemented by all rolling hash functions.
type Hash interface {
	// Value return the value of the hash/sign
	Value() uint64

	// Next calculate the hash of the next rolling window. The window is shifted with one byte
	Next(byte) uint64
}
