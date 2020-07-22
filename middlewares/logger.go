package middlewares

import (
	"github.com/MoonSHRD/logger"
	"net/http"
	"time"
)

// Logger logs the current request to the console printing the date, HTTP method, path and elapsed time
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next(res, req)
		logger.Infof("[%s] %q %v\n", req.Method, req.URL.String(), time.Since(start))
	}
}
