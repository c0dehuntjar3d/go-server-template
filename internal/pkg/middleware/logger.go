package middleware

import (
	"app/internal/pkg/logger"
	"app/internal/pkg/types"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Logging(next http.Handler, logger logger.Interface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		UUID := uuid.New().String()

		r = r.WithContext(
			context.WithValue(r.Context(), types.CtxKey("tx"), UUID),
		)

		logger.Info(
			fmt.Sprintf(
				"Request: [%s] -> Path: [%s] | UUID: %s",
				r.Method,
				r.RequestURI,
				UUID,
			),
		)

		next.ServeHTTP(w, r)

		logger.Debug(
			fmt.Sprintf(
				"Request Completed: [%s] -> Path: [%s] in [%v] | UUID: %s",
				r.Method,
				r.RequestURI,
				time.Since(start),
				UUID,
			),
		)
	})
}
