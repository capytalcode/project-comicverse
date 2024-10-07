package rerrors

import "net/http"

func MissingParameters(params []string) RouteError {
	return NewRouteError(http.StatusBadRequest, "Missing parameters", map[string]any{
		"missing_parameters": params,
	})
}

func MethodNowAllowed(method string, allowedMethods []string) RouteError {
	return NewRouteError(http.StatusMethodNotAllowed, "Method not allowed", map[string]any{
		"method":          method,
		"allowed_methods": allowedMethods,
	})
}
