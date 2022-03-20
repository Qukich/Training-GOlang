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
	app := fiber.New()
	go backgroundTask()

	app.Get("/api/v1/rate/:value/:bank", route.GetRate())
	app.Get("/api/v1/history/:value/:bank/t=*", route.GetHistory())

	app.Listen(":3000")
}

func Check_file(path string) bool {
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
	if name == "tinkoff" {
		fileName := fmt.Sprintf("Binary-course/%s_%d_%d.bin", name, currentTime.Month(), currentTime.Year())
		currentFilePath := fmt.Sprintf("C:/Users/Roman/GolandProjects/awesomeProject2/%s", fileName)
		if Check_file(currentFilePath) == false {
			fileDB, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			return &adapter.TAdapter{File: fileDB}
		} else {
			fileDB, err := os.OpenFile(fileName, os.O_CREATE, 0660)
			//fileDB, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			return &adapter.TAdapter{File: fileDB}
		}
	}
	if name == "sber" {
		fileName := fmt.Sprintf("Binary-course/%s_%d_%d.bin", name, currentTime.Month(), currentTime.Year())
		currentFilePath := fmt.Sprintf("C:/Users/Roman/GolandProjects/awesomeProject2/%s", fileName)
		if Check_file(currentFilePath) == false {
			fileDB, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			return &adapter.SAdapter{File: fileDB}
		} else {
			fileDB, err := os.OpenFile(fileName, os.O_CREATE, 0660)
			//fileDB, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			return &adapter.SAdapter{File: fileDB}
		}
	}
	return nil
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
	ticker := time.NewTicker(50 * time.Second)
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
