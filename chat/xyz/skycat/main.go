package main

import (
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func(){
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main(){
	// 新しいチャットルームを作る
	r := newRoom()
	// ルートアクセスされたら char.html を表示するようハンドリング
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// 
	http.Handle("/room")

	// ゴルーチンでチャットルームを起動する
	go r.run()

	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
