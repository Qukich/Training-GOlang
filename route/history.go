package route

import (
	"awesomeProject2/adapter"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

func GetHistory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		timeStamp, _ := strconv.ParseInt(c.Params("*"), 0, 64)
		mass2 := adapter.ReadBinary("USD", 2)
		fmt.Println(timeStamp)
		jsonMass2, err := json.Marshal(mass2)
		if err != nil {
			log.Fatal(err)
		}
		return c.Send(jsonMass2)
	}
}
