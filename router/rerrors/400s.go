package rerrors

import "net/http"

func NotFound() RouteError {
	return NewRouteError(http.StatusNotFound, "Not Found", map[string]any{})
}

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
