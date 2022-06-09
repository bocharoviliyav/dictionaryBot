package app

import (
	"dictionaryBot/log"
	repo "dictionaryBot/repository"
	"dictionaryBot/repository/bbolt"
	gtranslate "github.com/gilang-as/google-translate"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"os"
	"strconv"
)

var numericKeyboard = tg.NewInlineKeyboardMarkup(
	tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("3", "1"),
		tg.NewInlineKeyboardButtonData("2", "2"),
		tg.NewInlineKeyboardButtonData("1", "3"),
	),
)

func Run() {
	bot, err := tg.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	db, err := repo.InitBolt()
	if err != nil {
		log.Error("panic %s", err.Error())
	}
	storage := bbolt.NewWordsStorage(db)

	commandsRequest := tg.NewSetMyCommands(tg.BotCommand{Command: "dict", Description: "My dictionary"},
		tg.BotCommand{Command: "add", Description: "Add word"},
		tg.BotCommand{Command: "train", Description: "Training"})

	_, err = bot.Request(commandsRequest)
	if err != nil {
		panic(err)
	}

	updateConfig := tg.NewUpdate(0)
	updateConfig.Timeout = 10
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		log.Info("we've got a message %s", update.Message.Text)

		var msg tg.MessageConfig
		if !update.Message.IsCommand() { // ignore any non-command Messages

			value := gtranslate.Translate{
				Text: update.Message.Text,
				From: "en",
				To:   "ru",
			}
			translated, err := gtranslate.Translator(value)
			if err != nil {
				panic(err)
			}

			msg = tg.NewMessage(update.Message.Chat.ID, translated.Text)
			// msg.ReplyMarkup = numericKeyboard
			msg.ReplyToMessageID = update.Message.MessageID

		} else {
			msg = tg.NewMessage(update.Message.Chat.ID, "")

			// Extract the command from the Message.
			switch update.Message.Command() {
			case "dict":
				get, err := storage.Get(update.Message.From.ID, repo.UserWords)
				if err != nil {
					log.Error("can't get any words %s", err.Error())
				}
				msg.Text = get
			case "add":
				storage.Save(update.Message.From.ID, update.Message.Text+strconv.Itoa(rand.Int()), repo.UserWords)
				msg.Text = "Added"
			case "train":
				msg.Text = "Training"
			case "start":
				msg.Text = "Ты куда звонишь?"
			default:
				msg.Text = "I don't know that command"
			}

		}

		if _, err := bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
	}
}
