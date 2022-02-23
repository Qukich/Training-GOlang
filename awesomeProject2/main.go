package main

import (
	"awesomeProject2/Valute"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	TrackingId string `json:"trackingId"`
	ResultCode string `json:"resultCode"`
	payload    struct {
		LastUpdate struct {
			Milliseconds int64 `json:"milliseconds"`
		} `json:"lastUpdate"`
	} `json:"payload"`
}

func main() {
	/*app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")*/

	res, err := http.Get("https://api.tinkoff.ru/v1/currency_rates")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	fmt.Println(PrettyPrint(result))

}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/rate", Valute.GetValute)
	app.Get(" /api/v1/history", Valute.GetValuteTime)
}
