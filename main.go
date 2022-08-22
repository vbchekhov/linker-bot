package main

import (
	"encoding/json"
	"github.com/vbchekhov/skeleton"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

var Posts = MediaGroup{}

type MediaGroup map[int]ChannelPost

func (m MediaGroup) init() {
	f, _ := ioutil.ReadFile("posts/posts.json")
	json.Unmarshal(f, &m)
}
func (m MediaGroup) update() {
	bytes, _ := json.Marshal(m)
	ioutil.WriteFile("posts/posts.json", bytes, os.ModePerm)
}
func (m MediaGroup) append(id int, title, photo string) {

	if p, ok := m[id]; ok {
		p.Photo = append(p.Photo, photo)
		m[id] = p
	} else {
		m[id] = ChannelPost{
			title,
			time.Now().Format(time.RFC3339),
			[]string{photo}}
	}

}
func (m MediaGroup) sort(count string) MediaGroup {
	var arr []ChannelPost

	for k := range m {
		arr = append(arr, m[k])
	}

	sort.Slice(arr, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, arr[i].Time)
		t2, _ := time.Parse(time.RFC3339, arr[j].Time)

		return t1.UnixNano() > t2.UnixNano()
	})

	posts := MediaGroup{}
	for i := 0; i < len(arr); i++ {
		posts[strconv.Itoa(i)] = arr[i]
		if strconv.Itoa(i) == count {
			break
		}
	}

	return posts
}

type ChannelPost struct {
	Title string
	Time  string
	Photo []string
}

func download(url, filename string) *os.File {

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	f, _ := os.Create("posts/" + filename)
	io.Copy(f, resp.Body)

	defer f.Close()

	return f
}

func RunBot() {

	Posts.init()

	bot := skeleton.NewBot(config.Token)
	bot.HandleFunc("(.*)", func(c *skeleton.Context) bool {

		post := c.Update.Message
		if post.Photo != nil {
			p := *post.Photo

			photoID := p[len(p)-1].FileID
			directURL, _ := c.BotAPI.GetFileDirectURL(photoID)
			f := download(directURL, photoID+v.MimeType)

			Posts.append(post.MessageID, post.Caption, f.Name())
			Posts.update()

			return true
		}

		if post.Video != nil {
			v := *post.Video

			videoID := v.FileID
			directURL, _ := c.BotAPI.GetFileDirectURL(videoID)
			f := download(directURL, videoID+v.MimeType)

			Posts.append(post.MessageID, post.Caption, f.Name())
			Posts.update()

			return true
		}

		if post.Document != nil {

			v := *post.Document

			documentID := v.FileID
			directURL, _ := c.BotAPI.GetFileDirectURL(documentID)
			f := download(directURL, documentID+v.MimeType)

			Posts.append(post.MessageID, post.Caption, f.Name())
			Posts.update()

			return true

		}

		

		Posts.append(post.MessageID, post.Text, "")
		Posts.update()

		log.Printf("%+v", Posts)

		return true
	}).Border(skeleton.Private).Methods(skeleton.Messages)

	// log.Printf("%+v", Posts.array()  )

	bot.Run()
}
