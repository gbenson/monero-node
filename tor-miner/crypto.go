package miner

import (
	crypto_rand "crypto/rand"
	"golang.org/x/crypto/argon2"
)

func randomBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := crypto_rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func deriveKey(password, salt []byte, keySizeBytes uint32) []byte {
	return argon2.IDKey(password, salt,
		2,         // number of iterations to perform
		1024*1024, // amount of memory to use, in KiB
		4,         // number of lanes (threads)
		keySizeBytes)
}
