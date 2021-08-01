package client

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var pools map[string]bool

func init() {
	pools = map[string]bool{}
}

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

func GetNewFarms(client *resty.Client) {
	// get the tops pools from the subgraph
	payload := map[string]string{
		"operationName": "topPools",
		"query":         "query topPools {\n  pools(first: 50, orderBy: createdAtTimestamp, orderDirection: desc) {\n    id\n    __typename\n createdAtTimestamp\n }\n}\n",
	}
	res, err := client.R().SetBody(payload).Post("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3")
	if err != nil {
		log.Fatal("Failed to query subgraph", err.Error())
	}
	jsonParsed, err := gabs.ParseJSON(res.Body())
	if err != nil {
		log.Fatal("Failed to unmarshal json fromm subgraph", err.Error())
	}

	newpools := []string{}
	for _, child := range jsonParsed.Path("data.pools").Children() {
		pool := child.Path("id").Data().(string)
		if _, ok := pools[pool]; !ok {
			newpools = append(newpools, pool)
		}
	}

	if len(newpools) == 0 {
		log.Println("No new pools found")
		return
	}

	poolArrStr := ""
	for index, pool := range newpools {
		if index == 3 {
			break
		}
		poolArrStr += "\"" + pool + "\","
	}

	poolArrStr = strings.TrimRight(poolArrStr, ",")

	payload = map[string]string{
		"operationName": "pools",
		"query":         "query pools {\n  pools(\n    where: {id_in: [" + poolArrStr + "]}\n    orderBy: totalValueLockedUSD\n    orderDirection: desc\n  ) {\n    id\n    feeTier\n    liquidity\n    sqrtPrice\n    tick\n    token0 {\n      id\n      symbol\n      name\n      decimals\n      derivedETH\n      __typename\n    }\n    token1 {\n      id\n      symbol\n      name\n      decimals\n      derivedETH\n      __typename\n    }\n    token0Price\n    token1Price\n    volumeUSD\n    txCount\n    totalValueLockedToken0\n    totalValueLockedToken1\n    totalValueLockedUSD\n    __typename\n  }\n}\n",
	}

	res, err = client.R().SetBody(payload).Post("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3")
	if err != nil {
		log.Fatal("Failed to query subgraph", err.Error())
	}

	poolsDataJson, err := gabs.ParseJSON(res.Body())
	if err != nil {
		log.Fatalf("failed to parse pools response body")
	}

	// fmt.Println(poolsDataJson)
	for _, child := range poolsDataJson.Path("data.pools").Children() {
		token0 := child.Path("token0.symbol").Data().(string)
		token1 := child.Path("token1.symbol").Data().(string)
		fmt.Println(token0 + "/" + token1)

		payload = map[string]string{

			"title":            "New Farm Alter",
			"message":          "New Farm" + token0 + "/" + token1,
			"notificationType": "1",
		}
		// send notification
		_, err = client.R().SetBody(payload).Post("http://localhost:8080/notification")
		if err != nil {
			log.Fatal("Failed to query subgraph", err.Error())
		}

		break
	}

}
