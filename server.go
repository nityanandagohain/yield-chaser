package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nityanandagohain/yield-chaser/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
			client.GetAPY(restyClient, gaugeVector)
			time.Sleep(time.Second * 2)
		}
	}
}

func main() {
	log.Println("Starting application...")

	// register metric
	prometheus.MustRegister(gaugeVector)
	// endpoint for exposing metrics
	http.Handle("/metrics", promhttp.Handler())

	shutdownChannel := make(chan bool)
	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)
	go monitor(shutdownChannel, waitGroup, "zapper")

	http.ListenAndServe(":8090", nil)
	shutdownChannel <- true
	log.Println("Received quit. Sending shutdown and waiting on goroutines...")

	waitGroup.Wait()
	log.Println("Done.")
}
