package main

import (
	"awesomeProject2/adapter"
	"awesomeProject2/route"
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
	app := fiber.New()
	go backgroundTask()

	app.Get("/api/v1/rate/:value/:bank", route.GetRate())
	app.Get("/api/v1/history/:value/:bank/t=*", route.GetHistory())

	app.Listen(":3000")
}

func AdapterFactory(name string) adapter.Adapter {
	currentTime := time.Now()
	if name == "tinkoff" {
		fileName := fmt.Sprintf("%s_%d_%d.bin", name, currentTime.Month(), currentTime.Year())
		fileDB, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		return &adapter.TAdapter{File: fileDB}
	} else if name == "sber" {
		fileName := fmt.Sprintf("%s_%d_%d.bin", name, currentTime.Month(), currentTime.Year())
		fileDB, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		return &adapter.SAdapter{File: fileDB}
	}
	return nil
}

func backgroundTask() {
	tinkoffAdapter := AdapterFactory("tinkoff")
	if tinkoffAdapter == nil {
		log.Printf("Adapter not found")
	}
	sberAdapter := AdapterFactory("sber")
	if sberAdapter == nil {
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
