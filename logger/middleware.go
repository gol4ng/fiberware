package logger

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gol4ng/logger"
)

func New(log logger.LoggerInterface, configs ...Config) func(ctx *fiber.Ctx) error {
	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)
	config := fullFill(configs...)
	return func(ctx *fiber.Ctx) error {
		// Set error handler once
		once.Do(func() {
			errHandler = ctx.App().ErrorHandler
		})

		startTime := time.Now()

		if config.InjectOnContext {
			ctx.SetUserContext(LoggerToContext(ctx.UserContext(), log))
		}

		currentLoggerContext := logger.NewContext().
			Add("http_header", ctx.GetReqHeaders()).
			Add("http_kind", "server").
			Add("http_method", ctx.Method()).
			Add("http_url", ctx.OriginalURL()).
			Add("http_start_time", startTime.Format(time.RFC3339))

		if d, ok := ctx.Context().Deadline(); ok {
			currentLoggerContext.Add("http_request_deadline", d.Format(time.RFC3339))
		}

		defer func() {
			duration := time.Since(startTime)
			currentLoggerContext.
				Add("http_duration", duration.Seconds()).
				Add("http_response_header", ctx.GetRespHeaders())

			if err := recover(); err != nil {
				currentLoggerContext.Add("http_panic", err)
				log.Critical(fmt.Sprintf("http server panic %s %s [duration:%s]", ctx.Method(), ctx.OriginalURL(), duration), *currentLoggerContext.Slice()...)
				panic(err)
			}

			currentLoggerContext.Add("http_status", http.StatusText(ctx.Response().StatusCode())).
				Add("http_status_code", ctx.Response().StatusCode()).
				Add("http_response_length", len(ctx.Response().Body()))

			log.Log(
				fmt.Sprintf(
					"http server %s %s [status_code:%d, duration:%s, content_length:%d]",
					ctx.Method(), ctx.OriginalURL(), ctx.Response().StatusCode(), duration, len(ctx.Response().Body()),
				),
				config.LevelFunc(ctx.Response().StatusCode()),
				*currentLoggerContext.Slice()...,
			)
		}()

		log.Debug(fmt.Sprintf("http server received %s %s", ctx.Method(), ctx.OriginalURL()), *currentLoggerContext.Slice()...)

		chainErr := ctx.Next()

		if chainErr != nil {
			if err := errHandler(ctx, chainErr); err != nil {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
			}
		}
		return nil
	}
}

func fullFill(configs ...Config) Config {
	var config Config
	if len(configs) > 0 {
		config = configs[0]
	}
	if config.LevelFunc == nil {
		config.LevelFunc = levelFunc
	}
	return config
}

func levelFunc(statusCode int) logger.Level {
	switch {
	case statusCode < fiber.StatusBadRequest:
		return logger.InfoLevel
	case statusCode < fiber.StatusInternalServerError:
		return logger.WarningLevel
	}
	return logger.ErrorLevel
}
