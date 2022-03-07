package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func New(recorder Recorder, configs ...Config) func(ctx *fiber.Ctx) error {
	config := fullFill(configs...)
	return func(ctx *fiber.Ctx) error {
		handlerName := config.IdentifierProvider(ctx.Request())
		if !config.DisableInflightRequests {
			recorder.AddInflightRequests(ctx.Context(), handlerName, 1)
			defer recorder.AddInflightRequests(ctx.Context(), handlerName, -1)
		}

		start := time.Now()
		defer func() {
			code := strconv.Itoa(ctx.Response().StatusCode())
			if !config.RawStatus {
				code = fmt.Sprintf("%dxx", ctx.Response().StatusCode()/100)
			}

			recorder.ObserveHTTPRequestDuration(ctx.Context(), handlerName, time.Since(start), ctx.Method(), code)

			if !config.DisableResponseSize {
				recorder.ObserveHTTPResponseSize(ctx.Context(), handlerName, int64(len(ctx.Response().Body())), ctx.Method(), code)
			}
		}()

		return ctx.Next()
	}
}

type Config struct {
	// func that allows you to provide a strategy to identify/group metrics
	// you can group metrics by request host/url/... or app name ...
	// by default, we group metrics by request url
	IdentifierProvider func(req *fiber.Request) string
	// if set to true, each response status will be represented by a metrics
	// if set to false, response status codes will be grouped by first digit (204/201/200/... -> 2xx; 404/403/... -> 4xx)
	RawStatus bool
	// if set to true, recorder will add a responseSize metric
	DisableResponseSize bool
	// if set to true, recorder will add a metric representing the number of inflight requests
	DisableInflightRequests bool
}

// NewConfig returns a new metrics configuration with all options applied
func fullFill(configs ...Config) Config {
	var config Config
	if len(configs) > 0 {
		config = configs[0]
	}
	if config.IdentifierProvider == nil {
		config.IdentifierProvider = func(req *fiber.Request) string {
			return req.URI().String()
		}
	}
	return config
}
