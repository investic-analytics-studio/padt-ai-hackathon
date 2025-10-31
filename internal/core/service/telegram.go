package service

import (
	"context"
	"fmt"
	"log"

	"encoding/base64"
	"strings"

	tgbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/quantsmithapp/datastation-backend/config"
)

type telegramService struct {
	bot *tgbotapi.Bot
}

func NewTelegramService(token string) (*telegramService, error) {
	// log.Printf("Initializing Telegram bot with token: %s...", token[:10]+"...")
	bot, err := tgbotapi.New(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	service := &telegramService{
		bot: bot,
	}

	// Register handlers
	log.Println("Registering command handlers...")
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText, "/start", tgbotapi.MatchTypeExact, service.handleStart)
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText, "/get_id", tgbotapi.MatchTypeExact, service.handleGetID)
	log.Println("Command handlers registered successfully")

	return service, nil
}

func (s *telegramService) Start(ctx context.Context) error {
	log.Println("Starting Telegram bot...")
	s.bot.Start(ctx)
	log.Println("Bot is running and listening for messages...")
	return nil
}

func (s *telegramService) handleStart(ctx context.Context, b *tgbotapi.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.From.Username
	log.Printf("Received /start command from user %s (ID: %d) in chat %d", username, update.Message.From.ID, chatID)

	message := "Hello! Ready to receive signal\n" +
		"Use /get_id to view your Chat ID"

	log.Printf("Sending welcome message to chat %d", chatID)
	_, err := b.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})
	if err != nil {
		log.Printf("Error sending welcome message: %v", err)
	}
}

func (s *telegramService) handleGetID(ctx context.Context, b *tgbotapi.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	username := update.Message.From.Username
	log.Printf("Received /get_id command from user %s (ID: %d) in chat %d", username, userID, chatID)

	chatIDStr := fmt.Sprintf("%d", chatID)
	chatIDBase64 := base64.StdEncoding.EncodeToString([]byte(chatIDStr))
	// Remove any trailing '=' padding for URL friendliness
	chatIDBase64 = strings.TrimRight(chatIDBase64, "=")
	padtURL := config.GetConfig().Telegram.PadtURL
	url := fmt.Sprintf("%stelegram-alert?ChatID=%s", padtURL, chatIDBase64)

	messageEnterWebsite := "🚨 To start receiving PADT.AI signals, kindly press the button below and ensure a successful login on the website."
	messageEnterChatID := fmt.Sprintf("🚨 Chat ID: `%d`", chatID)
	messageOpenWebsite := "Alternative Way: Open this website and enter your Chat ID to start receiving PADT.AI signals."

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "✅ Confirm",
					URL:  url,
				},
			},
		},
	}
	log.Printf("Sending first message to chat %d", chatID)
	_, err := b.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID:      chatID,
		Text:        messageEnterWebsite,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Printf("Error sending first message: %v", err)
	}

	log.Printf("Sending second message to chat %d", chatID)
	_, err = b.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID:    chatID,
		Text:      messageEnterChatID,
		ParseMode: "Markdown",
	})
	if err != nil {
		log.Printf("Error sending second message: %v", err)
	}
	urlOpenWebsite := fmt.Sprintf("%stelegram-alert", padtURL)
	keyboardOpenWebsite := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "🔗 Open Website",
					URL:  urlOpenWebsite,
				},
			},
		},
	}
	_, err = b.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID:      chatID,
		Text:        messageOpenWebsite,
		ReplyMarkup: keyboardOpenWebsite,
	})
	if err != nil {
		log.Printf("Error sending second message: %v", err)
	}
}

func (s *telegramService) SendMessage(ctx context.Context, text string, chatID int64) error {
	log.Printf("Sending message to chat %d: %s", chatID, text)
	_, err := s.bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
