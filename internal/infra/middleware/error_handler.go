package middleware

import (
	"backend/internal"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	var appErr internal.ApplicationError

	switch e := err.(type) {
	case *fiber.Error:
		if e.Code >= 500 {
			appErr = internal.ErrorInternalServerError.WithMessage(e.Message).WithError(err)
		} else {
			appErr = internal.ApplicationError{HTTPStatusCode: e.Code, ErrorText: e.Message}
			appErr.ErrorCode = appErr.GetDefaultErrorCode()
		}
	case internal.ApplicationError:
		appErr = e
	default: // should never happen
		appErr = internal.ErrorInternalServerError.WithMessage(err.Error())
	}

	// client response
	return ctx.Status(appErr.HTTPStatusCode).JSON(appErr)
}
