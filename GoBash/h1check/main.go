package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/mod/sumdb/dirhash"
)

func main() {
	zip, mod := os.Args[1], os.Args[2]

	// go.sum 1 行目: ソースアーカイブ全体
	h, err := dirhash.HashZip(zip, dirhash.Hash1)
	if err != nil {
		panic(err)
	}
	fmt.Println("zip   ", h)

	// go.sum 2 行目: .mod 1 ファイルを名前 "go.mod" としてハッシュ
	h, err = dirhash.Hash1([]string{"go.mod"}, func(string) (io.ReadCloser, error) { return os.Open(mod) })
	if err != nil {
		panic(err)
	}
	fmt.Println("go.mod", h)
}
