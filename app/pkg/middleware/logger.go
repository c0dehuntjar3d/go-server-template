package middleware

import (
	"app/pkg/common"
	"app/pkg/logger"
	"app/pkg/types"
	"context"
	"fmt"
	"net/http"
	"time"
)

func Logging(next http.Handler, logger logger.Interface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		UUID := common.GenerateUUID()

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
