package metrics_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gol4ng/fiberware/metrics"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Options(t *testing.T) {
	config := metrics.NewConfig(
		metrics.WithSplitStatus(true),
		metrics.WithObserveResponseSize(false),
		metrics.WithMeasureInflightRequests(false),
		metrics.WithIdentifierProvider(func(req *fiber.Request) string {
			return "my-personal-identifier"
		}),
	)
	assert.Equal(t, true, config.SplitStatus)
	assert.Equal(t, false, config.ObserveResponseSize)
	assert.Equal(t, false, config.MeasureInflightRequests)
	assert.Equal(t, "my-personal-identifier", config.IdentifierProvider(nil))
}
