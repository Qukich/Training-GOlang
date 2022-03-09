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
	Name [3]byte `json:"name"`
	Sell float64 `json:"sell"`
	Time int64   `json:"time"`
}

type TAdapter struct {
	File *os.File
}

func (a *TAdapter) WriteRateToDatabase() error {
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
		if (p.Category == "C2CTransfers") && (len(p.FromCurrency.Name) == 3) && (p.ToCurrency.Name == "RUB") {
			var arr [3]byte
			copy(arr[:], p.FromCurrency.Name)
			d := Departure2{Name: arr, Sell: math.Round(p.Sell*10) / 10, Time: epochNow}
			binary.Write(&binBuf, binary.BigEndian, d)
			utils.WriteNextBytes(file, binBuf.Bytes())

			log.Printf("New rate time %d: sell: %f\n", d.Time, d.Sell)
			binBuf.Reset()
		}
	}

	return nil
}

func ReadBinary(Name string) (LastPart Departure) {
	file, err := os.Open("test2.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	NumberByte := 19
	m := Departure2{}
	var arr [3]byte

	for i := 1; i <= 10; i++ {
		data := utils.ReadLastBytes(file, int64(NumberByte), int64(i))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		copy(arr[:], Name)
		if arr == m.Name {
			break
		}
	}
	LastPart = Departure{Sell: m.Sell, Time: m.Time}

	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return LastPart
}

func ReadBinaryTime(Name string, Time int64) Departure {
	file, err := os.Open("test2.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	NumberByte := 19
	//
	m := Departure2{}
	//var nameValute []byte

	//for i := 0; i < int(count); i++ {
	data := utils.ReadNextBytes(file, NumberByte)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}
	//fmt.Println(string(m.Name[:]))
	if (string(m.Name[:]) == Name) && ((Time >= m.Time) && (Time < m.Time+10)) {
		return Departure{Sell: m.Sell, Time: m.Time}
		//break
	}
	//}
	return Departure{}
	//return date
}
