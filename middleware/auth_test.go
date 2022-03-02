package middleware_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gol4ng/fiberware/auth"
	"github.com/gol4ng/fiberware/middleware"
	"github.com/gol4ng/fiberware/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCredentialFinder(t *testing.T) {
	tests := []struct {
		authorizationHeader  string
		xAuthorizationHeader string
		expectedCredential   string
	}{
		{
			authorizationHeader:  "",
			xAuthorizationHeader: "",
			expectedCredential:   "",
		},
		{
			authorizationHeader:  "Foo",
			xAuthorizationHeader: "",
			expectedCredential:   "Foo",
		},
		{
			authorizationHeader:  "",
			xAuthorizationHeader: "Foo",
			expectedCredential:   "Foo",
		},
		{
			authorizationHeader:  "Foo",
			xAuthorizationHeader: "Bar",
			expectedCredential:   "Foo",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s%s", tt.authorizationHeader, tt.xAuthorizationHeader), func(t *testing.T) {
			request := &fiber.Request{}
			request.Header.Set(fiber.HeaderAuthorization, tt.authorizationHeader)
			request.Header.Set(auth.XAuthorizationHeader, tt.xAuthorizationHeader)

			assert.Equal(t, auth.Credential(tt.expectedCredential), middleware.DefaultCredentialFinder(request))
		})
	}
}

func TestDefaultErrorHandler(t *testing.T) {
	err := middleware.DefaultErrorHandler(errors.New("my_fake_error"), nil)
	e, ok := err.(*fiber.Error)

	assert.True(t, ok)
	assert.Equal(t, 401, e.Code)
	assert.EqualError(t, e, "my_fake_error")
}

func TestAuthentication(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Authentication(middleware.WithAuthenticateFunc(func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(auth.CredentialToContext(ctx.UserContext(), "my_allowed_credential"))
		return nil
	})))

	var innerContext context.Context
	handlerCalled := false

	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		innerContext = ctx.UserContext()
		return nil
	})
	request := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.NotEqual(t, request.Context(), innerContext) // decorated context
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_allowed_credential", auth.CredentialFromContext(innerContext))
}

func TestAuthenticationWithoutAuthenticator(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Authentication())

	var innerContext context.Context
	handlerCalled := false

	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		innerContext = ctx.UserContext()
		return nil
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "my_allowed_credential")
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.NotEqual(t, request.Context(), innerContext) // decorated context
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_allowed_credential", auth.CredentialFromContext(innerContext))
}

func TestAuthentication_Error(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Authentication(middleware.WithAuthenticateFunc(func(ctx *fiber.Ctx) error {
		return errors.New("my_authenticate_error")
	})))

	handlerCalled := false
	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		return nil
	})

	request := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.False(t, handlerCalled)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "my_authenticate_error", string(body))
}

func TestAuthentication_WithSuccessMiddleware(t *testing.T) {
	app := fiber.New()

	successMiddlewareCalled := false
	app.Use(middleware.Authentication(
		middleware.WithAuthenticateFunc(func(ctx *fiber.Ctx) error {
			ctx.SetUserContext(auth.CredentialToContext(ctx.UserContext(), "my_allowed_credential"))
			return nil
		}),
		middleware.WithSuccessMiddleware(func(ctx *fiber.Ctx) error {
			successMiddlewareCalled = true
			assert.Equal(t, "my_allowed_credential", auth.CredentialFromContext(ctx.UserContext()))
			return ctx.Next()
		}),
	))

	var innerContext context.Context
	handlerCalled := false
	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		innerContext = ctx.UserContext()
		return nil
	})

	request := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.True(t, successMiddlewareCalled)
	assert.NotEqual(t, request.Context(), innerContext) // decorated context
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_allowed_credential", auth.CredentialFromContext(innerContext))
}

func TestAuthentication_WithErrorHandler(t *testing.T) {
	app := fiber.New()

	errorHandlerCalled := false
	app.Use(middleware.Authentication(
		middleware.WithAuthenticateFunc(func(ctx *fiber.Ctx) error {
			return errors.New("my_authenticate_error")
		}),
		middleware.WithErrorHandler(func(err error, ctx *fiber.Ctx) error {
			errorHandlerCalled = true
			assert.EqualError(t, err, "my_authenticate_error")
			return nil
		}),
	))

	handlerCalled := false
	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		return nil
	})

	request := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.False(t, handlerCalled)
	assert.True(t, errorHandlerCalled)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestNewAuthenticateFunc(t *testing.T) {
	app := fiber.New()

	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", context.TODO(), "my_allowed_credential").Return("my_authenticate_credential", nil)

	app.Use(middleware.Authentication(
		middleware.WithAuthenticateFunc(middleware.NewAuthenticateFunc(middleware.WithAuthenticator(authenticator))),
	))

	var innerContext context.Context
	handlerCalled := false

	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		innerContext = ctx.UserContext()
		return nil
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Set("Authorization", "my_allowed_credential")
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.NotEqual(t, request.Context(), innerContext) // decorated context
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_authenticate_credential", auth.CredentialFromContext(innerContext))

	authenticator.AssertExpectations(t)
}

func TestNewAuthenticateFunc_WithCredentialFinder(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Authentication(
		middleware.WithAuthenticateFunc(middleware.NewAuthenticateFunc(
			middleware.WithCredentialFinder(func(request *fiber.Request) auth.Credential {
				return "my_finded_credential"
			}),
		)),
	))

	var innerContext context.Context
	handlerCalled := false

	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		innerContext = ctx.UserContext()
		return nil
	})
	request := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.NotEqual(t, request.Context(), innerContext) // decorated context
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "my_finded_credential", auth.CredentialFromContext(innerContext))
}

func TestNewAuthenticateFunc_Error(t *testing.T) {
	app := fiber.New()

	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", context.TODO(), "").Return("", errors.New("my_authenticator_error"))

	app.Use(middleware.Authentication(middleware.WithAuthenticateFunc(
		middleware.NewAuthenticateFunc(middleware.WithAuthenticator(authenticator))),
	))

	handlerCalled := false

	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		return nil
	})
	request := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.False(t, handlerCalled)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "my_authenticator_error", string(body))

	authenticator.AssertExpectations(t)
}
