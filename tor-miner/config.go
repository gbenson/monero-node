package miner

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	kdfSaltSizeBytes   = 16
	sealKeySizeBytes   = 32
	sealNonceSizeBytes = 24
)

//go:embed sealed.config
var sealedDefaultConfig []byte

type Config struct {
	Pool     APIEndpoint `json:"default_pool"`
	Monitor  APIEndpoint `json:"monitor"`
	Resolver Resolver    `json:"resolver,omitempty"`
}

func DefaultConfig(password string) (*Config, error) {
	salt := sealedDefaultConfig[:kdfSaltSizeBytes]
	msg := sealedDefaultConfig[kdfSaltSizeBytes:]

	fmt.Println("tor-miner: unpacking credentials")
	keyBytes := deriveKey([]byte(password), salt, sealKeySizeBytes)

	key := (*[sealKeySizeBytes]byte)(keyBytes)
	nonce := (*[sealNonceSizeBytes]byte)(msg[:sealNonceSizeBytes])
	encrypted := msg[sealNonceSizeBytes:]

	encoded, ok := secretbox.Open(nil, encrypted, nonce, key)
	if !ok {
		return nil, errors.New("bad passphrase")
	}

	cfg := &Config{}
	err := json.Unmarshal(encoded, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
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
