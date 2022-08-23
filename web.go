package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/gorilla/mux"
)

//go:embed web
var templatesHTML embed.FS

var funcs = template.FuncMap{
	"HumanTime": func(s string) string {
		now := time.Now()
		if parse, err := time.Parse(time.RFC3339, s); err != nil {
			return s
		} else {
			sub := now.Sub(parse)
			if sub.Hours() > 25 {
				return parse.Format("15:04 02 Jan 2006")
			}

			return time.Time{}.Add(sub).Format("15ч. 04мин. назад")
		}
	},
	"isImage":    func(s string) bool { return strings.HasSuffix(s, ".jpeg") },
	"isVideo":    func(s string) bool { return strings.HasSuffix(s, ".mp4") },
	"isFile": func(s string) bool { return !strings.HasSuffix(s, ".mp4") && !strings.HasSuffix(s, ".jpeg") },
}

var renderStorage = map[string]*template.Template{}
var render = func(name string, f template.FuncMap, patterns ...string) *template.Template {

	defined := []string{
		"webapp.gotm",
	}

	for i := range defined {
		patterns = append(patterns, defined[i])
	}

	if config.Debug {
		for i := range patterns {
			patterns[i] = "web/" + patterns[i]
		}
		fs, err := template.New(name).Funcs(funcs).ParseFiles(patterns...)
		if err != nil {
			log.Printf("Error render template %v", err)
		}
		return fs
	}

	var fs *template.Template
	if _, ok := renderStorage[name]; !ok {
		var err error

		if fs, err = template.New(name).Funcs(funcs).ParseFS(templatesHTML, patterns...); err != nil {
			log.Printf("Error render template %v", err)
		}
		renderStorage[name] = fs
	}
	return renderStorage[name]
}

func home(writer http.ResponseWriter, request *http.Request) {

	data := map[string]interface{}{
		"BotName": botName,
		"BotPic":  botPic,
		"Posts":   Posts.sort(),
	}

	render("index.html", funcs, "index.html").Execute(writer, data)
}

func runWeb() {

	r := mux.NewRouter()
	r.PathPrefix("/posts/").Handler(http.StripPrefix("/posts/", http.FileServer(http.Dir("./posts"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	r.HandleFunc("/", home).Methods(http.MethodGet)
	log.Print(http.ListenAndServe(":"+config.Port, r))
}
