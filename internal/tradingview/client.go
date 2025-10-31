package tradingview

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/quantsmithapp/datastation-backend/internal/constant"
)

const (
	signInURL = "https://www.tradingview.com/accounts/signin/"
	searchURL = "https://symbol-search.tradingview.com/symbol_search/"
	wsURL     = "wss://prodata.tradingview.com/socket.io/websocket" // Changed URL
	wsTimeout = 5 * time.Second
)

type TradingViewClient struct {
	token        string
	ws           *websocket.Conn
	session      string
	chartSession string
	wsDebug      bool
}

type HistoricalData struct {
	DateTime time.Time `json:"datetime"`
	Symbol   string    `json:"symbol"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   float64   `json:"volume"`
}

type Source2 struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SearchResult struct {
	Symbol       string   `json:"symbol"`
	FullName     string   `json:"full_name,omitempty"`
	Description  string   `json:"description"`
	Exchange     string   `json:"exchange"`
	Type         string   `json:"type"`
	CurrencyCode string   `json:"currency_code"`
	Prefix       string   `json:"prefix"`
	ProviderID   string   `json:"provider_id"`
	Source2      Source2  `json:"source2"`
	SourceID     string   `json:"source_id"`
	TypeSpecs    []string `json:"typespecs"`
}

type SearchResponse struct {
	SymbolsRemaining int            `json:"symbols_remaining"`
	Symbols          []SearchResult `json:"symbols"`
}

// Add new struct for search parameters
type SearchParams struct {
	Text          string `json:"text"`
	Exchange      string `json:"exchange"`
	Start         int    `json:"start"`       // For pagination
	Limit         int    `json:"limit"`       // Items per page
	SearchType    string `json:"search_type"` // e.g. "crypto"
	Lang          string `json:"lang"`
	Domain        string `json:"domain"`
	SortByCountry string `json:"sort_by_country"`
}

func NewTradingViewClient(username, password string) (*TradingViewClient, error) {
	client := &TradingViewClient{
		wsDebug: false,
	}

	// Auth
	token := client.auth(username, password)
	if token == "" {
		client.token = "unauthorized_user_token"
		// log warning about limited access
	} else {
		client.token = token
	}

	client.session = generateSession()
	client.chartSession = generateChartSession()

	return client, nil
}

func (c *TradingViewClient) auth(username, password string) string {
	if username == "" || password == "" {
		return ""
	}

	// Create request body
	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"remember": "on",
	})
	if err != nil {
		return ""
	}

	// Create request
	req, err := http.NewRequest("POST", signInURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return ""
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.tradingview.com")

	// Make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Parse response
	var result struct {
		User struct {
			AuthToken string `json:"auth_token"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	return result.User.AuthToken
}

func (c *TradingViewClient) createConnection() error {
	headers := http.Header{
		"Origin":          []string{"https://www.tradingview.com"}, // Changed origin
		"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		"Accept-Language": []string{"en-US,en;q=0.9"},
	}

	dialer := websocket.Dialer{
		HandshakeTimeout:  wsTimeout,
		EnableCompression: true,
		Proxy:             http.ProxyFromEnvironment,
	}

	ws, _, err := dialer.Dial(wsURL, headers)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %v", err)
	}

	// Set up ping handler
	ws.SetPingHandler(func(appData string) error {
		return ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(10*time.Second))
	})

	c.ws = ws
	return nil
}

func (c *TradingViewClient) GetHistoricalData(
	symbol string,
	exchange string,
	interval string,
	bars int,
	futContract *int,
	extendedSession bool,
) ([]HistoricalData, error) {
	// Format symbol
	formattedSymbol := formatSymbol(symbol, exchange, futContract)

	// Try up to 3 times
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		// Create connection
		if err := c.createConnection(); err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}

		data, err := c.getHistoricalDataWithConnection(formattedSymbol, interval, bars, extendedSession)
		if err != nil {
			c.ws.Close()
			lastErr = err
			time.Sleep(time.Second)
			continue
		}

		return data, nil
	}

	return nil, fmt.Errorf("failed after 1 attempts: %v", lastErr)
}

