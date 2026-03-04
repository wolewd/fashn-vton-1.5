package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"backend/utils"
)

type RateLimiterConfig struct {
	Requests int
	Window   time.Duration
}

type clientData struct {
	Count     int
	ExpiresAt time.Time
}

func NewRateLimiter(config RateLimiterConfig) echo.MiddlewareFunc {
	clients := make(map[string]*clientData)
	var mu sync.Mutex

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			now := time.Now()

			mu.Lock()
			data, exists := clients[ip]

			if !exists || now.After(data.ExpiresAt) {
				data = &clientData{
					Count:     1,
					ExpiresAt: now.Add(config.Window),
				}
				clients[ip] = data
			} else {
				data.Count++
			}

			mu.Unlock()

			if data.Count > config.Requests {
				return utils.JSONError(c, http.StatusTooManyRequests, "Rate limit exceeded", map[string]any{
					"limit": 5,
					"window": "1s",
})
			}

			return next(c)
		}
	}
}
