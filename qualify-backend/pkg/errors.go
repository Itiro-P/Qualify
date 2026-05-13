package pkg

import "time"

func NotFound(path string, message string) ErrorResponse {
	return ErrorResponse{
		Error:     "Not Found",
		Message:   message,
		Code:      "NOT_FOUND",
		Timestamp: time.Now(),
		Path:      path,
	}
}

func Unauthorized(path string, message string) ErrorResponse {
	return ErrorResponse{
		Error:     "Unauthorized",
		Message:   message,
		Code:      "UNAUTHORIZED",
		Timestamp: time.Now(),
		Path:      path,
	}
}

func Forbidden(path string, message string) ErrorResponse {
	return ErrorResponse{
		Error:     "Forbidden",
		Message:   message,
		Code:      "FORBIDDEN",
		Timestamp: time.Now(),
		Path:      path,
	}
}

func Conflict(path string, message string) ErrorResponse {
	return ErrorResponse{
		Error:     "Conflict",
		Message:   message,
		Code:      "CONFLICT",
		Timestamp: time.Now(),
		Path:      path,
	}
}

func BadRequest(path string, message string) ErrorResponse {
	return ErrorResponse{
		Error:     "Bad Request",
		Message:   message,
		Code:      "BAD_REQUEST",
		Timestamp: time.Now(),
		Path:      path,
	}
}

func ValidationFailed(path string, message string, validationErrors map[string]string) ErrorResponse {
	return ErrorResponse{
		Error:            "Validation Error",
		Message:          message,
		Code:             "VALIDATION_ERROR",
		Timestamp:        time.Now(),
		Path:             path,
		ValidationErrors: validationErrors,
	}
}

func Internal(path string, message string) ErrorResponse {
	return ErrorResponse{
		Error:     "Internal Server Error",
		Message:   message,
		Code:      "INTERNAL_ERROR",
		Timestamp: time.Now(),
		Path:      path,
	}
}
