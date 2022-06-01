package app

import (
	"dictionaryBot/log"
	repo "dictionaryBot/repository"
	"dictionaryBot/repository/bbolt"
	gtranslate "github.com/gilang-as/google-translate"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
)

var numericKeyboard = tg.NewInlineKeyboardMarkup(
	tg.NewInlineKeyboardRow(
		tg.NewInlineKeyboardButtonData("Х", "1"),
		tg.NewInlineKeyboardButtonData("У", "2"),
		tg.NewInlineKeyboardButtonData("Й", "3"),
	),
)

func Run() {
	//bot, err := tg.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	bot, err := tg.NewBotAPI("2113758243:AAEwwAD8ws1oRz7Lq_YGz9AnCwXLgOomans")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	db, err := repo.InitBolt()
	if err != nil {
		log.Error("PIZDEC %s", err.Error())
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
		if !update.Message.IsCommand() { // ignore any non-command Messages-

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

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
	}
}

/*2022/01/18 20:00:15 Endpoint: sendMessage, response: {"ok":true,"result":{"message_id":63,"from":{"id":2113758243,"is_bot":true,"first_name":"engvocabulary_bot","username":"eng_vocab_trainer_bot"},"chat":{"id":76982095,"first_name":
"\u0414\u043c\u0438\u0442\u0440\u0438\u0439","username":"Kiryakov","type":"private"},"date":1642521615,"reply_to_message":{"message_id":61,"from":{"id":76982095,"is_bot":false,"first_name":"\u0414\u043c\u0438\u0442\u0440\u0438\u0439
","username":"Kiryakov","language_code":"ru"},"chat":{"id":76982095,"first_name":"\u0414\u043c\u0438\u0442\u0440\u0438\u0439","username":"Kiryakov","type":"private"},"date":1642521615,"text":"1"},"text":"Prepare Uranus"}}
2022/01/18 20:00:17 Endpoint: getUpdates, response: {"ok":true,"result":[{"update_id":203068619,
"message":{"message_id":64,"from":{"id":76982095,"is_bot":false,"first_name":"\u0414\u043c\u0438\u0442\u0440\u0438\u0439","username":"Kiryakov","language_code":"ru"},"chat":{"id":76982095,"first_name":"\u0414\u043c\u0438\u0442\u0440
\u0438\u0439","username":"Kiryakov","type":"private"},"date":1642521617,"text":"1"}}]}
*/
