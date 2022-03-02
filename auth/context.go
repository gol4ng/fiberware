package auth

import (
	"context"
)

type ctxKey uint

const (
	credentialKey ctxKey = iota
)

func CredentialToContext(ctx context.Context, credential Credential) context.Context {
	return context.WithValue(ctx, credentialKey, credential)
}

func CredentialFromContext(ctx context.Context) Credential {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(credentialKey)
	if value == nil {
		return nil
	}
	credential, ok := value.(Credential)
	if !ok {
		return nil
	}

	return credential
}
