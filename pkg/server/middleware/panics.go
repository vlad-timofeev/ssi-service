package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tbd54566975/ssi-service/pkg/server/framework"
)

// Panics recovers from panics and converts the panic into an error
func Panics() framework.Middleware {
	mw := func(handler framework.Handler) framework.Handler {
		wrapped := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			v, ok := ctx.Value(framework.KeyRequestState).(*framework.RequestState)
			if !ok {
				return framework.NewShutdownError("request state missing from context.")
			}

			// defer a function to recover from a panic and set the err return
			// variable after the fact
			defer func() {
				if r := recover(); r != nil {
					// log the stack trace for this panic'd goroutine
					stack := debug.Stack()
					err = errors.Errorf("%s: \n%s", v.TraceID, stack)

					logrus.Infof("%s: PANIC : %s : \n%s", r, v.TraceID, stack)
				}
			}()

			// Call the next handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}

		return wrapped

	}

	return mw
}
