package config

type Config struct {
	Application       ApplicationConfig `mapstructure:"app"`
	StockDatabase     DatabaseConfig    `mapstructure:"stock_db"`
	AppDatabase       DatabaseConfig    `mapstructure:"app_db"`
	AnalyticDatabase  DatabaseConfig    `mapstructure:"analytic_db"`
	TimescaleDatabase TimescaleConfig   `mapstructure:"timescale_db"`
	PostgresDatabase  PostgresConfig    `mapstructure:"postgres_db"`
	CryptoDatabase    PostgresConfig    `mapstructure:"crypto_db"`
	Firebase          FirebaseConfig    `mapstructure:"firebase"`
	PageShow          PageShowConfig    `mapstructure:"page_show"`
	Jwt               JwtConfig         `mapstructure:"jwt"`
	Telegram          TelegramConfig    `mapstructure:"telegram"`
	Privy             PrivyConfig       `mapstructure:"privy"`
	CryptoTradingBot  CryptoTradingBotConfig `mapstructure:"crypto_trading_bot"`
}

type ApplicationConfig struct {
	Name     string `mapstructure:"name"`
	Port     string `mapstructure:"port"`
	LogLevel int    `mapstructure:"log_level"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type TimescaleConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}
type PostgresConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
	Unix     string `mapstructure:"unix"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

type FirebaseConfig struct {
	Credential string `mapstructure:"credential"`
}

type PageShowConfig struct {
	CryptoLitePage     bool `mapstructure:"cypto_lite_page"`
	TwitterPage        bool `mapstructure:"twitter_page"`
	SentimentPage      bool `mapstructure:"sentimemt_page"`
	StatsPage          bool `mapstructure:"stats_page"`
	SectorPage         bool `mapstructure:"secttor_page"`
	OverviewMarketPage bool `mapstructure:"overview_market_page"`
	CopyTradePage      bool `mapstructure:"copy_trade_page"`
	GenesisPage        bool `mapstructure:"genesis_page"`
}

type JwtConfig struct {
	Secret string `mapstructure:"secret"`
}

type TelegramConfig struct {
	BotToken string `mapstructure:"bot_token"`
	RunBot   bool   `mapstructure:"run_bot"`
	PadtURL  string `mapstructure:"padt_url"`
}

type PrivyConfig struct {
	AppID                    string `mapstructure:"app_id"`
	AppSecret                string `mapstructure:"app_secret"`
	BASE64_PKCS8_PRIVATE_KEY string `mapstructure:"base64_pkcs8_private_key"`
	EthClient                string `mapstructure:"eth_client"`
	USDCSmartContract        string `mapstructure:"usdc_smart_contract"`
	MaxCopytradeUsers        int    `mapstructure:"max_copytrade_users"`
}

type CryptoTradingBotConfig struct {
	BaseURL string `mapstructure:"baseURL"`
	Token   string `mapstructure:"token"`
}
