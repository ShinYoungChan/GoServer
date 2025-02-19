package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// GET 방식 예제
	app.Get("/fiber", func(c *fiber.Ctx) error {
		return c.SendString("Hello, get fiber!")
	})

	// GET 요청: 동적 경로 (URL 파라미터)
	app.Get("/user/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString("Hello, " + name)
	})

	app.Listen(":8080")
}
