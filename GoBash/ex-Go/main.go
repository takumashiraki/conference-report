package main // これはGoにとって、「このパッケージは実行可能プログラムです。」という意味になります。

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, Go server!")
	})

	log.Println("server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
