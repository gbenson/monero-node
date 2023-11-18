package miner

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	kdfSaltSizeBytes   = 16
	sealKeySizeBytes   = 32
	sealNonceSizeBytes = 24
)

type Config struct {
	Pool    APIEndpoint `json:"pool"`
	Monitor APIEndpoint `json:"monitor"`
}

func (c *Config) Seal(password string) ([]byte, error) {
	encoded, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	salt, err := randomBytes(kdfSaltSizeBytes)
	if err != nil {
		return nil, err
	}

	fmt.Println("Starting KDF")
	start := time.Now()
	buf := deriveKey([]byte(password), salt, sealKeySizeBytes)
	limit := time.Now()
	fmt.Println("KDF took:", limit.Sub(start))
	key := (*[sealKeySizeBytes]byte)(buf)

	buf, err = randomBytes(sealNonceSizeBytes)
	if err != nil {
		return nil, err
	}

	nonce := (*[sealNonceSizeBytes]byte)(buf)

	buf = append(salt, buf...)
	return secretbox.Seal(buf, encoded, nonce, key), nil
}
