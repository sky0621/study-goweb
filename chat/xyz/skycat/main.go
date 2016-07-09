package main

import (
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte(`
			<html>
				<head>
					<title>チャット</title>
				</head>
				<body>
					チャットしましょう！
				</body>
			</html>
		`))
	})

	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
