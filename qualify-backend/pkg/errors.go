package pkg

import "time"

func newError(path, httpError, code, message string) ErrorResponse {
	return ErrorResponse{
		Error:     httpError,
		Message:   message,
		Code:      code,
		Timestamp: time.Now(),
		Path:      path,
	}
}

func NotFound(path, message string) ErrorResponse {
	return newError(path, "Not Found", "NOT_FOUND", message)
}

func UnprocessableEntity(path, message string) ErrorResponse {
	return newError(path, "Unprocessable Entity", "UNPROCESSABLE_ENTITY", message)
}

func Unauthorized(path, message string) ErrorResponse {
	return newError(path, "Unauthorized", "UNAUTHORIZED", message)
}

func Forbidden(path, message string) ErrorResponse {
	return newError(path, "Forbidden", "FORBIDDEN", message)
}

func Conflict(path, message string) ErrorResponse {
	return newError(path, "Conflict", "CONFLICT", message)
}

func BadRequest(path, message string) ErrorResponse {
	return newError(path, "Bad Request", "BAD_REQUEST", message)
}

func Internal(path, message string) ErrorResponse {
	return newError(path, "Internal Server Error", "INTERNAL_ERROR", message)
}

func ValidationFailed(path, message string, validationErrors map[string]string) ErrorResponse {
	r := newError(path, "Validation Error", "VALIDATION_ERROR", message)
	r.ValidationErrors = validationErrors
	return r
}
