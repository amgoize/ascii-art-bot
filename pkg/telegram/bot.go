package telegram

import (
	"ascii-art-server/pkg/ascii"
	"ascii-art-server/pkg/img"
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	userPhoto map[int64]image.Image
}

func NewBot(token string) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Println("Бот запущен")
	return &Bot{
		api:       botAPI,
		userPhoto: make(map[int64]image.Image),
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		// Обработка отправленного фото
		if update.Message != nil && update.Message.Photo != nil {
			log.Println("Получено изображение")

			photo, err := img.DownloadPhoto(b.api, update.Message)
			if err != nil {
				log.Printf("Ошибка скачивания: %v", err)
				continue
			}

			// Сохранение фото
			chatID := update.Message.Chat.ID
			b.userPhoto[chatID] = photo

			// Отправка кнопок для выбора
			msg := tgbotapi.NewMessage(chatID, "Выберите тип ASCII-арта:")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Ч/Б", "ascii_bw"),
					tgbotapi.NewInlineKeyboardButtonData("Цветной", "ascii_color"),
				),
			)
			msg.ReplyMarkup = keyboard
			b.api.Send(msg)
		} else if update.Message != nil && update.Message.Command() == "start" {
			chatID := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatID, "Привет! Я бот для создания ASCII-артов. Отправь мне изображение.")
			b.api.Send(msg)
		} else if update.Message != nil && update.Message.Photo == nil {
			chatID := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, отправьте фото.")
			b.api.Send(msg)
		}

		// Обработка ответа на кнопку
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			data := update.CallbackQuery.Data

			photo, ok := b.userPhoto[chatID]
			if !ok {
				log.Println("Фото не найдено")
				continue
			}

			// Преобразование image.Image в []byte
			var buf bytes.Buffer
			err := jpeg.Encode(&buf, photo, nil)
			if err != nil {
				log.Printf("Ошибка кодирования: %v", err)
				continue
			}

			if data == "ascii_bw" {
				asciiArt := ascii.ConvertToASCIIArt(photo)
				file := tgbotapi.FileBytes{
					Name:  "ascii_art.txt",
					Bytes: []byte(asciiArt),
				}
				doc := tgbotapi.NewDocument(chatID, file)
				doc.Caption = "Вот твой Ч/Б ASCII-арт"
				b.api.Send(doc)
			} else if data == "ascii_color" {
				htmlArt := ascii.ConvertToColorASCIIArt(photo, "html")
				file := tgbotapi.FileBytes{
					Name:  "ascii_art.html",
					Bytes: []byte(htmlArt),
				}
				doc := tgbotapi.NewDocument(chatID, file)
				doc.Caption = "Вот твой цветной ASCII-арт (нужно открыть в браузере)"
				b.api.Send(doc)
			}

			// Удаление кнопки после выбора
			edit := tgbotapi.NewEditMessageReplyMarkup(
				chatID,
				update.CallbackQuery.Message.MessageID,
				tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
			)
			b.api.Send(edit)
		}
	}
}
