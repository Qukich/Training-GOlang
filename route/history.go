package route

import (
	"awesomeProject2/adapter"
	"github.com/gofiber/fiber/v2"
)

func GetHistory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticker := c.Params("value")
		timestamp := c.FormValue("t")
		bankCode := c.Params("bank")

		rate := adapter.GetRateByTimestamp(ticker, timestamp, bankCode, 19)

		return c.JSON(fiber.Map{
			"rate": rate.Sell,
			"timestamp": rate.Time,
		})
	}
}
