package mrpack

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

type HashType string

const (
	HashSHA1   HashType = "sha1"
	HashSHA512 HashType = "sha512"
)

func VerifyFile(path string, hashes map[string]string) (bool, error) {
	if len(hashes) == 0 {
		return true, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open for hash: %w", err)
	}
	defer f.Close()

	var hasher hash.Hash
	var expected string

	if h, ok := hashes["sha512"]; ok {
		hasher = sha512.New()
		expected = h
	} else if h, ok := hashes["sha1"]; ok {
		hasher = sha1.New()
		expected = h
	} else {
		return true, nil
	}

	if _, err := io.Copy(hasher, f); err != nil {
		return false, fmt.Errorf("hash read: %w", err)
	}

	computed := hex.EncodeToString(hasher.Sum(nil))
	return computed == expected, nil
}

func HashFileSHA1(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func HashFileSHA512(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
