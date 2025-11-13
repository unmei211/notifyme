package http_middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CorrelationGenerateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		id := uuid.New().String()

		c.Response().Header().Set(echo.HeaderXCorrelationID, id)
		mutateReq := req.WithContext(context.WithValue(req.Context(), echo.HeaderXCorrelationID, id))
		c.SetRequest(mutateReq)

		return next(c)
	}
}
