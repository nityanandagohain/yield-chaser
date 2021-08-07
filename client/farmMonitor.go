package client

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/go-resty/resty/v2"
)

var uniswapPools map[string]bool
var quickswapPools map[string]bool
var balancerPools map[string]bool

func init() {
	uniswapPools = map[string]bool{}
	quickswapPools = map[string]bool{}
	balancerPools = map[string]bool{}
}

type PoolResponse struct {
	ID          string
	Name        string
	TVL         string
	DailyVolume string
	APY         float32
	Platform    string
	Network     string
}

func GetUniSwapPools(client *resty.Client) *PoolResponse {
	// get the tops pools from the subgraph
	payload := map[string]string{
		"operationName": "topPools",
		"query":         "query topPools {\n  pools(first: 20, orderBy: createdAtTimestamp, orderDirection: desc) {\n    id\n    __typename\n createdAtTimestamp\n }\n}\n",
	}
	res, err := client.R().SetBody(payload).Post("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3")
	if err != nil {
		log.Println("Failed to query subgraph", err.Error())
		return nil
	}
	jsonParsed, err := gabs.ParseJSON(res.Body())
	if err != nil {
		log.Println("Failed to unmarshal json fromm subgraph", err.Error())
		return nil
	}

	newpools := []string{}
	for _, child := range jsonParsed.Path("data.pools").Children() {
		pool := child.Path("id").Data().(string)
		if _, ok := uniswapPools[pool]; !ok {
			newpools = append(newpools, pool)
		}
	}

	if len(newpools) == 0 {
		log.Println("No new pools found")
		return nil
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
		log.Println("Failed to query subgraph", err.Error())
		return nil
	}

	poolsDataJson, err := gabs.ParseJSON(res.Body())
	if err != nil {
		log.Println("failed to parse pools response body")
		return nil
	}
	pool := poolsDataJson.Path("data.pools").Children()[0]
	poolResponse := PoolResponse{
		ID:       pool.Path("id").Data().(string),
		Name:     pool.Path("token0.symbol").Data().(string) + "/" + pool.Path("token1.symbol").Data().(string),
		TVL:      pool.Path("totalValueLockedUSD").Data().(string),
		Platform: "uniswap",
		Network:  "Ethereum mainnet",
	}

	uniswapPools[pool.Path("id").Data().(string)] = true

	return &poolResponse
}

func GetQuickSwapPools(client *resty.Client) *PoolResponse {
	// get the tops pools from the subgraph
	payload := map[string]string{
		"query": "{\n pairs(first: 5, orderBy: createdAtTimestamp, orderDirection: desc) {\n id,\n token0 {\n symbol,\n name\n },\n token1{\n symbol,\n name\n },\n volumeUSD\n totalSupply\n }\n}",
	}
	res, err := client.R().SetBody(payload).Post("https://api.thegraph.com/subgraphs/name/henrydapp/quickswap")
	if err != nil {
		log.Println("Failed to query subgraph", err.Error())
		return nil
	}
	jsonParsed, err := gabs.ParseJSON(res.Body())
	if err != nil {
		log.Println("Failed to unmarshal json fromm subgraph", err.Error())
		return nil
	}

	// newpools := []string{}
	var pool *gabs.Container
	for _, child := range jsonParsed.Path("data.pairs").Children() {
		poolid := child.Path("id").Data().(string)
		if _, ok := quickswapPools[poolid]; !ok {
			pool = child
			break
		}
	}

	if pool == nil {
		log.Println("No new pools found")
		return nil
	}

	poolResponse := PoolResponse{
		ID:       pool.Path("id").Data().(string),
		Name:     pool.Path("token0.symbol").Data().(string) + "/" + pool.Path("token1.symbol").Data().(string),
		TVL:      pool.Path("volumeUSD").Data().(string),
		Platform: "QuickSwap",
		Network:  "polygon mainnet",
	}

	quickswapPools[pool.Path("id").Data().(string)] = true

	return &poolResponse
}

func GetBalancerPools(client *resty.Client) *PoolResponse {
	// get the tops pools from the subgraph
	payload := map[string]string{
		"query": "query { pools (first: 10, orderBy: \"createTime\", orderDirection: \"desc\", where: {totalShares_gt: 0.01, id_not_in: [\"\"], poolType_not: \"Element\", tokensList_contains: []}, skip: 0) { id poolType swapFee tokensList totalLiquidity totalSwapVolume totalSwapFee createTime totalShares owner factory amp tokens { symbol address balance weight } } }",
	}
	res, err := client.R().SetBody(payload).Post("https://api.thegraph.com/subgraphs/name/balancer-labs/balancer-polygon-v2")
	if err != nil {
		log.Println("Failed to query subgraph", err.Error())
		return nil
	}
	jsonParsed, err := gabs.ParseJSON(res.Body())
	if err != nil {
		log.Println("Failed to unmarshal json fromm subgraph", err.Error())
		return nil
	}

	// newpools := []string{}
	var pool *gabs.Container
	for _, child := range jsonParsed.Path("data.pools").Children() {
		poolid := child.Path("id").Data().(string)
		if _, ok := balancerPools[poolid]; !ok {
			pool = child
			break
		}
	}

	if pool == nil {
		log.Println("No new pools found")
		return nil
	}

	tokenName := ""
	for _, child := range pool.Path("tokens").Children() {
		tokenName += "/" + child.Path("symbol").Data().(string)
	}
	tokenName = strings.Trim(tokenName, "/")

	poolResponse := PoolResponse{
		ID:       pool.Path("id").Data().(string),
		Name:     tokenName,
		TVL:      pool.Path("totalLiquidity").Data().(string),
		Platform: "Balancer",
		Network:  "Polygon mainnet",
	}

	balancerPools[pool.Path("id").Data().(string)] = true

	return &poolResponse
}

func SendNotification(client *resty.Client, pool *PoolResponse) {
	payload := map[string]string{
		"title":            "New farm alert!",
		"message":          "Found new farm " + pool.Name + " on " + pool.Platform + ", " + pool.Network + " at " + fmt.Sprint(time.Now().Format(time.RFC850)) + ".",
		"notificationType": "1",
	}
	// send notification
	_, err := client.R().SetBody(payload).Post("https://floating-hollows-80327.herokuapp.com/notification")
	if err != nil {
		log.Println("Failed to query subgraph", err.Error())
	}
}

func GetNewFarms(client *resty.Client) {

	log.Println("Fetching new farms")
	uniswapPool := GetUniSwapPools(client)

	if uniswapPool != nil {
		log.Println("uniswap pool: ", uniswapPool)

		SendNotification(client, uniswapPool)
	}

	quickSwapPool := GetQuickSwapPools(client)
	if quickSwapPool != nil {
		log.Println("qcpool : ", quickSwapPool)
		SendNotification(client, quickSwapPool)
	}

	balancerPool := GetBalancerPools(client)
	if balancerPool != nil {
		log.Println("balancer pool : ", balancerPool)
	}
}
