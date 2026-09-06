# 疑問集

## 疑問1

> go mod init でなんで最初から go.mod ができるんだ？


ここはGo Modulesを理解するうえでかなり重要です。

結論から言うと、

> **`go mod init` は「このディレクトリをGo Moduleとして扱う」という宣言をするコマンドだから、宣言を書いた `go.mod` が最初に必要になる**

ということです。

### `go mod init` がやっていること

例えば、

```bash
go mod init example.com/go-server
```

を実行すると、Goは現在のディレクトリに

```text
go.mod
```

を作って、

```go
module example.com/go-server

go 1.XX
```

と書きます。

つまり、

```text
go mod init
      ↓
「このプロジェクトの名前はこれです」
      ↓
go.mod を作る
      ↓
GoがこのディレクトリをModuleとして認識できる
```

という流れです。

---

### なぜ `go.mod` が必要なのか？

Goからすると、例えばこのディレクトリに

```text
my-project/
├── main.go
├── user/
│   └── user.go
└── ...
```

があったとしても、

**「これはどのGo Moduleに属しているプロジェクトなの？」**

という情報がありません。

そこで `go.mod` に

```go
module example.com/go-server
```

と書く。

これによってGoは、

```text
example.com/go-server
        │
        ├── main
        │
        └── user
             ↑
       example.com/go-server/user
```

という**Moduleの境界と名前空間**を認識できます。

go mod init の引数は

「このディレクトリをどんなimport pathのModuleとして扱うか」

を指定している、と考えると分かりやすいです。

---

### そして依存関係もここに書かれる

例えばchiを使うと、

```go
module example.com/go-server

go 1.XX

require github.com/go-chi/chi/v5 v5.x.x
```

となります。

ここで初めて、

```text
go.mod
  │
  ├── このプロジェクト自身
  │      └── module example.com/go-server
  │
  └── 依存しているModule
         └── github.com/go-chi/chi/v5
```

という情報をGoが持てるわけです。

なので、今回の勉強では **`go.mod` は「依存関係を書くファイル」だけではなく、「このプロジェクトは何というGo Moduleなのかを宣言するファイル」**と捉えるのが重要です。

そしてここから一歩進むと、

> **じゃあ `go.mod` に `require` を書いたら、GoはどうやってそのModuleを見つけて、どこにダウンロードしているの？**

という疑問が出てきます。
ここが `GOPROXY` → Module Cache → `go.sum` につながって、今回の登壇テーマの核心に入っていきます。

## 疑問2

> "fmt", "log", "net/http" などの標準パッケージは go.mod と go.sum に何も書いてなくても動く。その理由は？

はい、正確には **Go の標準パッケージは Go のインストール先、つまり `GOROOT` 配下に最初から置かれている** という意味です。

今回コマンドで確認すると、あなたの環境ではこうでした。

```bash
go env GOROOT GOPATH GOMOD GOMODCACHE GOVERSION
```

結果の重要部分:

```txt
GOROOT=/usr/local/go
GOMOD='/Users/shiraki_takuma/ghq/github.com/takumashiraki/conference-report/GoBash/ex-Go/go.mod'
GOVERSION=go1.26.8
```

つまり Go 本体は `/usr/local/go` にあります。

そして標準パッケージを `go list` で見ると:

```bash
go list -f '{{.ImportPath}} {{.Standard}} {{.Dir}}' fmt os strings net/http
```

結果:

```txt
fmt true /usr/local/go/src/fmt
os true /usr/local/go/src/os
strings true /usr/local/go/src/strings
net/http true /usr/local/go/src/net/http
```

ここがポイントです。

`fmt` や `net/http` は:

```txt
Standard=true
Dir=/usr/local/go/src/...
```

になっています。

つまり Go は `import "fmt"` を見たときに、

```go
import "fmt"
```

「これは標準パッケージだ」と判断して、`go.mod` ではなく:

```txt
/usr/local/go/src/fmt
```

を見に行きます。

一方、あなたの自作モジュールはこうでした。

```bash
go list -f 'ImportPath={{.ImportPath}} Standard={{.Standard}} Module={{if .Module}}{{.Module.Path}}@{{.Module.Version}}{{else}}<nil>{{end}} Dir={{.Dir}}' fmt net/http github.com/takumashiraki/conference-report/GoBash/ex-Go
```

結果:

```txt
ImportPath=fmt Standard=true Module=<nil> Dir=/usr/local/go/src/fmt
ImportPath=net/http Standard=true Module=<nil> Dir=/usr/local/go/src/net/http
ImportPath=github.com/.../ex-Go Standard=false Module=github.com/.../ex-Go@ Dir=/Users/.../GoBash/ex-Go
```

標準パッケージは `Module=<nil>` です。

つまり、標準パッケージは **Go Module の依存ではない** ということです。だから `go.mod` に `require fmt` のようには書かれません。

内部の流れをざっくり言うと:

```txt
main.go
  ↓
import "fmt"
  ↓
Go コマンドが import path を解決
  ↓
"fmt" は標準パッケージとして認識される
  ↓
GOROOT/src/fmt を使う
  ↓
go.mod / go.sum には追加しない
```

今の `main.go` は:

```go
import (
	"fmt"
	"log"
	"net/http"
)
```

