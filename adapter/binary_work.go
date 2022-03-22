package adapter

import (
	"awesomeProject2/utils"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unsafe"
)

func GetRate(ticker string, file *os.File) (*Departure, error) {

	//var m Departure2
	//m := Departure2{}

	var arr [3]byte
	copy(arr[:], ticker)

	for i := 1; i < 20; i++ {
		m := Departure2{}
		data := utils.ReadLastBytes(file, /*int64(unsafe.Sizeof(m))*/19, int64(i))
		buffer := bytes.NewBuffer(data)
		_ = binary.Read(buffer, binary.BigEndian, &m)

		log.Printf("%s,%f,%d\n\n\n", m.Name, m.Sell, m.Time)

		if m.Name == arr {
			return &Departure{Sell: m.Sell, Time: m.Time}, nil
		}
	}

	return nil, errors.New("ticker not found")
}

func GetRateByTimestamp(ticker string, Time string, NameBank string, NumberByte int) Departure {
	currentTime := time.Now()
	FileName := fmt.Sprintf("./Binary-course/%s_%d_%d.bin", NameBank, currentTime.Month(), currentTime.Year())
	file, err := os.Open(FileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	m := Departure2{}
	stat, err := os.Stat(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	size := stat.Size()
	timestamp, err := strconv.ParseInt(Time, 10, 64)
	if err != nil {
		panic(err)
	}

	/*for (string(m.Name[:]) != Name) && ((time <= m.Time) && (time >= m.Time+5)) {
		data := utils.ReadNextBytes(file, NumberByte)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.ReadTime failed", err)
		}
		return Departure{Sell: m.Sell, Time: m.Time}
	}*/

	for i := 0; i < int(size); i++ {
		data := utils.ReadNextBytes(file, int64(unsafe.Sizeof(m)))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.ReadTime failed", err)
		}
		if (string(m.Name[:]) == ticker) && ((timestamp >= m.Time) && (timestamp <= m.Time+5)) {
			return Departure{Sell: m.Sell, Time: m.Time}
		}
	}
	return Departure{}
}
