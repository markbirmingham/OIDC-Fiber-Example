package main

import (
	"encoding/gob"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

var (
	store = session.New()
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := NewAuthenticator()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	engine := html.New("./template", ".html")

	router := fiber.New(fiber.Config{
		Views: engine,
	})

	router.Use(logger.New())

	// Use gob register to store custom types in our cookies
	gob.Register(map[string]interface{}{})

	router.Static("/public", "./static")

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Render("home", nil)
	})

	router.Get("/login", Login(auth))
	router.Get("/callback", Callback(auth))
	router.Get("/user", isAuthenticated, User)
	router.Get("/protected", isAuthenticated, Protected)
	router.Get("/logout", Logout)
	router.Get("/bye", func(c *fiber.Ctx) error {
		return c.Render("logout", nil)
	})

	log.Fatal(router.Listen(":9000"))
}

func isAuthenticated(c *fiber.Ctx) error {

	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	profile := session.Get("profile")

	if profile == nil {
		return c.Redirect("/", fiber.StatusSeeOther)
	}
	return c.Next()
}
