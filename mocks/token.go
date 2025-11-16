package mocks

import (
	"timo/helper"
)

type MockJwtToken struct {
	Claims *helper.Claims
	Token  *string
	Err    error
}

func (m *MockJwtToken) Create(i *helper.Claims) (*string, error) {
	return m.Token, m.Err
}

func (m *MockJwtToken) Extract(tokenString string) (*helper.Claims, error) {
	return m.Claims, m.Err
}
