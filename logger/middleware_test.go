package logger_test

import (
	"context"
	"errors"
	"github.com/gol4ng/logger"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	fiberware_logger "github.com/gol4ng/fiberware/logger"
	"github.com/gol4ng/fiberware/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	app := fiber.New()
	l := &mocks.LoggerInterface{}
	l.On(
		"Debug",
		"http server received GET /",
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
	)
	l.On(
		"Log",
		mock.AnythingOfType("string"),
		logger.InfoLevel,
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
	)

	app.Use(fiberware_logger.New(l, fiberware_logger.Config{InjectOnContext: true}))

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
	assert.Equal(t, l, fiberware_logger.LoggerFromContext(innerContext))

	l.AssertExpectations(t)
}

func TestNew_WithError(t *testing.T) {
	app := fiber.New()
	l := &mocks.LoggerInterface{}
	l.On(
		"Debug",
		"http server received GET /",
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
	)
	l.On(
		"Log",
		mock.AnythingOfType("string"),
		logger.ErrorLevel,
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
		mock.AnythingOfType("logger.Field"),
	)

	app.Use(fiberware_logger.New(l, fiberware_logger.Config{InjectOnContext: true}))

	var innerContext context.Context
	handlerCalled := false

	app.Get("/", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		innerContext = ctx.UserContext()
		return errors.New("my_custom_error")
	})
	request := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.NotEqual(t, request.Context(), innerContext) // decorated context
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, l, fiberware_logger.LoggerFromContext(innerContext))

	l.AssertExpectations(t)
}
