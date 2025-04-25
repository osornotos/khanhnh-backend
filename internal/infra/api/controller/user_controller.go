package controller

import (
	"backend/internal/application/interface"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type UserController struct {
	userService _interface.UserService
}

func NewUserController(userService _interface.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) Login(ctx *fiber.Ctx) error {
	token, err := u.userService.Login(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"data":   token,
		"status": "ok",
	})
}
