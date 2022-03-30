package adapter

import (
	"awesomeProject2/utils"
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
	"unsafe"
)

func GetRate(ticker string, file *os.File) (*Departure, error) {

	//var m Departure2
	//m := Departure2{}

	var arr [8]byte
	copy(arr[:], ticker)

	for i := 1; i < 20; i++ {
		m := Departure2{}
		data := utils.ReadLastBytes(file, int64(unsafe.Sizeof(m)), int64(i))
		buffer := bytes.NewBuffer(data)
		_ = binary.Read(buffer, binary.BigEndian, &m)

		log.Printf("%s,%f,%d\n\n\n", m.Name, m.Sell, m.Time)

		if m.Name == arr {
			return &Departure{Sell: m.Sell, Time: m.Time}, nil
		}
	}

	return nil, errors.New("ticker not found")
}

func GetRateByTimestamp(ticker string, file *os.File, timestamp int64) (*Departure, error) {
	//m := Departure2{}
	return &Departure{}, nil

	//lastTimestampFile := utils.ReadLastBytes(file, int64(unsafe.Sizeof(m)), 1)
	//firstTimestampFile := utils.ReadNextBytes(file, int64(unsafe.Sizeof(m)))
	//bufferFirst := bytes.NewBuffer(firstTimestampFile)
	//bufferLast := bytes.NewBuffer(lastTimestampFile)
	//
	//firstBorder := math.Abs(timestamp - int64(firstTimestampFile))

	//надо понять, к какому из двух значений ближе, и от него искать

	/*for (string(m.Name[:]) != Name) && ((time <= m.Time) && (time >= m.Time+5)) {
		data := utils.ReadNextBytes(file, NumberByte)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.ReadTime failed", err)
		}
		return Departure{Sell: m.Sell, Time: m.Time}
	}*/

	/*for i := 0; i < int(size); i++ {
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
	return Departure{}*/

}
