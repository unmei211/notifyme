package http_middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ApiVersioningMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		headerApiVersion := headers.Get("Api-Version")

		req.URL.Path = fmt.Sprintf("/%s%s", headerApiVersion, req.URL.Path)

		return next(c)
	}
}
