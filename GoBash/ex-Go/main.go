package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5" // サーバーを起動するパッケージ
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Hello, Go server!")
	})

	log.Println("Go Bash vol.3")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=vscode&wt=none&l=javascript&width=570&ds=true&dsyoff=0px&dsblur=0px&wc=true&wa=false&pv=0px&ph=0px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=const%2520express%2520%253D%2520require%28%2522express%2522%29%253B%2520%252F%252F%2520%25E3%2582%25B5%25E3%2583%25BC%25E3%2583%2590%25E3%2583%25BC%25E3%2582%2592%25E8%25B5%25B7%25E5%258B%2595%25E3%2581%2599%25E3%2582%258B%25E3%2583%2591%25E3%2583%2583%25E3%2582%25B1%25E3%2583%25BC%25E3%2582%25B8%250Aconst%2520app%2520%253D%2520express%28%29%253B%250A%250Aapp.get%28%2522%252F%2522%252C%2520%28req%252C%2520res%29%2520%253D%253E%2520%257B%250A%2520%2520res.send%28%2522Hello%252C%2520Node%2520server%21%255Cn%2522%29%253B%250A%257D%29%253B%250A%250Aapp.listen%288080%252C%2520%28%29%2520%253D%253E%2520%257B%250A%2520%2520console.log%28%2522Go%2520Bash%2520vol.3%2522%29%253B%250A%257D%29%253B
