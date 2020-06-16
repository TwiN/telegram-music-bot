package main

import (
	"github.com/TwinProduction/telegram-music-bot/config"
	"github.com/TwinProduction/telegram-music-bot/youtube"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"strings"
	"time"
)

var (
	youtubeService *youtube.Service
	activeTasks    = 0
)

func main() {
	config.Load()
	youtubeService = youtube.NewService(config.Get().MaximumAudioDurationInSeconds)
	bot, err := telebot.NewBot(telebot.Settings{
		Token:     config.Get().TelegramToken,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeMarkdown,
	})
	if err != nil {
		panic(err)
	}
	bot.Handle("/yt", HandleYoutubeCommand(bot))
	bot.Handle("/youtube", HandleYoutubeCommand(bot))
	err = bot.SetCommands([]telebot.Command{
		{
			Text:        "yt",
			Description: "Search for a clip on Youtube and convert it to MP3",
		},
		{
			Text:        "youtube",
			Description: "Search for a clip on Youtube and convert it to MP3",
		},
	})
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Starting telegram bot @%s", bot.Me.Username)
	bot.Start()
	defer bot.Stop()
}

func HandleYoutubeCommand(bot *telebot.Bot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		query := extractQueryFromText(m.Text)
		if len(query) < 5 {
			_, _ = bot.Reply(m, "Query must be at least 5 characters long\ne.g.: `/yt what is love`")
			return
		}
		if activeTasks >= config.Get().MaximumActiveTasks {
			_, _ = bot.Reply(m, "Ask me again in a few seconds, I'm busy.")
			return
		}
		activeTasks++
		defer func() {
			activeTasks--
		}()
		chatName := extractNameFromChat(m.Chat)
		log.Printf("[@%s] User @%s requested \"%s\"", chatName, m.Sender.Username, query)
		statusMessage, _ := bot.Reply(m, "⌛ Give me a moment...")
		media, err := youtubeService.SearchAndDownload(query)
		if err != nil {
			_, _ = bot.Reply(m, "I ran into an error and couldn't complete your request :(")
			return
		}
		statusMessage, _ = bot.Edit(statusMessage, "✅ Successfully downloaded.\n⌛ Uploading file...")
		defer os.Remove(media.FilePath)
		_, err = bot.Send(m.Chat, &telebot.Audio{
			File: telebot.File{
				FileLocal: media.FilePath,
			},
			Duration:  int(media.Duration.Seconds()),
			Title:     media.Title,
			Caption:   media.URL,
			Performer: media.Uploader,
			FileName:  media.FilePath,
		})
		if err != nil {
			log.Printf("[@%s] Ran into an error trying to process request from User @%s for query \"%s\": %s", chatName, m.Sender.Username, query, err.Error())
			statusMessage, _ = bot.Edit(statusMessage, "❌ Ran into an error trying to process your query!")
			return
		}
		statusMessage, _ = bot.Edit(statusMessage, "✅ File uploaded successfully!")
		go func(bot *telebot.Bot, statusMessage *telebot.Message) {
			time.Sleep(5 * time.Second)
			_ = bot.Delete(statusMessage)
		}(bot, statusMessage)
		log.Printf("[@%s] Successfully completed request for user @%s's \"%s\" query", chatName, m.Sender.Username, query)
	}
}

func extractQueryFromText(text string) string {
	var query string
	if strings.HasPrefix(text, "/youtube") {
		query = strings.TrimSpace(strings.TrimPrefix(text, "/youtube"))
	} else if strings.HasPrefix(text, "/yt") {
		query = strings.TrimSpace(strings.TrimPrefix(text, "/yt"))
	}
	query = strings.ReplaceAll(query, "\"", "")
	query = strings.ReplaceAll(query, "'", "")
	query = strings.ReplaceAll(query, "`", "")
	query = strings.ReplaceAll(query, "\\", "")
	return query
}

func extractNameFromChat(chat *telebot.Chat) string {
	identifier := chat.Username
	if len(identifier) == 0 {
		identifier = chat.Title
	}
	return identifier
}
