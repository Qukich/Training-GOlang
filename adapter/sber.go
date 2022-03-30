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

var rates map[string]float64

type ResponseSber struct {
	Valute map[string]SberTickerInfo
}

type SberTickerInfo struct {
	Value float64 `json:"value"`
}

type SAdapter struct {
	File *os.File
}

type DepartureSber struct {
	Name [3]byte `json:"name"`
	Sell float64 `json:"sell"`
	Time int64   `json:"time"`
}

func init() {
	rates = make(map[string]float64)
}

func (a *SAdapter) GetCode() string {
	return "sber"
}

func (a *SAdapter) WriteRateToFile() error {
	file := a.File

	res, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var result ResponseSber
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	epochNow := time.Now().Unix()

	for ticker, obj := range result.Valute {
		if utils.StringInArray(ticker, []string{"USD", "AMD"}) {
			var binBuf bytes.Buffer
			var arr [3]byte
			copy(arr[:], ticker)
			sell := math.Round(obj.Value * 10) / 10
			needWriteToDatabase := true

			if lastSell, ok := rates[ticker]
			ok {
				if lastSell == sell {
					needWriteToDatabase = false
				}
			} else {
				rates[ticker] = sell
			}

			if needWriteToDatabase {
				tempDeparture := DepartureSber{Name: arr, Sell: sell, Time: epochNow}
				binary.Write(&binBuf, binary.BigEndian, tempDeparture)
				utils.WriteNextBytes(file.Name(), binBuf.Bytes())
				log.Printf("New rate [sber] %s --- time %d: sell: %f\n", arr, tempDeparture.Time, tempDeparture.Sell)
				binBuf.Reset()
			} else {
				log.Printf("test")
			}
		}
	}

	return nil
}

func (a *SAdapter) CloseFile() error {
	return nil
}

func (a *SAdapter) GetRateFromFile(ticker string) (*Departure, error) {
	return GetRate(ticker, a.File)
}

func (a *SAdapter) GetRateByTimestampFromFile(ticker string, timestamp int64) (*Departure, error) {
	return GetRateByTimestamp(ticker, a.File, timestamp)
}
