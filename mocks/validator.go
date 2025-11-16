package mocks

import (
	"context"
	"timo/helper"
)

type MockTokenValidator struct {
	Payload *helper.Payload
	Err     error
}

func (m *MockTokenValidator) Validate(ctx context.Context, idToken string) (*helper.Payload, error) {
	return m.Payload, m.Err
}
