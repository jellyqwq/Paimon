package coronavirus

import (
	"fmt"
	"log"

	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "github.com/jellyqwq/Paimon/requests"
	// "github.com/jellyqwq/Paimon/tools"
)

func Test(t *testing.T) {
	Core, err := Entry()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(Core)
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	var c int = 1
	for k := range Core.ProvinceDetailed {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(k, fmt.Sprintf("virus-%v", k)))
		// 每四个块合并row到keyboard中并重置row
		if c % 4 == 0 {
			keyboard = append(keyboard, row)
			row = nil
			c = 0
		}
		c += 1
	}
	
	log.Println(keyboard)
}
