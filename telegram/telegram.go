package telegram

import (
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
	"strconv"
	"virtuals-tracker/virtuals"
)

var _bot *telego.Bot

func GetBot() (*telego.Bot, error) {
	if _bot != nil {
		return _bot, nil
	}

	botToken := os.Getenv("BOT_TOKEN")
	bot, err := telego.NewBot(botToken,
		telego.WithHealthCheck(),
		telego.WithDefaultLogger(false, true))
	if err != nil {
		return nil, err
	}

	_bot = bot
	return _bot, nil
}

func SendNotif(data virtuals.Data) error {
	fmt.Println("Sending notif for", data)

	bot, err := GetBot()
	if err != nil {
		return err
	}

	cid, err := GetChannelID()
	if err != nil {
		return err
	}

	msg := tu.Message(*cid, fmt.Sprintf("%s has reached `%f`", data.Name, data.McapInVirtual))
	_, err = bot.SendMessage(msg)

	if err != nil {
		return err
	}

	return nil
}

func GetChannelID() (*telego.ChatID, error) {
	cid := os.Getenv("CHANNEL_ID")
	if cid == "" {
		return nil, fmt.Errorf("CHANNEL_ID environment variable not set")
	}

	cidf, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return nil, err
	}

	id := tu.ID(cidf)
	return &id, nil
}
