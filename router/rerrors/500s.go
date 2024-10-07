package rerrors

import (
	"errors"
	"net/http"
)

func InternalError(errs ...error) RouteError {
	return NewRouteError(http.StatusInternalServerError, "Internal server error", map[string]any{
		"errors": errors.Join(errs...),
	})
}
