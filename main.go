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

func main() {
	banks := []string{"tinkoff", "sber"}

	var adapters adapter.Adapters

	for _, bank := range banks {
		a := AdapterFactory(bank)
		if a == nil {
			log.Fatalf("Adapter %s not found", bank)
		}
		adapters = append(adapters, AdapterFactory(bank))
	}

	app := fiber.New()

	go backgroundTask(adapters)
	defer closeDBFiles(adapters)

	app.Get("/api/v1/rate/:value/:bank", route.GetRate(adapters))
	app.Get("/api/v1/history/:value/:bank", route.GetHistory(adapters))

	log.Fatal(app.Listen(":3000"))
}

func closeDBFiles(adapters []adapter.Adapter) {
	for _, a := range adapters {
		_ = a.CloseFile()
	}
}

func getDBFileName(name string) string {
	year, month, _ := time.Now().Date()
	return fmt.Sprintf("./Binary-course/%s_%d_%d.bin", name, month, year)
}

func AdapterFactory(name string) adapter.Adapter {
	if name == "tinkoff" {
		fileDB, err := os.OpenFile(getDBFileName(name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			log.Fatal(err)
		}

		return &adapter.TAdapter{File: fileDB}
	}
	if name == "sber" {
		fileDB, err := os.OpenFile(getDBFileName(name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			log.Fatal(err)
		}

		return &adapter.SAdapter{File: fileDB}
	}
	return nil
}

func backgroundTask(adapters []adapter.Adapter) {
	ticker := time.NewTicker(50 * time.Second)
	for range ticker.C {
		for _, a := range adapters {
			err := a.WriteRateToFile()
			if err != nil {
				log.Print(err)
			}
		}
	}
}
