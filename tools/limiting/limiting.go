package limiting

import (
	"net/http"

	"golang.org/x/time/rate"
)

var defaultLimiter = NewLimiter(1, 1)

func NewLimiter(limit float64, burst int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(limit), burst)
}

func SetLimiter(l *rate.Limiter) {
	defaultLimiter = l
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if defaultLimiter.Allow() == false {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
