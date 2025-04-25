package middleware

import (
	"github.com/gofiber/fiber/v2"
	"khanhnh-backend/internal"
	"khanhnh-backend/internal/domain/auth"
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
