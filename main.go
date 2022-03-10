package main

import (
	"awesomeProject2/adapter"
	"awesomeProject2/route"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strconv"
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

	app.Get("/api/v1/history/:value/:bank/t=*", func(c *fiber.Ctx) error {
		timeStamp, _ := strconv.ParseInt(c.Params("*"), 0, 64)
		mass2 := adapter.ReadBinary("USD")
		fmt.Println(timeStamp)
		jsonMass2, err := json.Marshal(mass2)
		if err != nil {
			log.Fatal(err)
		}
		return c.Send(jsonMass2)
	})

	app.Get("/api/v1/rate/:value", route.GetRate())

	app.Listen(":3000")
}

func AdapterFactory(name string) adapter.Adapter {
	if name == "tinkoff" {
		fileDB, err := os.Create("tinkoff_03_2022.bin")
		defer fileDB.Close()
		if err != nil {
			log.Fatal(err)
		}
		return &adapter.TAdapter{File: fileDB}
	} else if name == "sber" {
		fileDB, err := os.Create("sber_03_2022.bin")
		defer fileDB.Close()
		if err != nil {
			log.Fatal(err)
		}
		return &adapter.SAdapter{File: fileDB}
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
