package fiberware

import (
	"github.com/gofiber/fiber/v2"
)

func NextHandler(ctx *fiber.Ctx) error {
	return ctx.Next()
}
