package main

import (
	"github.com/gofiber/fiber/v2"
)

func Callback(auth *Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			panic(err)
		}

		if c.Query("state") != session.Get("state") {
			c.Status(fiber.StatusBadRequest).SendString("Invalid state parameter.")
			return nil
		}

		token, err := auth.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Failed to exchange an authorization code for a token.")
		}

		idToken, err := auth.VerifyIDToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to verify ID Token.")
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)

		if err := session.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Redirect("/user", fiber.StatusTemporaryRedirect)
	}
}
