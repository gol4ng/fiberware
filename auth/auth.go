package auth

import (
	"github.com/gofiber/fiber/v2"
)

// New middleware delegate the authentication process to the AuthenticateFunc
func New(options ...AuthOption) fiber.Handler {
	config := NewAuthConfig(options...)
	return func(ctx *fiber.Ctx) error {
		if err := config.authenticateFunc(ctx); err != nil {
			return config.errorHandler(err, ctx)
		}
		if config.successMiddleware != nil {
			return config.successMiddleware(ctx)
		}
		return ctx.Next()
	}
}

type CredentialFinder func(request *fiber.Request) Credential
type AuthenticateFunc func(ctx *fiber.Ctx) error
type ErrorHandler func(err error, ctx *fiber.Ctx) error

// AuthOption defines a interceptor middleware configuration option
type AuthOption func(*AuthConfig)

type AuthConfig struct {
	authenticateFunc  AuthenticateFunc
	errorHandler      ErrorHandler
	successMiddleware fiber.Handler
}

func (o *AuthConfig) apply(options ...AuthOption) {
	for _, option := range options {
		option(o)
	}
}

func NewAuthConfig(options ...AuthOption) *AuthConfig {
	opts := &AuthConfig{
		errorHandler: DefaultErrorHandler,
	}
	opts.apply(options...)
	if opts.authenticateFunc == nil {
		opts.authenticateFunc = NewAuthenticateFunc()
	}
	return opts
}

func DefaultErrorHandler(err error, ctx *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusUnauthorized, err.Error())
}

// WithAuthenticateFunc will configure AuthenticateFunc option
func WithAuthenticateFunc(authenticateFunc AuthenticateFunc) AuthOption {
	return func(config *AuthConfig) {
		config.authenticateFunc = authenticateFunc
	}
}

// WithErrorHandler will configure ErrorHandler option
func WithErrorHandler(errorHandler ErrorHandler) AuthOption {
	return func(config *AuthConfig) {
		config.errorHandler = errorHandler
	}
}

// WithSuccessMiddleware will configure successMiddleware option
func WithSuccessMiddleware(middleware fiber.Handler) AuthOption {
	return func(config *AuthConfig) {
		config.successMiddleware = middleware
	}
}

// NewAuthenticateFunc is an AuthenticateFunc that find, authenticate and hydrate credentials on the request context
func NewAuthenticateFunc(options ...AuthFuncOption) AuthenticateFunc {
	config := NewAuthFuncConfig(options...)
	return func(ctx *fiber.Ctx) error {
		userCtx := ctx.UserContext()
		credential := config.credentialFinder(ctx.Request())
		if config.authenticator != nil {
			creds, err := config.authenticator.Authenticate(userCtx, credential)
			if err != nil {
				return err
			}
			credential = creds
		}
		ctx.SetUserContext(CredentialToContext(userCtx, credential))
		return nil
	}
}

// AuthFuncOption defines a AuthenticateFunc configuration option
type AuthFuncOption func(*AuthFuncConfig)

type AuthFuncConfig struct {
	authenticator    Authenticator
	credentialFinder CredentialFinder
}

func (o *AuthFuncConfig) apply(options ...AuthFuncOption) {
	for _, option := range options {
		option(o)
	}
}

func NewAuthFuncConfig(options ...AuthFuncOption) *AuthFuncConfig {
	opts := &AuthFuncConfig{
		credentialFinder: DefaultCredentialFinder,
	}
	opts.apply(options...)
	return opts
}

func DefaultCredentialFinder(request *fiber.Request) Credential {
	return FromHeader(request)()
}

// WithAuthenticator will configure Authenticator option
func WithAuthenticator(authenticator Authenticator) AuthFuncOption {
	return func(config *AuthFuncConfig) {
		config.authenticator = authenticator
	}
}

// WithCredentialFinder will configure AuthenticateFunc option
func WithCredentialFinder(credentialFinder CredentialFinder) AuthFuncOption {
	return func(config *AuthFuncConfig) {
		config.credentialFinder = credentialFinder
	}
}
