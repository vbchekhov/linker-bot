package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var Posts = MediaGroup{}

type MediaGroup map[string]Post

func (m MediaGroup) init() {
	os.Mkdir(config.Folder, os.ModePerm)

	f, err := ioutil.ReadFile(config.Folder + "posts.json")
	if err != nil {
		os.Create(config.Folder + "posts.json")
		f, _ = ioutil.ReadFile(config.Folder + "posts.json")
	}
	json.Unmarshal(f, &m)
}

func (m MediaGroup) update() {
	bytes, _ := json.Marshal(m)
	ioutil.WriteFile(config.Folder+"posts.json", bytes, os.ModePerm)
}

func (m MediaGroup) append(id string, post Post) {

	if p, ok := m[id]; ok {
		p.Photo = append(p.Photo, post.Photo...)
		p.Video = append(p.Video, post.Video...)
		p.Document = append(p.Document, post.Document...)
		p.Entries = append(p.Entries, post.Entries...)

		m[id] = p
	} else {
		m[id] = post

	}

	m.update()
}

func (m MediaGroup) sort() MediaGroup {
	var arr []Post

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
		
	}

	return posts
}
func download(url, filename string) *os.File {

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	f, _ := os.Create(config.Folder + filename)
	io.Copy(f, resp.Body)

	defer f.Close()

	return f
}

type Post struct {
	Title    string
	Time     string
	Photo    []string
	Video    []string
	Document []string
	Entries  []tgbotapi.MessageEntity
	Metadata Metadata
}

type Metadata struct {
	Title     string
	UserName  string
	Avatar    string
	Group     string
	Url       string
	MessageId int
}
