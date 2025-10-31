package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/api"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/infra"
	"github.com/quantsmithapp/datastation-backend/internal/core/service"
	"github.com/quantsmithapp/datastation-backend/pkg/logger"
	"github.com/quantsmithapp/datastation-backend/pkg/util"

	_ "github.com/quantsmithapp/datastation-backend/docs" // Generated docs
)

// @title Crypto Studio Backend API
// @version 1.0
// @description API documentation for Crypto Studio Backend
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email admintool@investicstudio.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /api/v2/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your-token}" to authenticate
var conf config.ApplicationConfig

func init() {
	config.InitConfig()
	conf = config.GetConfig().Application
}

func init() {
	logger.InitGlobalLogger()
}

func init() {
	// infra.InitCloudStorage()
	infra.InitStockDB()
	infra.InitAnalyticDB()

	// infra.InitTimescaleDB()
	if err := infra.InitPostgresDB(); err != nil {
		logger.Error(err)
	}

	if err := infra.InitCryptoDB(); err != nil {
		logger.Error(err)
	}

	infra.InitFirebaseClient()

	// Initialize Telegram service
	if config.GetConfig().Telegram.BotToken != "" {
		// Check if we should run the bot based on configuration
		if !config.GetConfig().Telegram.RunBot {
			logger.Info("Skipping Telegram bot initialization (run_bot is set to false in config)")
			return
		}

		// Create a lock file
		lockFile := "/tmp/telegram_bot.lock"
		file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			logger.Error(fmt.Errorf("failed to create lock file: %w", err))
			return
		}

		// Try to acquire an exclusive lock
		// err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		// if err != nil {
		// 	logger.Error(fmt.Errorf("another instance is already running: %w", err))
		// 	file.Close()
		// 	return
		// }

		// Keep the file open and locked while the program runs
		go func() {
			defer file.Close()
			// defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

			telegramService, err := service.NewTelegramService(config.GetConfig().Telegram.BotToken)
			if err != nil {
				logger.Error(fmt.Errorf("failed to initialize Telegram service: %w", err))
				return
			}

			// Start the Telegram bot
			if err := telegramService.Start(context.Background()); err != nil {
				logger.Error(fmt.Errorf("failed to start Telegram bot: %w", err))
			}
		}()
		// TODO: add telegram bot to send message to telegram group
	}
}

func main() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 20 * time.Second,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	// Log every request
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		logger.Infof("[Request] %s %s | %d | %s",
			c.Method(), c.OriginalURL(), c.Response().StatusCode(), duration)

		return err
	})
	app.Get("/swagger/*", swagger.HandlerDefault)
	api.InitAPI(app)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	serverShutdown := make(chan struct{})

	go func() {
		<-c
		logger.Info("🛑 Received shutdown signal, gracefully stopping...")

		// Graceful shutdown within 10 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown fiber with context
		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.Error(fmt.Errorf("error during shutdown: %w", err))
		}

		// Additional cleanup for Cloud Run
		if config.GetConfig().Telegram.RunBot {
			// Remove lock file if exists
			lockFile := "/tmp/telegram_bot.lock"
			if err := os.Remove(lockFile); err != nil && !os.IsNotExist(err) {
				logger.Error(fmt.Errorf("error removing lock file: %w", err))
			}
		}

		serverShutdown <- struct{}{}
	}()

	addr := getAddress()
	logger.Infof("🚀 %v started at %v", conf.Name, addr)

	if err := app.Listen(addr); err != nil {
		logger.Fatal(fmt.Errorf("fiber server error: %w", err))
	}

	<-serverShutdown
	logger.Info("✅ Cleanup done. Server exited.")
}

func getAddress() string {
	addr := ":8080"
	if !util.EmptyString(conf.Port) {
		addr = fmt.Sprintf(":%v", conf.Port)
	}

	return addr
}
