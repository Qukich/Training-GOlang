package route

import (
	"awesomeProject2/adapter"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

func GetHistory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		mass2 := adapter.ReadBinaryTime(c.Params("value"), c.Params("*"), c.Params("bank"), 19)
		jsonMass2, err := json.Marshal(mass2)
		if err != nil {
			log.Fatal(err)
		}
		return c.Send(jsonMass2)
	}
}
