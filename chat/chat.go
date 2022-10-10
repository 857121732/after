package chat

import (
	"log"
	"math/rand"

	"github.com/kevwan/chatbot/bot"
	"github.com/kevwan/chatbot/bot/adapters/logic"
	"github.com/kevwan/chatbot/bot/adapters/storage"
)

var global *bot.ChatBot

func InitChat(corpora string, tops int) *bot.ChatBot {
	store, err := storage.NewSeparatedMemoryStorage(corpora)
	if err != nil {
		log.Fatal(err)
	}

	global = &bot.ChatBot{
		LogicAdapter: logic.NewClosestMatch(store, tops),
	}
	return global
}

func GetResponse(text string) string {
	if global == nil {
		return "Emmmmmmmm......"
	}

	res := global.GetResponse(text)
	if len(res) == 0 {
		return "娃没听懂耶。。。在学了在学了！！！"
	}
	return res[rand.Intn(len(res))].Content
}
