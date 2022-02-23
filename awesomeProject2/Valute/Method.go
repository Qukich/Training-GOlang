package Valute

import "github.com/gofiber/fiber/v2"

func GetValute(c *fiber.Ctx) error {
	return c.Send([]byte("USD"))
}

func GetValuteTime(c *fiber.Ctx) error {
	return c.Send([]byte("USD time"))
}
