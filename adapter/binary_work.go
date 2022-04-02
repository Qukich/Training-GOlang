package adapter

import (
	"awesomeProject2/utils"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"unsafe"
)

func GetRate(ticker string, file *os.File) (*Departure, error) {
	var arr [8]byte
	copy(arr[:], ticker)

	for i := 1; i < 20; i++ {
		m := DepartureBank{}
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

func GetRateByTimestamp(file *os.File, timestamp int64) (*Departure, error) {
	firstStruct := DepartureBank{}
	secondStruct := DepartureBank{}

	lastTimestampFile := utils.ReadLastBytes(file, int64(unsafe.Sizeof(firstStruct)), 1)
	firstTimestampFile := utils.ReadNextBytes(file, int64(unsafe.Sizeof(secondStruct)), 1)
	bufferFirst := bytes.NewBuffer(firstTimestampFile)
	bufferLast := bytes.NewBuffer(lastTimestampFile)
	err1 := binary.Read(bufferFirst, binary.BigEndian, &firstStruct)
	if err1 != nil {
		log.Println(err1)
	}
	err2 := binary.Read(bufferLast, binary.BigEndian, &secondStruct)
	if err2 != nil {
		log.Println(err2)
	}

	if timestamp < firstStruct.Time || timestamp > secondStruct.Time {
		return nil, fmt.Errorf("rate not found")
	}

	firstBorder := math.Abs(float64(timestamp - firstStruct.Time))
	secondBorder := math.Abs(float64(secondStruct.Time - timestamp))

	log.Printf("%f %f", firstBorder, secondBorder)

	if firstBorder > secondBorder {
		log.Printf("firstBorder > secondBorder")
		i := 1
		last := DepartureBank{}
		var current DepartureBank

		for {
			tmp := utils.ReadNextBytes(file, int64(unsafe.Sizeof(current)), int64(i))
			bufferTmp := bytes.NewBuffer(tmp)
			err := binary.Read(bufferTmp, binary.BigEndian, &current)
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("%+v", current)
			if current.Time == 0 {
				break
			}
			log.Printf("%+v\n%+v", current, last)

			if (current.Time - timestamp) < (last.Time - timestamp) {
				return &Departure{
					Sell: current.Sell,
					Time: current.Time,
				}, nil
			}
			last = current
			i++
		}

	} else if secondBorder > firstBorder {
		log.Printf("secondBorder > firstBorder")
		i := 1
		last := DepartureBank{}
		var current DepartureBank

		for {
			tmp := utils.ReadLastBytes(file, int64(unsafe.Sizeof(current)), int64(i))
			bufferTmp := bytes.NewBuffer(tmp)
			err := binary.Read(bufferTmp, binary.BigEndian, &current)
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("%+v", current)
			if current.Time == 0 {
				break
			}
			log.Printf("%+v\n%+v", current, last)

			if (current.Time - timestamp) > (last.Time - timestamp) {
				return &Departure{
					Sell: current.Sell,
					Time: current.Time,
				}, nil
			}
			last = current
			i++
		}
	}
	return &Departure{}, nil
}
