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
		m := DepartureTinkoff{}
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
	firstStruct := DepartureTinkoff{}
	secondStruct := DepartureTinkoff{}
	//return &Departure{}, nil

	lastTimestampFile := utils.ReadLastBytes(file, int64(unsafe.Sizeof(firstStruct)), 1)
	firstTimestampFile := utils.ReadNextBytes(file, int64(unsafe.Sizeof(secondStruct)))
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

		//return &Departure{}, nil
	} else if secondBorder > firstBorder {
		log.Printf("secondBorder > firstBorder")
		i := 1
		last := DepartureTinkoff{}
		var current DepartureTinkoff

		for {
			tmp := utils.ReadLastBytes(file, int64(unsafe.Sizeof(current)), int64(i))
			bufferTmp := bytes.NewBuffer(tmp)
			err := binary.Read(bufferTmp, binary.BigEndian, &current)
			if err != nil {
				log.Println(err)
				break
			}
			//if last == nil {
			//	last = &current
			//}
			log.Printf("%+v", current)
			//
			if current.Time == 0 {
				break
			}
			//if last.Time > 0 {
				log.Printf("%+v\n%+v", current, last)

				if (current.Time - timestamp) > (last.Time - timestamp) {
					return &Departure{
						Sell: current.Sell,
						Time: current.Time,
					}, nil
				}
			//}
			//

			last = current
			i++
		}
		//return &Departure{}, nil
	}


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
	return &Departure{}, nil
}
