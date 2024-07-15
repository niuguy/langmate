package api

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/niuguy/langmate/llm"
)

type BotServer struct {
	bot          *tgbotapi.BotAPI
	updates      tgbotapi.UpdatesChannel
	openaiClient *llm.OpenAIClient
}

func NewBotServer(openaiClient *llm.OpenAIClient) (*BotServer, error) {
	TELEBOT_API_TOKEN := os.Getenv("TELEBOT_API_TOKEN")
	if TELEBOT_API_TOKEN == "" {
		return nil, fmt.Errorf("TELEBOT_API_TOKEN environment variable is not set")
	}

	bot, err := tgbotapi.NewBotAPI(TELEBOT_API_TOKEN)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return &BotServer{
		bot:          bot,
		updates:      updates,
		openaiClient: openaiClient,
	}, nil
}

func (s *BotServer) Start() {
	for update := range s.updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			msg.Text = "Welcome to the bot! How can I help you?"
		case "/help":
			msg.Text = "You can send me any message and I will echo it back!"
		default:
			transferredText, err := s.openaiClient.TransferText(update.Message.Text, "en")
			if err != nil {
				log.Printf("Error transferring text: %v", err)
				transferredText = "Sorry, I couldn't process your request."
			}
			msg.Text = transferredText
		}

		s.bot.Send(msg)
	}
}
