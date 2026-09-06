# README

## インストール

> // indirect は「このモジュールのコードが直接 import していない依存」を示すマーカーです。Go 1.17 の module graph pruning 以降、メインモジュールの go.mod にはビルドに必要な依存が indirect も含めて展開して書かれるようになりました。

Go 1.17 より前のgoを入れるので、 [go1.26.8](https://go.dev/dl/go1.26.8.darwin-arm64.pkg) を入れる

### PATH設定

```zsh
open -e ~/.zshrc
```

```.zshrc
# go

export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$HOME/go/bin
```

```zsh
source ~/.zshrc
```

go1.26.8になっていたらOK

```zsh
go version
go version go1.26.8 darwin/arm64
```

## リポジトリ作成


```zsh
cd GoBash/ex-Go

go mod init github.com/takumashiraki/conference-report/GoBash/ex-Go
```

## サーバーの作成

main.go 作成

```zsh
touch main.go
```

<details>

<summary>Tips for collapsed sections</summary>

```go
package main

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
```

</details>

サーバーを起動
```zsh
go run .
```

別のタブ作成

```zsh
❯ curl http://localhost:8080
Hello, Go server!
```