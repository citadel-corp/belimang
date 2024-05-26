package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/citadel-corp/go-project-template/internal/common/response"
	"github.com/rs/zerolog/log"
)

func PanicRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				if r != http.ErrAbortHandler {
					log.Error().Msg(fmt.Sprintf("Recovered from panic: %s", string(debug.Stack())))
				}
				response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
					Message: "Internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
