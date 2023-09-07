package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/v2"
)

func SetTraceId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			tr := apm.TransactionFromContext(c.Context())
			traceId := tr.TraceContext().Trace.String()
			c.Set("x-elastic-trace-id", traceId)
		}()

		return c.Next()
	}
}
