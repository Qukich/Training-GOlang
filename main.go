package main

import (
	"awesomeProject2/adapter"
	"awesomeProject2/route"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"
)

type BinHeader struct {
	Version [50]byte
	Year    uint
}

func main() {

	tinkoffAdapter := AdapterFactory("tinkoff")
	if tinkoffAdapter == nil {
		log.Printf("Adapter not found")
	}
	defer tinkoffAdapter.CloseDB()
	go backgroundTask()

	app := fiber.New()

	app.Get("/api/v1/rate/:value", route.GetRate())
	app.Get("/api/v1/history/:value/:bank/t=*", route.GetHistory())

	app.Listen(":3000")
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func AdapterFactory(name string) adapter.Adapter {
	currentTime := time.Now()
	fileName := fmt.Sprintf("%s_%d_%d", name, currentTime.Year(), currentTime.Month())
	currentFilePath := fmt.Sprintf("C:/Users/м/GolandProjects/Training-GOlang/%s_%d_%d", name, currentTime.Year(), currentTime.Month())

	if FileExists(currentFilePath) == false {
		fileDB, err := os.Create(fileName)
		defer fileDB.Close()
		if err != nil {
			log.Fatal(err)
		}
		return &adapter.TAdapter{File: fileDB}
	} else {
		fileDB, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		return &adapter.TAdapter{File: fileDB}
	}

	return nil
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func backgroundTask() {
	tinkoffAdapter := AdapterFactory("tinkoff")
	if tinkoffAdapter == nil {
		log.Printf("Adapter not found")
	}
	sberAdapter := AdapterFactory("sber")
	if tinkoffAdapter == nil {
		log.Printf("Adapter not found")
	}
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		err := tinkoffAdapter.WriteRateToDatabase()
		if err != nil {
			log.Print(err)
		}
		err2 := sberAdapter.WriteRateToDatabase()
		if err2 != nil {
			log.Print(err2)
		}
	}
}
