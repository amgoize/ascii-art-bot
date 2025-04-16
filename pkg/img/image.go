package img

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
)

func DownloadPhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (image.Image, error) {
	photo := message.Photo[len(message.Photo)-1]
	fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
	file, err := bot.GetFile(fileConfig)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get("https://api.telegram.org/file/bot" + bot.Token + "/" + file.FilePath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fileExt := strings.ToLower(file.FilePath[len(file.FilePath)-4:])
	if fileExt == ".jpg" || fileExt == "jpeg" {
		return jpeg.Decode(resp.Body)
	} else if fileExt == ".png" {
		return png.Decode(resp.Body)
	}

	return nil, err
}
