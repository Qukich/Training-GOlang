package route

import (
	"awesomeProject2/adapter"
	"github.com/gofiber/fiber/v2"
)

func GetRate(adapters []adapter.Adapter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticker := c.Params("value")
		bankCode := c.Params("bank")

		var bank adapter.Adapter
		for _, a := range adapters {
			if a.GetCode() == bankCode {
				bank = a
				break
			}
		}

		if bank == nil {
			return c.JSON(fiber.Map{
				"error": "Bank not found",
			})
		}

		rate, err := bank.GetRateFromFile(ticker)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"rate": rate.Sell,
			"timestamp": rate.Time,
		})
	}
}
