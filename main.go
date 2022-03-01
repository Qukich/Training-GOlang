package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Rates Rates `json:"rates"`
}

type Rates []struct {
	Category     string   `json:"category"`
	FromCurrency Currency `json:"fromCurrency"`
	ToCurrency   Currency `json:"toCurrency"`
	Sell         float64  `json:"sell"`
}

type Currency struct {
	Name string `json:"name"`
}

type Departure struct {
	Sell float64 `json:"sell"`
	Time int64   `json:"time"`
}

type Departure2 struct {
	Name [3]byte `json:"name"` //USD = 0, EUR = 1
	Sell float64 `json:"sell"`
	Time int64   `json:"time"`
}

var result Response
var epochNow int64
var arr [3]byte

func main() {

	app := fiber.New()
	res, err := http.Get("https://api.tinkoff.ru/v1/currency_rates")
	if err != nil {
		log.Fatal(err)
	}
	go backgroundTask(res)
	time.Sleep(2 * time.Second)

	app.Get("/api/v1/rate/:value", func(c *fiber.Ctxы) error {
		mass1 := ReadBinary("USD")
		jsonMass, err := json.Marshal(mass1)
		if err != nil {
			log.Fatal(err)
		}
		return c.Send(jsonMass)
	})

	app.Get("/api/v1/history/", func(c *fiber.Ctx) error {
		return c.Send([]byte("В разработке"))
	})

	app.Listen(":3000")
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func backgroundTask(res *http.Response) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &result)
	epochNow = time.Now().Unix()

	file, err := os.Create("test2.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	var bin_buf bytes.Buffer

	for _, p := range result.Payload.Rates {
		if (p.Category == "C2CTransfers") && ((p.FromCurrency.Name == "USD") || (p.FromCurrency.Name == "EUR")) && (p.ToCurrency.Name == "RUB") {
			copy(arr[:], p.FromCurrency.Name)
			binary.Write(&bin_buf, binary.BigEndian, Departure2{Name: arr, Sell: math.Round(p.Sell*10) / 10, Time: epochNow})
			writeNextBytes(file, bin_buf.Bytes())
		}
	}

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(body, &result)
		epochNow = time.Now().Unix()

		for _, p := range result.Payload.Rates {
			if (p.Category == "C2CTransfers") && ((p.FromCurrency.Name == "USD") || (p.FromCurrency.Name == "EUR")) && (p.ToCurrency.Name == "RUB") {
				copy(arr[:], p.FromCurrency.Name)
				binary.Write(&bin_buf, binary.BigEndian, Departure2{Name: arr, Sell: math.Round(p.Sell*10) / 10, Time: epochNow})
				writeNextBytes(file, bin_buf.Bytes())
			}
		}
	}
}

func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func ReadBinary(Name string) (date Departure2) {
	file, err := os.Open("test2.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := Departure2{}

	for i := 0; i < 3; i++ {
		data := readNextBytes(file, 19)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		fmt.Println(m)
	}

	return m
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
