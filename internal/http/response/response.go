package response

import (
	"github.com/gofiber/fiber/v2"
)

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
	Details any       `json:"details,omitempty,omitzero"`
}

func BadRequest(c *fiber.Ctx, err error, message string) error {
	return sendError(c, fiber.StatusBadRequest, ErrorCodeBadRequest, err, message)
}

func Unauthorized(c *fiber.Ctx, err error, message string) error {
	return sendError(c, fiber.StatusUnauthorized, ErrorCodeUnauthorized, err, message)
}

func Internal(c *fiber.Ctx, err error, message string) error {
	return sendError(c, fiber.StatusInternalServerError, ErrorCodeInternal, err, message)
}

// sendError форматирует и отправляет ошибку [ErrorResponse] в ответе HTTP.
// Параметр err нужен для логирования и может быть nil, если ошибка не требуется.
// Если параметр details не пустой, то он будет добавлен в ответ как поле Details.
// Параметр message используется для описания ошибки и будет отображаться в ответе.
func sendError(c *fiber.Ctx, status int, code ErrorCode, err error, message string, details ...any) error {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	var respDetails []any
	for _, detail := range details {
		if detail != nil {
			respDetails = append(respDetails, details)
		}
	}

	response.Details = respDetails

	c.Status(status).JSON(response)

	return err
}
