package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestLogger(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Path()

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		log.Debugf("[%s] %s %d %v", method, path, status, duration)
		return err
	}
}
