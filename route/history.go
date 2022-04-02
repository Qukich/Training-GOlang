package route

import (
	"awesomeProject2/adapter"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetHistory(adapters []adapter.Adapter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		timestamp := c.FormValue("t")
		timestampInt64, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": "Bank not found",
			})
		}
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

		rate, err := bank.GetRateByTimestampFromFile(timestampInt64)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"rate":      rate.Sell,
			"timestamp": rate.Time,
		})
	}
}
