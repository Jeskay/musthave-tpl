package utils

import "crypto/sha256"

// HashBytes applies sha256 hash function to given data and key, returns the result.
func HashBytes(data []byte, key string) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	if _, err := h.Write([]byte(key)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