func (c *TradingViewClient) getHistoricalDataWithConnection(
	symbol string,
	interval string,
	bars int,
	extendedSession bool,
) ([]HistoricalData, error) {
	defer c.ws.Close()

	// Send messages
	messages := []struct {
		function string
		args     []interface{}
	}{
		{"set_auth_token", []interface{}{c.token}},
		{"chart_create_session", []interface{}{c.chartSession, ""}},
		{"resolve_symbol", []interface{}{c.chartSession, "symbol_1", fmt.Sprintf(`={"symbol":"%s","adjustment":"splits","session":"%s"}`,
			symbol, map[bool]string{false: "regular", true: "extended"}[extendedSession])}},
		{"create_series", []interface{}{c.chartSession, "s1", "s1", "symbol_1", interval, bars}},
	}

	// Send messages with delay
	for _, msg := range messages {
		if err := c.sendMessage(msg.function, msg.args); err != nil {
			return nil, fmt.Errorf("failed to send message %s: %v", msg.function, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Read response
	var rawData strings.Builder
	done := make(chan bool)
	errChan := make(chan error)

	go func() {
		for {
			_, message, err := c.ws.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					if rawData.Len() > 0 {
						done <- true
						return
					}
				}
				errChan <- err
				return
			}

			rawData.Write(message)
			rawData.WriteString("\n")

			if strings.Contains(string(message), "series_completed") ||
				strings.Contains(string(message), "timescale_update") {
				done <- true
				return
			}
		}
	}()

	select {
	case <-done:
		return createDataFrame(rawData.String(), symbol)
	case err := <-errChan:
		return nil, fmt.Errorf("websocket error: %v", err)
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("timeout waiting for data")
	}
}

func (c *TradingViewClient) sendMessage(function string, args []interface{}) error {
	msg := createMessage(function, args)
	if c.wsDebug {
		fmt.Println(msg)
	}

	// Set write deadline
	c.ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := c.ws.WriteMessage(websocket.TextMessage, []byte(prependHeader(msg)))
	c.ws.SetWriteDeadline(time.Time{}) // Reset deadline
	return err
}

// Helper functions
func generateSession() string {
	return "qs_" + randomString(12)
}

func generateChartSession() string {
	return "cs_" + randomString(12)
}

func randomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func formatSymbol(symbol, exchange string, contract *int) string {
	if strings.Contains(symbol, ":") {
		return symbol
	}
	if contract != nil {
		return fmt.Sprintf("%s:%s%d!", exchange, symbol, *contract)
	}
	return fmt.Sprintf("%s:%s", exchange, symbol)
}

func createMessage(function string, params []interface{}) string {
	msg := map[string]interface{}{
		"m": function,
		"p": params,
	}
	bytes, _ := json.Marshal(msg)
	return string(bytes)
}

func prependHeader(msg string) string {
	return fmt.Sprintf("~m~%d~m~%s", len(msg), msg)
}

func createDataFrame(rawData, symbol string) ([]HistoricalData, error) {
	re := regexp.MustCompile(`"s":\[(.+?)\}\]`)
	matches := re.FindStringSubmatch(rawData)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no data found")
	}

	dataPoints := strings.Split(matches[1], ",{")
	var result []HistoricalData
	for _, point := range dataPoints {
		fields := strings.FieldsFunc(strings.Trim(point, "{}[]"), func(r rune) bool {
			return r == '[' || r == ':' || r == ',' || r == ']'
		})
		if len(fields) < 9 {
			continue
		}

		// Parse and convert timestamp
		ts, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}
		if ts > 10000000000 { // Convert milliseconds to seconds if needed
			ts = ts / 1000
		}
		datetime := time.Unix(int64(ts), 0).UTC() // Convert to UTC

		// Parse OHLCV values
		open, _ := strconv.ParseFloat(fields[4], 64)
		high, _ := strconv.ParseFloat(fields[5], 64)
		low, _ := strconv.ParseFloat(fields[6], 64)
		close, _ := strconv.ParseFloat(fields[7], 64)
		volume, _ := strconv.ParseFloat(fields[8], 64)

		// Store data with the DateTime in the desired format
		data := HistoricalData{
			DateTime: datetime,
			Symbol:   symbol,
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			Volume:   volume,
		}
		result = append(result, data)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid data points found")
	}

	return result, nil
}

