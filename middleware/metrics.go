package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gol4ng/fiberware/metrics"
)

func Metrics(recorder metrics.Recorder, options ...metrics.Option) func(ctx *fiber.Ctx) error {
	config := metrics.NewConfig(options...)
	return func(ctx *fiber.Ctx) error {
		handlerName := config.IdentifierProvider(ctx.Request())
		if config.MeasureInflightRequests {
			recorder.AddInflightRequests(ctx.Context(), handlerName, 1)
			defer recorder.AddInflightRequests(ctx.Context(), handlerName, -1)
		}

		start := time.Now()
		defer func() {
			code := strconv.Itoa(ctx.Response().StatusCode())
			if !config.SplitStatus {
				code = fmt.Sprintf("%dxx", ctx.Response().StatusCode()/100)
			}

			recorder.ObserveHTTPRequestDuration(ctx.Context(), handlerName, time.Since(start), ctx.Method(), code)

			if config.ObserveResponseSize {
				recorder.ObserveHTTPResponseSize(ctx.Context(), handlerName, int64(len(ctx.Response().Body())), ctx.Method(), code)
			}
		}()

		return ctx.Next()
	}
}
