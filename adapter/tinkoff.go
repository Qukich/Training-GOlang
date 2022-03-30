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
	Sell float64 `json:"sell"`
	Time int64   `json:"time"`
	Name [8]byte `json:"name"`
}

type TAdapter struct {
	File *os.File
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
	var binBuf bytes.Buffer

	for _, p := range result.Payload.Rates {
		//&& (len(p.FromCurrency.Name) == 3)
		if (p.Category == "C2CTransfers") && (p.ToCurrency.Name == "RUB") {
			var arr [8]byte
			copy(arr[:], p.FromCurrency.Name)
			tempDeparture := Departure2{Name: arr, Sell: math.Round(p.Sell*10) / 10, Time: epochNow}
			binary.Write(&binBuf, binary.BigEndian, tempDeparture)
			utils.WriteNextBytes(file.Name(), binBuf.Bytes())

			log.Printf("New rate [tinkoff] %s --- time %d: sell: %f\n", arr, tempDeparture.Time, tempDeparture.Sell)
			binBuf.Reset()
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
