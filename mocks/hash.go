package mocks

import (
	"errors"
)

type MockHasher struct {
	ShouldFail bool
}

func (h *MockHasher) Compare(hash string, password string) error {
	if h.ShouldFail {
		return errors.New("simulated error compare")
	}
	return nil
}

func (h *MockHasher) Hash(password string) (string, error) {
	if h.ShouldFail {
		return "", errors.New("simulated error hashing")
	}
	return password, nil
}
