package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		stop := time.Now()
		latency := stop.Sub(start)

		req := c.Request()
		res := c.Response()

		fmt.Printf("[%s] %s %s %s %d %s\n",
			stop.Format(time.RFC3339),
			c.RealIP(),
			req.Method,
			req.RequestURI,
			res.Status,
			latency,
		)

		return err
	}
}
