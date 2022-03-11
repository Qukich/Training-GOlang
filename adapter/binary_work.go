package adapter

import (
	"awesomeProject2/utils"
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func ReadBinary(Name string, NumberByte int) (LastPart Departure) {
	file, err := os.Open("sber_3_2022.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
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
