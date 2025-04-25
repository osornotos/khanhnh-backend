package middleware

import (
	"backend/internal"
	"backend/internal/domain/auth"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Authorization(tokenProvider auth.TokenProvider[auth.Claims]) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := GetRequestToken(ctx)
		if token == "" {
			return internal.ErrorUnauthorized.WithMessage("invalid token")
		}
		claims, err := tokenProvider.Verify(ctx.Context(), token)
		if err != nil {
			return internal.ErrorUnauthorized.WithMessage("invalid token")
		}

		// save user info
		ctx.Locals("user", claims)
		return ctx.Next()

	}
}

func GetRequestToken(c *fiber.Ctx) string {
	token, _ := c.GetReqHeaders()["Authorization"]
	splitToken := strings.Split(token, "Bearer")
	if len(splitToken) != 2 {
		return ""
	}
	return strings.TrimSpace(splitToken[1])
}
