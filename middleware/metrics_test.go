package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	prom "github.com/gol4ng/fiberware/metrics/prometheus"
	"github.com/gol4ng/fiberware/middleware"
	"github.com/gol4ng/fiberware/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetrics(t *testing.T) {
	app := fiber.New()
	recorder := &mocks.Recorder{}

	handlerIdentifier := "http://example.com/path"
	recorder.On("AddInflightRequests",
		mock.AnythingOfType("*fasthttp.RequestCtx"),
		handlerIdentifier,
		1,
	)
	recorder.On("AddInflightRequests",
		mock.AnythingOfType("*fasthttp.RequestCtx"),
		handlerIdentifier,
		-1,
	)
	recorder.On("ObserveHTTPRequestDuration",
		mock.AnythingOfType("*fasthttp.RequestCtx"),
		handlerIdentifier,
		mock.AnythingOfType("time.Duration"),
		"GET",
		"2xx",
	)
	recorder.On("ObserveHTTPResponseSize",
		mock.AnythingOfType("*fasthttp.RequestCtx"),
		handlerIdentifier,
		mock.AnythingOfType("int64"),
		"GET",
		"2xx",
	)

	app.Use(
		middleware.Metrics(recorder),
	)

	handlerCalled := false
	app.Get("/path", func(ctx *fiber.Ctx) error {
		handlerCalled = true
		return nil
	})

	request := httptest.NewRequest("GET", "/path", nil)
	resp, err := app.Test(request)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	recorder.AssertExpectations(t)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleMetrics() {
	app := fiber.New()

	recorder := prom.NewRecorder(prom.Config{}).RegisterOn(nil)
	app.Use(middleware.Metrics(recorder))

	app.Get("/path", func(ctx *fiber.Ctx) error {
		return nil
	})
	//Output:
}
