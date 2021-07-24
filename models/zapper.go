package models

type ZapperResponse []struct {
	Address         string  `json:"address"`
	TokenAddress    string  `json:"tokenAddress"`
	ContractAddress string  `json:"contractAddress"`
	Decimals        int     `json:"decimals"`
	Symbol          string  `json:"symbol"`
	Value           string  `json:"value"`
	Label           string  `json:"label"`
	Supply          float64 `json:"supply"`
	Liquidity       float64 `json:"liquidity"`
	PricePerToken   float64 `json:"pricePerToken"`
	Protocol        string  `json:"protocol"`
	ProtocolDisplay string  `json:"protocolDisplay"`
	ProtocolSymbol  string  `json:"protocolSymbol"`
	Tokens          []struct {
		Address  string  `json:"address"`
		Decimals int     `json:"decimals"`
		Symbol   string  `json:"symbol"`
		Price    float64 `json:"price"`
		Reserve  float64 `json:"reserve"`
	} `json:"tokens"`
	Fee       float64 `json:"fee"`
	FeeVolume float64 `json:"feeVolume"`
	Volliq    float64 `json:"volliq"`
	Volume    float64 `json:"volume"`
	DailyROI  float64 `json:"dailyROI"`
	WeeklyROI float64 `json:"weeklyROI"`
	YearlyROI float64 `json:"yearlyROI"`
}
