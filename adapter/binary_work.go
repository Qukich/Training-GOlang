package adapter

import (
	"awesomeProject2/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
)

func ReadBinary(Name string, NumberByte int, NameBank string) (LastPart Departure) {
	FileName := fmt.Sprintf("%s_3_2022.bin", NameBank)
	file, err := os.Open(FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	m := Departure2{}
	var arr [3]byte

	for i := 1; i <= 3; i++ {
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

func ReadBinaryTime(Name string, Time string, NameBank string, NumberByte int) Departure {
	FileName := fmt.Sprintf("%s_3_2022.bin", NameBank)
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
	time, err := strconv.ParseInt(Time, 10, 64)
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(size); i++ {
		data := utils.ReadNextBytes(file, NumberByte)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		if (string(m.Name[:]) == Name) && ((time >= m.Time) && (time <= m.Time+5)) {
			return Departure{Sell: m.Sell, Time: m.Time}
			break
		}
	}
	return Departure{}
}