// Modify SearchSymbol to support pagination
func (c *TradingViewClient) SearchSymbol(params SearchParams) (*SearchResponse, error) {
	if params.Limit == 0 {
		params.Limit = 50 // Default limit
	}
	if params.Lang == "" {
		params.Lang = "en"
	}
	if params.Domain == "" {
		params.Domain = "production"
	}
	if params.SearchType == "" {
		params.SearchType = "crypto"
	}

	baseURL := "https://symbol-search.tradingview.com/symbol_search/v3/"

	// Build query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Add query parameters
	q := u.Query()
	q.Set("text", params.Text)
	q.Set("exchange", params.Exchange)
	q.Set("start", strconv.Itoa(params.Start))
	q.Set("limit", strconv.Itoa(params.Limit))
	q.Set("search_type", params.SearchType)
	q.Set("lang", params.Lang)
	q.Set("domain", params.Domain)
	q.Set("sort_by_country", params.SortByCountry)
	q.Set("hl", "1")
	u.RawQuery = q.Encode()

	// Create request with headers
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.tradingview.com")
	req.Header.Set("Referer", "https://www.tradingview.com/")

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Read and clean response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Clean HTML tags from response
	cleanBody := strings.ReplaceAll(string(body), "<em>", "")
	cleanBody = strings.ReplaceAll(cleanBody, "</em>", "")

	// Parse the response
	var searchResponse SearchResponse
	if err := json.Unmarshal([]byte(cleanBody), &searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &searchResponse, nil
}

// Add this new method that uses REST API instead of WebSocket
func (c *TradingViewClient) GetHistoricalDataREST(
	symbol string,
	exchange string,
	interval constant.TradingViewInterval,
	bars int,
) ([]HistoricalData, error) {
	// Format the symbol
	formattedSymbol := formatSymbol(symbol, exchange, nil)

	// Add debug logging
	fmt.Printf("Getting historical data for symbol: %s, interval: %s, bars: %d\n",
		formattedSymbol, interval, bars)

	// Build the URL
	baseURL := "https://tradingview.com/chart/data-server/"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Add query parameters
	q := u.Query()
	q.Set("symbol", formattedSymbol)
	q.Set("resolution", string(interval))
	q.Set("limit", strconv.Itoa(bars))
	u.RawQuery = q.Encode()

	// Debug print the full URL
	fmt.Printf("Request URL: %s\n", u.String())

	// Create request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://www.tradingview.com")
	req.Header.Set("Referer", "https://www.tradingview.com/")

	// Make the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// After getting response, add debug logging
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Debug print the response
	fmt.Printf("Response body: %s\n", string(body))

	// Parse the response
	var data struct {
		T []int64   `json:"t"` // timestamps
		O []float64 `json:"o"` // open
		H []float64 `json:"h"` // high
		L []float64 `json:"l"` // low
		C []float64 `json:"c"` // close
		V []float64 `json:"v"` // volume
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v, body: %s", err, string(body))
	}

	// Debug print the parsed data length
	fmt.Printf("Received data points: %d\n", len(data.T))

	// Convert to HistoricalData
	var result []HistoricalData
	for i := range data.T {
		timestamp := data.T[i]
		datetime := time.Unix(timestamp, 0)

		fmt.Printf("Processing data point %d: timestamp=%d, datetime=%v\n",
			i, timestamp, datetime.Format(time.RFC3339))

		result = append(result, HistoricalData{
			DateTime: datetime,
			Symbol:   formattedSymbol,
			Open:     data.O[i],
			High:     data.H[i],
			Low:      data.L[i],
			Close:    data.C[i],
			Volume:   data.V[i],
		})
	}

	fmt.Printf("Processed %d data points successfully\n", len(result))
	return result, nil
}
