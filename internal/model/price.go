package model

type Get24HoursPriceBySymbolsReq struct {
	Symbols []string `json:"symbols"` // ["BTC", "ETH"]
}

type ScannerData struct {
	S string        `json:"s"`
	D []interface{} `json:"d"`
}

type ScannerResponse struct {
	TotalCount int           `json:"totalCount"`
	Data       []ScannerData `json:"data"`
}
