package v1

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/pkg/errors"
)

func (server *Server) sentryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		value, ok := c.Get("error")
		if !ok {
			return
		}
		if value == nil {
			return
		}
		hub := sentrygin.GetHubFromContext(c)
		if hub == nil {
			return
		}
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("integration", "errors")
			if errCapture, ok := value.(error); ok {
				if isSkipError(errCapture) {
					return
				}
				hub.CaptureException(errCapture)
			} else {
				hub.CaptureMessage(fmt.Sprintf("error message: %s", value.(string)))
			}
		})
	}
}

func isSkipError(err error) bool {
	if errors.Is(err, errs.ErrProductNotFound) {
		return true
	}
	if errors.Is(err, validator.ErrPassed) {
		return true
	}
	if errors.Is(err, errs.ErrStoreNotFound) {
		return true
	}
	return false
}
