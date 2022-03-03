package auth

import (
	"github.com/gofiber/fiber/v2"
)

// New middleware delegate the authentication process to the AuthenticateFunc
func New(cfg ...Config) fiber.Handler {
	config := fullFill(cfg...)
	return func(ctx *fiber.Ctx) error {
		if err := config.AuthenticateFunc(ctx); err != nil {
			return config.ErrorHandler(err, ctx)
		}
		if config.SuccessMiddleware != nil {
			return config.SuccessMiddleware(ctx)
		}
		return ctx.Next()
	}
}

type CredentialFinder func(request *fiber.Request) Credential
type AuthenticateFunc func(ctx *fiber.Ctx) error
type ErrorHandler func(err error, ctx *fiber.Ctx) error

type Config struct {
	AuthenticateFunc  AuthenticateFunc
	ErrorHandler      ErrorHandler
	SuccessMiddleware fiber.Handler
}

func fullFill(configs ...Config) Config {
	var config Config
	if len(configs) > 0 {
		config = configs[0]
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultErrorHandler
	}
	if config.AuthenticateFunc == nil {
		config.AuthenticateFunc = NewAuthenticateFunc()
	}
	return config
}

func DefaultErrorHandler(err error, ctx *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusUnauthorized, err.Error())
}

// NewAuthenticateFunc is an AuthenticateFunc that find, authenticate and hydrate credentials on the request context
func NewAuthenticateFunc(configs ...FuncConfig) AuthenticateFunc {
	config := fullFillFuncConfig(configs...)
	return func(ctx *fiber.Ctx) error {
		userCtx := ctx.UserContext()
		credential := config.CredentialFinder(ctx.Request())
		if config.Authenticator != nil {
			creds, err := config.Authenticator.Authenticate(userCtx, credential)
			if err != nil {
				return err
			}
			credential = creds
		}
		ctx.SetUserContext(CredentialToContext(userCtx, credential))
		return nil
	}
}

type FuncConfig struct {
	Authenticator    Authenticator
	CredentialFinder CredentialFinder
}

func fullFillFuncConfig(configs ...FuncConfig) FuncConfig {
	var config FuncConfig
	if len(configs) > 0 {
		config = configs[0]
	}
	if config.CredentialFinder == nil {
		config.CredentialFinder = DefaultCredentialFinder
	}
	return config
}

func DefaultCredentialFinder(request *fiber.Request) Credential {
	return FromHeader(request)()
}
