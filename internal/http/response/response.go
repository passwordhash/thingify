package response

import "github.com/gofiber/fiber/v2"

type ErrorCode string

const (
	ErrorCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrorCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrorCodeConflict     ErrorCode = "CONFLICT"
	ErrorCodeBadRequest   ErrorCode = "BAD_REQUEST"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

func BadRequest(c *fiber.Ctx, err error, message string) error {
	return sendError(c, fiber.StatusBadRequest, ErrorCodeBadRequest, message, err.Error())
}

func Unauthorized(c *fiber.Ctx, err error, message string) error {
	return sendError(c, fiber.StatusUnauthorized, ErrorCodeUnauthorized, message, err.Error())
}

func Internal(c *fiber.Ctx, err error, message string) error {
	return sendError(c, fiber.StatusInternalServerError, ErrorCodeInternal, message)
}
func sendError(c *fiber.Ctx, status int, code ErrorCode, message string, details ...any) error {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		response.Details = details[0]
	}

	return c.Status(status).JSON(response)
}
