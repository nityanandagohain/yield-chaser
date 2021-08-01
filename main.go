package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nityanandagohain/yield-chaser/client"
	"github.com/nityanandagohain/yield-chaser/controllers"
	"github.com/prometheus/client_golang/prometheus"
)

var server = controllers.Server{}

var gaugeVector = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "APY_CHANGE_GAGUE",
		Help: "No of request handled by Ping handler",
	},
	[]string{"APY_CHANGE"},
)

func monitor(shutdownChannel chan bool, waitGroup *sync.WaitGroup, name string) {
	log.Println("Starting work goroutine...")
	defer waitGroup.Done()
	restyClient := resty.New()
	for {
		select {
		case <-shutdownChannel:
			log.Println("Shutting down channel: ", name)
			return
		default:
			log.Println("Feching new data from : ", name)
			client.GetNewFarms(restyClient)
			time.Sleep(time.Hour * 1)
		}
	}
}

func main() {
	log.Println("Starting application...")

	shutdownChannel := make(chan bool)
	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)
	go monitor(shutdownChannel, waitGroup, "zapper")

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8090" //localhost
	}

	server.Initialize(os.Getenv("DATABASE_URL"), ":"+port)
	shutdownChannel <- true
	log.Println("Received quit. Sending shutdown and waiting on goroutines...")

	waitGroup.Wait()
	log.Println("Done.")
}
