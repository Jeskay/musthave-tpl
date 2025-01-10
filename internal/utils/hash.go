package utils

import "crypto/sha256"

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
