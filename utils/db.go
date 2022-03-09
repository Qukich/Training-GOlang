package utils

import (
	"log"
	"os"
)

func WriteNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func ReadNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func ReadLastBytes(file *os.File, number int64, i int64) []byte {
	bytes := make([]byte, number)
	stat, err := os.Stat(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	start := stat.Size() - number*i
	_, err = file.ReadAt(bytes, start)
	return bytes
}
