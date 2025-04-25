package controller

import (
	"backend/internal"
	"backend/internal/infra/util"
	"github.com/gofiber/fiber/v2"
)

type transformFn[T any] func(*T) error

func parseRequestParams[T any](ctx *fiber.Ctx, transforms ...transformFn[T]) (T, error) {
	var params T
	if err := ctx.BodyParser(&params); err != nil {
		return params, internal.ErrorBadRequest.WithMessage(err.Error())
	}
	if len(transforms) > 0 {
		for _, transform := range transforms {
			if err := transform(&params); err != nil {
				return params, internal.ErrorBadRequest.WithMessage(err.Error())
			}
		}
	}
	if err := util.Validator.Struct(&params); err != nil {
		return params, internal.ErrorBadRequest
	}
	return params, nil
}

func parseRequestQuery[T any](ctx *fiber.Ctx) (T, error) {
	var params T
	if err := ctx.QueryParser(&params); err != nil {
		return params, internal.ErrorBadRequest.WithMessage(err.Error())
	}
	if err := util.Validator.Struct(&params); err != nil {
		return params, internal.ErrorBadRequest
	}
	return params, nil
}