なので、Go はそれぞれを `/usr/local/go/src/...` から解決しています。

さらに `fmt` の実体も確認できました。

```bash
ls /usr/local/go/src/fmt
```

中には:

```txt
doc.go
errors.go
format.go
print.go
scan.go
...
```

がありました。

つまり `fmt.Println` などは魔法ではなく、Go 本体に含まれている `/usr/local/go/src/fmt` の Go コードを使っています。

`go.mod` / `go.sum` が担当するのは主にこちらです。

```go
import "github.com/gin-gonic/gin"
```

このような **外部モジュール** は Go 本体に含まれていないので、`go.mod` に:

```txt
require github.com/gin-gonic/gin v...
```

のように記録されます。

まとめると:

```txt
標準パッケージ
  → GOROOT/src から読む
  → Go 本体に同梱
  → go.mod / go.sum に書かれない

外部パッケージ
  → module cache などから読む
  → go.mod / go.sum でバージョン管理される
```

なので、「標準パッケージは Go 本体に同梱されている」とは、実際には **`/usr/local/go/src/fmt` のような場所に最初からソースコードがあり、Go コマンドがそこを特別扱いして探している** という意味です。

## 疑問3

> 上記の内容実態がどこにあるのかは理解できた。
> しかし、どうやって標準パッケージなのかを判断しているのかがあまりわかんないです。

判断の根本は **`go.mod` ではなく、`import path の形` と `GOROOT/src に実体があるか`** です。

今回見た Go 本体の実装では、入口はこのあたりです。

[search.go](/usr/local/go/src/cmd/go/internal/search/search.go:519)

```go
func IsStandardImportPath(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	elem := path[:i]
	return !strings.Contains(elem, ".")
}
```

つまり Go はまず import path の **最初の要素に `.` があるか** を見ます。

例えば:

```txt
fmt                  → 最初の要素 fmt        → ドットなし → 標準候補
net/http             → 最初の要素 net        → ドットなし → 標準候補
github.com/gin-gonic/gin → 最初の要素 github.com → ドットあり → 標準候補ではない
example.com/foo      → 最初の要素 example.com → ドットあり → 外部モジュール側
```

なので `import "fmt"` は、この時点で「標準パッケージかもしれない」と判定されます。

ただし、**ドットがないだけでは標準パッケージ確定ではありません**。次に実体確認があります。

[build.go](/usr/local/go/src/cmd/go/internal/modload/build.go:37)

```go
func findStandardImportPath(path string) string {
	if search.IsStandardImportPath(path) {
		if modindex.IsStandardPackage(cfg.GOROOT, cfg.BuildContext.Compiler, path) {
			return filepath.Join(cfg.GOROOT, "src", path)
		}
	}
	return ""
}
```

ここで:

```txt
1. IsStandardImportPath(path)
2. modindex.IsStandardPackage(GOROOT, compiler, path)
```

の両方を満たしたら標準パッケージです。

さらに `modindex.IsStandardPackage` のフォールバック実装を見ると、実際には `GOROOT/src/<import path>` に `.go` ファイルがあるか確認しています。

[gc.go](/usr/local/go/src/internal/goroot/gc.go:19)

```go
dir := filepath.Join(goroot, "src", path)
dirents, err := readDir(dir)
if err != nil {
	return false
}
for _, dirent := range dirents {
	if strings.HasSuffix(dirent.Name(), ".go") {
		return true
	}
}
return false
```

つまり `fmt` の場合はこうです。

```txt
import "fmt"
  ↓
最初の要素は "fmt"
  ↓
"." を含まない
  ↓
標準パッケージ候補
  ↓
/usr/local/go/src/fmt を見る
  ↓
.go ファイルがある
  ↓
標準パッケージとして扱う
```

実際に確認したコマンド:

```bash
go list -f '{{.ImportPath}} {{.Standard}} {{.Dir}}' fmt net/http
```

結果:

```txt
fmt true /usr/local/go/src/fmt
net/http true /usr/local/go/src/net/http
```

一方で、存在しないけどドットなしの名前を指定すると:

```bash
go list notexistpkg
```

結果:

```txt
package notexistpkg is not in std (/usr/local/go/src/notexistpkg)
```

これはかなり重要です。

`notexistpkg` は `.` がないので **標準候補** として見られます。でも `/usr/local/go/src/notexistpkg` に実体がないので、「std にない」と言われます。

逆にドットありだと:

```bash
go list example.com/notexist/pkg
```

結果:

```txt
no required module provides package example.com/notexist/pkg; to add it:
	go get example.com/notexist/pkg
```

これは `example.com` に `.` があるので、標準ライブラリ候補ではなく **外部モジュール** として扱われています。

なので結論はこうです。

```txt
go.mod にあるか？
  → 標準パッケージ判定の根本ではない

import の書き方？
  → import path の文字列は判定に使われる

実際の判定
  → まず import path の最初の要素に "." がないかを見る
  → 次に GOROOT/src/<path> に Go パッケージが存在するか見る
  → 存在すれば Standard=true
```

`go.mod` はその後の外部モジュール解決で使われるもの、という理解が近いです。標準パッケージはその前に `GOROOT` 側で解決されるので、`require fmt` のような記述が不要になります。