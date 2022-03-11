package route

import (
	"awesomeProject2/adapter"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

func GetRate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		mass1 := adapter.ReadBinary(c.Params("value"), 19, c.Params("bank"))
		jsonMass, err := json.Marshal(mass1)
		if err != nil {
			log.Fatal(err)
		}
		return c.Send(jsonMass)
	}
}
