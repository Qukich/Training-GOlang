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

type ResponseSber struct {
	Valute Valute `json:"valute"`
}

type Valute struct {
	USD USD `json:"USD"`
}

type USD struct {
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

func (a *SAdapter) WriteRateToDatabase() error {
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
	var binBuf bytes.Buffer
	var arr [3]byte

	copy(arr[:], "USD")
	d := DepartureSber{Name: arr, Sell: math.Round(result.Valute.USD.Value*10) / 10, Time: epochNow}
	binary.Write(&binBuf, binary.BigEndian, d)
	utils.WriteNextBytes(file.Name(), binBuf.Bytes())
	log.Printf("New rate time %d: sell: %f\n", d.Time, d.Sell)
	binBuf.Reset()

	return nil
}

func (a *SAdapter) CloseDB() error {
	return nil
}
