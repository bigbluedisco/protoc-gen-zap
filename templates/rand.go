package templates

import "crypto/rand"

const lowerCharBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

// LowerString generates a fixed-sized random string
// composed of lowercase letters and digits.
func lowerString(n int) string {
	output := make([]byte, n)
	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)
	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(lowerCharBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])
		// random % 64
		randomPos := random % uint8(l)
		// put into output
		output[pos] = lowerCharBytes[randomPos]
	}
	return string(output)
}
