package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/vbchekhov/skeleton"
)

var botName = ""
var botPic = ""

func runBot() {

	bot := skeleton.NewBot(config.Token)
	bot.HandleFunc("(.*)", saveMessage).Border(skeleton.Private).Methods(skeleton.Messages)
	bot.Run()
}

func saveMessage(c *skeleton.Context) bool {

	message := c.Update.Message
	post := Post{
		Time: time.Now().Format(time.RFC3339),
	}

	id := ""

	// private chat
	if message.ForwardFromChat == nil {
		id = fmt.Sprintf("%d-%d", message.Chat.ID, c.Update.UpdateID)

		chat, _ := c.BotAPI.GetChat(message.Chat.ChatConfig())

		photoID := chat.Photo.BigFileID
		directURL, _ := c.BotAPI.GetFileDirectURL(photoID)
		f := download(directURL, photoID+".jpeg")

		post.Metadata = Metadata{
			Title:    message.Chat.FirstName + " " + message.Chat.LastName,
			UserName: message.Chat.UserName,
			Avatar:   f.Name(),
			Group:    fmt.Sprintf("https://t.me/%s", message.Chat.UserName),
		}
	}

	//  public chat
	if message.ForwardFromChat != nil {
		id = fmt.Sprintf("%d-%d", message.ForwardFromChat.ID, message.ForwardDate)

		chat, _ := c.BotAPI.GetChat(message.ForwardFromChat.ChatConfig())

		photoID := chat.Photo.BigFileID
		directURL, _ := c.BotAPI.GetFileDirectURL(photoID)
		f := download(directURL, photoID+".jpeg")

		post.Metadata = Metadata{
			Title:    message.ForwardFromChat.Title,
			UserName: message.ForwardFromChat.UserName,
			Avatar:   f.Name(),
			Group:    fmt.Sprintf("https://t.me/%s", message.ForwardFromChat.UserName),
			Url:      fmt.Sprintf("https://t.me/%s/%d", message.ForwardFromChat.UserName, message.ForwardFromMessageID),
			MessageId: message.ForwardFromMessageID,
		}

	}

	if message.Entities != nil {
		post.Entries = append(post.Entries, *message.Entities...)
	}

	if message.Photo != nil {
		p := *message.Photo

		photoID := p[len(p)-1].FileID
		directURL, _ := c.BotAPI.GetFileDirectURL(photoID)
		f := download(directURL, photoID+".jpeg")

		post.Photo = append(post.Photo, f.Name())
		post.Title = message.Caption

		Posts.append(id, post)

		return true
	}

	if message.Video != nil {
		v := *message.Video

		videoID := v.FileID
		directURL, _ := c.BotAPI.GetFileDirectURL(videoID)
		f := download(directURL, videoID+".mp4")

		post.Video = append(post.Video, f.Name())
		post.Title = message.Caption

		Posts.append(id, post)

		return true
	}

	if message.Document != nil {

		d := *message.Document

		documentID := d.FileID
		directURL, _ := c.BotAPI.GetFile(tgbotapi.FileConfig{FileID: documentID})
		f := download(directURL.Link(config.Token), d.FileName)

		post.Document = append(post.Document, f.Name())
		post.Title = message.Caption

		Posts.append(id, post)

		return true
	}

	post.Title = message.Text

	Posts.append(id, post)

	return true
}

