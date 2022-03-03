package main

import (
	"awesomeProject2/adapter"
	"awesomeProject2/route"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"strconv"
	"time"
)

type BinHeader struct {
	Version [50]byte
	Year uint
}

//var epochNow int64
//var count uint16 = 0

func main() {

	tinkoffAdapter := AdapterFactory("tinkoff")
	if tinkoffAdapter == nil {
		log.Printf("Adapter not found")
	}
	defer tinkoffAdapter.CloseDB()

	app := fiber.New()

	go backgroundTask()

	//time.Sleep(2 * time.Second)

	app.Get("/api/v1/history/:value/:bank/t=*", func(c *fiber.Ctx) error {
		timeStamp, _ := strconv.ParseInt(c.Params("*"), 0, 64)
		mass2 := ReadBinaryTime(c.Params("value"), timeStamp)
		jsonMass2, err := json.Marshal(mass2)
		if err != nil {
			log.Fatal(err)
		}
		return c.Send(jsonMass2)
	})

	app.Get("/api/v1/rate/:value", route.GetRate(tinkoffAdapter))

	//app.Get("/api/v1/rate/:value", func(c *fiber.Ctx) error {
	//	mass1 := ReadBinary(c.Params("value"))
	//	jsonMass, err := json.Marshal(mass1)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	//string(mass1.Name[:])
	//	return c.Send(jsonMass)
	//})

	app.Listen(":3000")
}

func AdapterFactory(name string) adapter.Adapter {
	if name == "tinkoff" {
		fileDB, err := os.Create("test2.bin")
		//defer fileDB.Close()
		if err != nil {
			log.Fatal(err)
		}
		return &adapter.TAdapter{ File: fileDB }
	} else if name == "sber" {
		return &adapter.SAdapter{}
	}
	return nil
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func backgroundTask(file *os.File) {

	//err := tinkoffAdapter.WriteRateToDatabase()

	//err := TinkoffGetRate(file)
	//if err != nil {
	//	log.Print(err)
	//	return
	//}

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		err = TinkoffGetRate(file)
		if err != nil {
			log.Print(err)
		}
	}
}




