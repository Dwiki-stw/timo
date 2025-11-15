package helper

import (
	"context"

	"google.golang.org/api/idtoken"
)

type Payload struct {
	GoogleID string
	Email    string
	Name     string
}

type TokenValidator interface {
	Validate(ctx context.Context, idToken string) (*Payload, error)
}

type GoogleValidator struct {
	clientID string
}

func NewGoogleValidator(clientID string) TokenValidator {
	return &GoogleValidator{clientID: clientID}
}

func (v *GoogleValidator) Validate(ctx context.Context, idToken string) (*Payload, error) {
	payload, err := idtoken.Validate(ctx, idToken, v.clientID)
	if err != nil {
		return nil, err
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	googleID := payload.Subject

	return &Payload{
		GoogleID: googleID,
		Email:    email,
		Name:     name,
	}, nil
}
