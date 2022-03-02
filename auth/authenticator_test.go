package auth_test

import (
	"context"
	"errors"
	"github.com/gol4ng/fiberware/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticatorFunc_Authenticate(t *testing.T) {
	called := false
	a := auth.AuthenticatorFunc(func(context.Context, auth.Credential) (auth.Credential, error) {
		called = true
		return "data", errors.New("my_error")
	})

	creds, err := a.Authenticate(context.TODO(), "my_input_credential")

	assert.True(t, called)
	assert.EqualError(t, err, "my_error")
	assert.Equal(t, "data", creds)
}
