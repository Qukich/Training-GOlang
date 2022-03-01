package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
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
	Name [3]byte `json:"name"`
	Sell float64 `json:"sell"`
	Time int64   `json:"time"`
}

var result Response
var epochNow int64
var arr [3]byte
var count uint16 = 0

func main() {

	app := fiber.New()
	res, err := http.Get("https://api.tinkoff.ru/v1/currency_rates")
	if err != nil {
		log.Fatal(err)
	}
	go backgroundTask(res)
	time.Sleep(2 * time.Second)
	app.Get("/api/v1/history/:value/t=*", func(c *fiber.Ctx) error {
		timeStamp, _ := strconv.ParseInt(c.Params("*"), 0, 64)
		mass2 := ReadBinaryTime(c.Params("value"), timeStamp)
		jsonMass2, err := json.Marshal(mass2)
		if err != nil {
			log.Fatal(err)
		}
		return c.Send(jsonMass2)
	})

	app.Get("/api/v1/rate/:value", func(c *fiber.Ctx) error {
		mass1 := ReadBinary(c.Params("value"))
		jsonMass, err := json.Marshal(mass1)
		if err != nil {
			log.Fatal(err)
		}
		//string(mass1.Name[:])
		return c.Send(jsonMass)
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
	var binBuf bytes.Buffer

	for _, p := range result.Payload.Rates {
		if (p.Category == "C2CTransfers") && ((p.FromCurrency.Name == "USD") || (p.FromCurrency.Name == "EUR")) && (p.ToCurrency.Name == "RUB") {
			copy(arr[:], p.FromCurrency.Name)
			binary.Write(&binBuf, binary.BigEndian, Departure2{Name: arr, Sell: math.Round(p.Sell*10) / 10, Time: epochNow})
			writeNextBytes(file, binBuf.Bytes())
			//fmt.Println(binBuf)
			binBuf.Reset()
			count++
		}
	}

	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(body, &result)
		epochNow = time.Now().Unix()

		for _, p := range result.Payload.Rates {
			if (p.Category == "C2CTransfers") && ((p.FromCurrency.Name == "USD") || (p.FromCurrency.Name == "EUR")) && (p.ToCurrency.Name == "RUB") {
				copy(arr[:], p.FromCurrency.Name)
				binary.Write(&binBuf, binary.BigEndian, Departure2{Name: arr, Sell: math.Round(p.Sell*10) / 10, Time: epochNow})
				writeNextBytes(file, binBuf.Bytes())
				//fmt.Println(binBuf)
				binBuf.Reset()
				count++
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

func ReadBinary(Name string) (date Departure) {
	file, err := os.Open("test2.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := Departure2{}
	//var nameValute []byte

	for i := count; i > 0; i-- {
		data := readNextBytes(file, 19)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		//fmt.Println(string(m.Name[:]))
		if (string(m.Name[:]) == Name) && (m.Time == epochNow) {
			date = Departure{Sell: m.Sell, Time: m.Time}
			break
		}
	}
	return date
}

func ReadBinaryTime(Name string, Time int64) (date Departure) {
	file, err := os.Open("test2.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := Departure2{}
	//var nameValute []byte

	for i := 0; i < int(count); i++ {
		data := readNextBytes(file, 19)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		//fmt.Println(string(m.Name[:]))
		if (string(m.Name[:]) == Name) && ((Time >= m.Time) && (Time < m.Time+10)) {
			date = Departure{Sell: m.Sell, Time: m.Time}
			break
		}
	}
	return date
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
