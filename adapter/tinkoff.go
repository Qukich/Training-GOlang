package adapter

import (
	"awesomeProject2/utils"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

var ratesTinkoff map[string]float64

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

type DepartureTinkoff struct {
	Sell float64 `json:"sell"`
	Time int64   `json:"time"`
	Name [8]byte `json:"name"`
}

type TAdapter struct {
	File *os.File
}

func init() {
	ratesTinkoff = make(map[string]float64)
}

func (a *TAdapter) GetCode() string {
	return "tinkoff"
}

func (a *TAdapter) WriteRateToFile() error {
	file := a.File

	res, err := http.Get("https://api.tinkoff.ru/v1/currency_rates")
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	epochNow := time.Now().Unix()

	for _, obj := range result.Payload.Rates {
		if utils.StringInArray(obj.FromCurrency.Name, []string{"USD", "EUR", "GBP"}) && (obj.Category == "C2CTransfers") && (obj.ToCurrency.Name == "RUB") {
			var binBuf bytes.Buffer
			var arr [8]byte
			copy(arr[:], obj.FromCurrency.Name)

			sell := math.Round(obj.Sell*10) / 10
			needWriteToDatabase := true

			if lastsell, ok := ratesTinkoff[obj.FromCurrency.Name]; ok {
				if lastsell == sell {
					needWriteToDatabase = false
				}
			} else {
				ratesTinkoff[obj.FromCurrency.Name] = sell
			}

			if needWriteToDatabase {
				tempDeparture := DepartureTinkoff{Name: arr, Sell: sell, Time: epochNow}
				binary.Write(&binBuf, binary.BigEndian, tempDeparture)
				utils.WriteNextBytes(file.Name(), binBuf.Bytes())

				log.Printf("New rate [tinkoff] %s --- time %d: sell: %f\n", arr, tempDeparture.Time, tempDeparture.Sell)
				binBuf.Reset()
			} else {
				log.Printf("The course is already in the file Tinkoff")
			}
		}
	}
	return nil
}

func (a *TAdapter) CloseFile() error {
	return nil
}

func (a *TAdapter) GetRateFromFile(ticker string) (*Departure, error) {
	return GetRate(ticker, a.File)
}

func (a *TAdapter) GetRateByTimestampFromFile(ticker string, timestamp int64) (*Departure, error) {
	return GetRateByTimestamp(ticker, a.File, timestamp)
}
