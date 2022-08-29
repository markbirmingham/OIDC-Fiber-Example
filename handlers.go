package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Login(auth *Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		state, err := generateRandomState()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())

		}

		session, err := store.Get(c)
		if err != nil {
			panic(err)
		}

		session.Set("state", state)
		if err := session.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Redirect(auth.AuthCodeURL(state), fiber.StatusTemporaryRedirect)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func Logout(c *fiber.Ctx) error {
	logoutUrl, err := url.Parse(os.Getenv("OIDC_PROVIDER_URL") + os.Getenv("OIDC_DOMAIN") + os.Getenv("OIDC_LOGOUT_URL"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())

	}

	scheme := c.Protocol()

	returnTo, err := url.Parse(scheme + "://" + c.Hostname() + "/bye")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	parameters := url.Values{}
	parameters.Add("redirect_uri", returnTo.String())
	parameters.Add("client_id", os.Getenv("OIDC_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	return c.Redirect(logoutUrl.String(), fiber.StatusTemporaryRedirect)
}

func User(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	profile := session.Get("profile")

	return c.Status(fiber.StatusOK).Render("user", profile)
}

func Protected(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	profile := session.Get("profile")

	return c.Status(fiber.StatusOK).Render("protected", profile)
}
