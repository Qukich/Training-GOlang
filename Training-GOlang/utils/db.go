package utils

import (
	"log"
	"os"
)

func WriteNextBytes(file string, bytes []byte) {
	f, err := os.OpenFile(file, os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err = f.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

func ReadNextBytes(file *os.File, numberByte int64, i int64) []byte {
	bytes := make([]byte, numberByte)
	start := numberByte * i
	_, err := file.ReadAt(bytes, start)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func ReadLastBytes(file *os.File, numberByte int64, i int64) []byte {
	bytes := make([]byte, numberByte)
	stat, err := os.Stat(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	start := stat.Size() - numberByte*i
	_, err = file.ReadAt(bytes, start)
	//log.Printf("Last: %s\n", bytes)
	return bytes
}
