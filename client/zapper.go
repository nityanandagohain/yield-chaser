package client

import (
	"log"
	"math/rand"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func GetAPY(client *resty.Client, gaugeVector *prometheus.GaugeVec) {
	// result := &ZapperResponse{}
	// _, err := client.R().
	// 	EnableTrace().
	// 	SetResult(result).
	// 	Get("https://api.zapper.fi/v1/pool-stats/1inch?network=ethereum&api_key=96e0cc51-a62e-42ca-acee-910ea7d2a241")
	// if err != nil {
	// 	log.Fatalf("Failed to create resty client : %v", err)
	// }
	// calculate APY

	// for now creating dummy data only for 1INCH pool

	min := 10
	max := 30
	val1 := rand.Intn(max-min) + min
	val2 := rand.Intn(max-min) + min
	log.Println(val1)
	gaugeVector.WithLabelValues("inch").Set(float64(val1))
	gaugeVector.WithLabelValues("eth").Set(float64(val2))
}
