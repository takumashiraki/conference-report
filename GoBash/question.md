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

## 疑問4

> 標準パッケージを入れると、go.sum が生まれた。
> - go.sum は誰がいつ作るのか？
> - 内部ではどうやって作るのか？

```sh
cd GoBash/ex-Go

go get github.com/go-chi/chi/v5
go: downloading github.com/go-chi/chi/v5 v5.3.2
go: downloading github.com/go-chi/chi v1.5.5
go: added github.com/go-chi/chi/v5 v5.3.2

ls
README.md   go.mod      go.sum      main.go     question.md

cat go.sum
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
```

ただし `fmt` や `net/http` のような **標準ライブラリだけ**なら、外部モジュールを取得しないので基本的に `go.sum` は増えません。`github.com/go-chi/chi/v5` は標準パッケージではなく、外部モジュールです。

内部ではだいたいこう動きます。

1. `go.mod` を見て、必要なモジュールとバージョンを決める
2. `GOPROXY`、または GitHub などの VCS からモジュールを取得する
3. モジュールの `.zip` と `go.mod` をローカルのモジュールキャッシュに保存する
4. 取得した内容から暗号学的ハッシュを計算する
5. 既存の `go.sum` または checksum database、通常は `sum.golang.org`、と照合する
6. 問題なければ `go.sum` に記録する

この 6 ステップが具体的に何をしているのかを、実際にモジュールキャッシュを覗きながら追ってみます。この 6 ステップを起動するのは `go get` / `go build` / `go mod tidy` のいずれでもよく、外部モジュールが必要になった時点でステップ 1 を `modload`、ステップ 2 以降を `modfetch` が担当します。

### ステップ1: 必要なモジュールとバージョンを決める

ここでは実は 2 つの別処理が走っています。

**(a) import path からモジュール候補を切り出す**

`github.com/go-chi/chi/v5` という import path を見ても、go コマンドは「どこまでがモジュール名で、どこからがパッケージのサブディレクトリか」を知りません。

```go
import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5" // import path
)
```

そこで go コマンドは、次の 2 段階でモジュールを決めます（[Resolving a package to a module](https://go.dev/ref/mod#resolve-pkg-mod)）。

| 段階 | 照合する相手 | ネットワーク |
| --- | --- | --- |
| (a-1) ビルドリスト照合 | 手元のビルドリスト | 使わない |
| (a-2) プロキシ問い合わせ | `GOPROXY` | 使う。(a-1) で見つからず、`go get` / `go mod tidy` のときだけ |

**(a-1) ビルドリスト照合**

ビルドリストは、`go.mod` から決まる「ビルドに使うモジュールとバージョンの一覧」です（決め方は (b) の MVS）。`go list -m all` で確認できます。

```bash
go list -m all
github.com/takumashiraki/conference-report/GoBash/ex-Go # ← メインモジュール（バージョンなし）
github.com/go-chi/chi/v5 v5.3.2                          # ← 依存モジュール（モジュールパス バージョン）
```

go コマンドはまず、ビルドリストの中から、モジュールパス（左列）が import path の接頭辞になっているモジュールを探します。
完全一致も接頭辞に含まれ、接頭辞からはみ出た部分はモジュール内のサブディレクトリになります。

※ `GOPROXY=off` でネットワークを切っても、ビルドリストだけで解決できました。

```bash
GOPROXY=off go list -f '{{.ImportPath}} → {{.Module.Path}} {{.Module.Version}}' \
  github.com/go-chi/chi/v5 github.com/go-chi/chi/v5/middleware
github.com/go-chi/chi/v5 → github.com/go-chi/chi/v5 v5.3.2
github.com/go-chi/chi/v5/middleware → github.com/go-chi/chi/v5 v5.3.2
```

`go.mod` に `require` があれば、この段階でモジュールが決まります。後続のプロキシには問い合わせません。

**(a-2) プロキシ問い合わせ**

ビルドリストで見つからない場合、`go get` と `go mod tidy` は新しいモジュールを探しに行きます。
最初の `go get` の時点では `go.mod` に `require` がなく、ビルドリストにはメインモジュールしかないので、こちらに進みます。

import path を末尾から 1 要素ずつ削った接頭辞それぞれを、パッケージを提供しうるモジュールパスの候補にします。`GOPROXY` の各エントリに対して、候補ごとに最新バージョンを要求します。公式リファレンスの例では、これらの要求は 1 つのプロキシに対して並列に送られます。

```txt
github.com/go-chi/chi/v5   の最新バージョン
github.com/go-chi/chi      の最新バージョン
github.com/go-chi          の最新バージョン
github.com                 の最新バージョン
```

モジュールキャッシュが空の状態で `go get -x` を実行すると、go コマンドが実際にプロキシへ送った要求が見えます。最新バージョンを知るために叩いているのは `@v/list`（バージョン一覧）で、バージョンが返ってくる（要求に成功する）のは 2 つでした。

```bash
go get -x github.com/go-chi/chi/v5 2>&1 | grep '/@v/list: '
# get https://proxy.golang.org/github.com/@v/list: 404 Not Found
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/list: 200 OK
# get https://proxy.golang.org/github.com/go-chi/@v/list: 404 Not Found
# get https://proxy.golang.org/github.com/go-chi/chi/@v/list: 200 OK
```

（行末の所要時間は省略しています。4 つの要求は並列に送られるので、並び順は実行するたびに変わります。）

要求に成功したモジュールパスについては、最新バージョンのモジュールを取得し、要求されたパッケージを含むかどうかを確認します。パッケージを含むモジュールが複数あれば、パスが最も長いモジュールを使います。

これが、先ほどのログに出てきた謎の 1 行の正体です。

```sh
go: downloading github.com/go-chi/chi/v5 v5.3.2
go: downloading github.com/go-chi/chi v1.5.5   ← これ
```

`github.com/go-chi/chi` の中に `v5` というパッケージディレクトリがあるかもしれないので、確認のために取得しています。実際、キャッシュには両方残っています。

`chi chi@v1.5.5` の `chi` は、chi v1 のことではありません。中を見ると `v5@v5.3.2/` が入っていました。つまり `chi/` は `chi/v5@v5.3.2` の親ディレクトリです。「両方残っています」という記述は正しいです。ただ、`chi` と `chi@v1.5.5` が v1 系の 2 つに見えて誤読されやすくなっています。

```sh
ls $(go env GOMODCACHE)/github.com/go-chi/
chi   chi@v1.5.5
```

**(b) MVS（Minimal Version Selection）**

- 候補が確定したら、モジュールグラフをたどってバージョンを決めます。
- MVS はメインモジュールからグラフをたどり、各モジュールについて要求された中で**最も高い**バージョンを記録します。
- たどり終えた時点で記録されている最も高いバージョンの集合がビルドリストになります。これが、すべての要求を満たす最小のバージョンです。
- 「最も高い」と「最小」は矛盾しません。たとえば A が chi v5.0.0 以上、B が v5.2.0 以上を要求していれば、両方を満たす最小のバージョンは v5.2.0 です。それより新しいバージョンが公開されていても、誰も要求していなければ選びません。

> （[Minimal version selection (MVS)](https://go.dev/ref/mod#minimal-version-selection)）

MVS を行うために、go コマンドは依存モジュールの複数のバージョンの `go.mod` を読み込むことがあります。プロキシから取得するときは `go.mod` をモジュールの残りの内容とは別に取得します。そのため後述の proxy プロトコルでは `.mod` と `.zip` が別エンドポイントになっていて、`go.sum` にも `/go.mod` 付きの行（`go.mod` だけのハッシュ）と付かない行（`.zip` の中身のハッシュ）の 2 行が載ります（[go.sum files](https://go.dev/ref/mod#go-sum-files)）。

### ステップ2: モジュールプロキシまたは VCS から取得する

デフォルトの設定はこうなっています。

```sh
go env GOPROXY
https://proxy.golang.org,direct

# https://proxy.golang.org モジュールのキャッシュサーバー
# , フォールバックのチェーン (「1 番目で取れなければ 2 番目を試す」という順番付きの候補リスト)
# direct は「プロキシ(https://proxy.golang.org)を使わず直接 VCS(Git) を叩く」という特殊な値で、プロキシが 404 / 410 を返したら次に進みます。
```

|場面|経路|例|
|---|---|---|
|公開モジュール（GitHub の OSS など）|**プロキシ経由**|`github.com/go-chi/chi/v5`|
|プライベートモジュール（社内リポジトリ）|**direct**|このマシンでは `gitlab.tokyo.optim.co.jp/*`|
|公開モジュールだが、ミラーが配信を拒否した場合|direct（フォールバック）|公式ヘルプによると、主に法的な理由のとき|

つまり取得経路は 2 通りあります。違いは「`git` を誰が実行するか」です。

```txt
【プロキシ経由（普段）】
 go コマンド ──HTTP GET──▶ proxy.golang.org ──git──▶ github.com/go-chi/chi
                            （git の実行と zip 化はプロキシが代行する）

【direct】
 go コマンド ──git ls-remote / git fetch──▶ github.com/go-chi/chi
             （手元の PC が git を実行し、zip も自分で組み立てる）
```

#### プロキシ経由の場合

ここでいうプロキシは `proxy.golang.org` のことです。Go チームが運営し Google がホストしている、公開モジュールのミラー（キャッシュサーバー）です。
GitHub などから取ってきたソースを zip にして保存しておき、go コマンドに HTTP で配ります。社内ネットワークの HTTP プロキシ（`HTTP_PROXY`）とは別物です。

作者がアップロードする npm registry とは違い、初めて要求されたバージョンはプロキシ自身が元リポジトリから取ってきます。いったんキャッシュされると、作者が元リポジトリでリリースを消しても取得でき続けます（[proxy.golang.org](https://proxy.golang.org/) の FAQ）。

プロキシへのリクエストは、この 5 種類の HTTP GET だけです。

```txt
GET /<module>/@v/list              バージョン一覧
GET /<module>/@v/<version>.info    メタデータ (JSON)
GET /<module>/@v/<version>.mod     そのバージョンの go.mod
GET /<module>/@v/<version>.zip     ソースアーカイブ (zip)
GET /<module>/@latest              最新版の解決（タグ付きバージョンが無いときだけ使う）
```

バージョン無指定の `go get` だと、まず `@v/list` でバージョン一覧を取り、その中で最も高いリリースについて `.info` → `.mod` → `.zip` の順に叩きます。`.mod` と `.zip` の間には `sum.golang.org` へのチェックサム照合が挟まります。

空のモジュール・空のキャッシュで実行し、`chi/v5` の取得に関わる行だけを残して URL を縮めた出力です。

```txt
❯ cd "$(mktemp -d)" && go mod init demo
❯ GOMODCACHE=$(mktemp -d) GOFLAGS=-modcacherw go get -x github.com/go-chi/chi/v5 2>&1

# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/list 200 OK        ← ① バージョン一覧
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/v5.3.2.info 200 OK ← ② 最も高いリリースのメタデータ
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/v5.3.2.mod 200 OK  ← ③ go.mod
# get https://sum.golang.org/lookup/github.com/go-chi/chi/v5@v5.3.2 200 OK    ← ④ チェックサム照合
go: downloading github.com/go-chi/chi/v5 v5.3.2
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/v5.3.2.zip              200 OK  ← ⑤ ソースアーカイブ
go: added github.com/go-chi/chi/v5 v5.3.2
```

① の一覧の末尾はこうなっていて、最も高い `v5.3.2` が選ばれます。

```txt
❯ curl -s https://proxy.golang.org/github.com/go-chi/chi/v5/@v/list | sort -V | tail -3
v5.3.0
v5.3.1
v5.3.2
```

実際の出力には、ほかに次の行も混ざります（上では省きました）。

- 親パス（`github.com/go-chi/chi`、`github.com/go-chi`、`github.com`）の `@v/list`。`chi/v5` が短いパスのモジュールの中の `v5/` ディレクトリではないかを確かめるため、並列に問い合わせます。`github.com/go-chi/chi` だけが 200 を返し、その最新版 `v1.5.5` の `.mod` / `.info` / `.zip` も取得します
- `sum.golang.org/tile/...`。チェックサム照合で、ログの木の一部を取りに行く行です
- リクエスト開始の行（`# get <URL>`）と所要時間

`@latest` は一度も叩かれていません。go コマンドのソース（`cmd/go/internal/modload/query.go`）を読むと、一覧にリリースもプレリリースも無いときにだけ `@latest` を呼んでいます。タグを一度も打っていないリポジトリを `go get` したときがこれにあたり、プロキシが最新コミットから疑似バージョンを作って返します。

`.info` には取得元の証跡が残っています。

```sh
cat $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.info
```

```json
{"Version":"v5.3.2","Time":"2026-08-20T09:37:52Z",
 "Origin":{"VCS":"git","URL":"https://github.com/go-chi/chi",
 "Hash":"38939062c5df4d3e8814aad1a488983112627ced","Ref":"refs/tags/v5.3.2"}}
```

プロキシ自身が「git の `refs/tags/v5.3.2`、コミット `3893906` から作った zip だ」と申告しているわけです。

#### direct の場合

VCS は Version Control System の略で、Git や Mercurial などのバージョン管理システムのことです。`direct` は、プロキシを通さずに go コマンドが VCS を直接実行して取得する経路です。上の `.info` にある「git でタグを引き、コミットを取り出して zip にする」作業を、手元の PC で go コマンド自身が行います。公開サーバーに対して使う VCS は、既定で `git` と `hg` だけです（`go help vcs`）。

direct になるのは次の 3 つのときです。

- `GOPROXY` で前にあるプロキシが 404 / 410 を返したとき。500 やタイムアウトでは次へ進まず、そこで止まります。区切りを `,` ではなく `|` にすると、どのエラーでも次へ進みます（`cmd/go/internal/modfetch/proxy.go` の `fallBackOnError`）
- `GOPROXY=direct` と明示したとき
- モジュールパスが `GOPRIVATE`（または `GONOPROXY`）に一致するとき。社内リポジトリなどが該当します（`go help private`）

1 つ目のフォールバックは、何を要求しても 404 を返すプロキシを立てると観察できます。

```bash
# 空ディレクトリを配信するだけのサーバー = 何を要求しても 404 を返すプロキシ
mkdir -p /tmp/empty && python3 -m http.server 18080 --directory /tmp/empty &

# 普段のキャッシュを汚さないよう、使い捨ての GOMODCACHE に取得する
GOMODCACHE=/tmp/modcache GOFLAGS=-modcacherw GOPROXY=http://localhost:18080,direct \
  go mod download -x github.com/go-chi/chi/v5@v5.3.2 2>&1 | grep -E '^# get .*: |^cd .*git (ls-remote|.*fetch)'
```

```txt
# get http://localhost:18080/github.com/go-chi/chi/v5/@v/v5.3.2.info: 404 File not found
cd …/cache/vcs/672f73a6…; git ls-remote -q --end-of-options https://github.com/go-chi/chi
cd …/cache/vcs/672f73a6…; git -c protocol.version=2 fetch -f --depth=1 --end-of-options origin refs/tags/v5.3.2:refs/tags/v5.3.2
# get http://localhost:18080/github.com/go-chi/chi/v5/@v/v5.3.2.mod: 404 File not found
# get http://localhost:18080/sumdb/sum.golang.org/supported: 404 File not found
# get https://sum.golang.org/lookup/github.com/go-chi/chi/v5@v5.3.2: 200 OK
# get http://localhost:18080/github.com/go-chi/chi/v5/@v/v5.3.2.zip: 404 File not found
```

（パスと行末の所要時間、`sum.golang.org/tile/...` の行は省略しています。）

上から順に読むとこうなります。

1. プロキシに `.info` を要求して 404 が返る
2. 次の候補 `direct` に進み、`git ls-remote` でタグ一覧を取り、`git fetch --depth=1` で `v5.3.2` のタグだけを取得する
3. `.mod` と `.zip` も毎回まずプロキシに要求し、404 なので、取得済みの git リポジトリから取り出す
4. checksum database への問い合わせも同じ順で進む。プロキシの `/sumdb/sum.golang.org/supported` が 404 なので、`sum.golang.org` へ直接問い合わせる

経路は違っても、取得結果は同じです。`.info` の `Origin` は、プロキシ経由のときとまったく同じ内容でした。go コマンドが自前で組み立てた zip の `.ziphash` も、`go.sum` と同じ `h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=` です。ステップ4 で見るとおり、H1 は zip のバイト列ではなくファイル名と中身から計算するので、誰が zip を作っても一致します。違いは、direct だと `$GOMODCACHE/cache/vcs/` に git のベアリポジトリが残ることです。

`direct` を外して `GOPROXY=http://localhost:18080` だけにすると、次の候補がないので 404 のまま失敗します。

```txt
go: github.com/go-chi/chi/v5@v5.3.2: reading http://localhost:18080/github.com/go-chi/chi/v5/@v/v5.3.2.info: 404 File not found
```

社内モジュールをプロキシに漏らしたくない場合（3 つ目の入口）の制御はこのあたりです。

```sh
GOPRIVATE   GONOPROXY   GONOSUMDB   GOVCS
```

### ステップ3: モジュールキャッシュに保存する

保存先は `GOMODCACHE`（既定は `~/go/pkg/mod`）です。

```sh
go env GOMODCACHE
/Users/opm008296/go/pkg/mod
```

この直下に、役割の違う 2 つの置き場が同じ階層で並んでいます。

```txt
~/go/pkg/mod/
├── cache/                          ← 層1: 取得したものをそのまま置く
│   ├── download/
│   │   ├── github.com/go-chi/chi/v5/@v/
│   │   │   ├── list                ←   手元にあるバージョンの一覧
│   │   │   ├── v5.3.2.info         ←   バージョンと取得元（git のコミット）
│   │   │   ├── v5.3.2.mod          ←   go.mod 単体
│   │   │   ├── v5.3.2.zip          ←   ソースアーカイブ（ソース一式の zip）
│   │   │   ├── v5.3.2.ziphash      ←   zip のハッシュ（go.sum と同じ値）
│   │   │   └── v5.3.2.lock
│   │   └── sumdb/                  ←   チェックサム DB の応答（ステップ5）
│   └── vcs/                        ←   direct 取得時の git リポジトリ（ステップ2）
│
└── github.com/go-chi/chi/          ← 層2: zip を展開したソース（ビルドが読む）
    ├── v5@v5.2.5/
    ├── v5@v5.3.1/
    └── v5@v5.3.2/                  ←   chi.go, mux.go, go.mod, middleware/ ...
```

`go get` 1 回で、両方が次の順にできます。

```txt
1. プロキシ（または git）から .info / .mod / .zip を取得し、層1に保存する
2. zip のハッシュを計算して go.sum と照合し、.ziphash に記録する
3. zip を層2に展開し、読み取り専用にする
4. 以降の go build / go test は層2のソースだけを読む（zip は開かない）
```

#### 1 コマンドで両方の層を見る

`go mod download -json` を使うと、1 つのモジュールが両方の層のどこに置かれているかがまとめて出力されます。

```sh
go mod download -json github.com/go-chi/chi/v5@v5.3.2
{
	"Path": "github.com/go-chi/chi/v5",
	"Version": "v5.3.2",
	"Info": "/Users/opm008296/go/pkg/mod/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.info",
	"GoMod": "/Users/opm008296/go/pkg/mod/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.mod",
	"Zip": "/Users/opm008296/go/pkg/mod/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.zip",
	"Dir": "/Users/opm008296/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2",
	"Sum": "h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=",
	...
}
```

`Info` / `GoMod` / `Zip` が層1、`Dir` が層2 です。`Sum` は `go.sum` に書かれる値で、層1の `.ziphash` の中身と同じです。

#### 層1: `cache/download`（生データ）

外側から順に `ls` していくと、層1の中身が見えます。

```sh
ls $(go env GOMODCACHE)/cache
download  lock  vcs

ls $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/
list          v5.2.5.zip      v5.3.1.mod      v5.3.2.lock
v5.2.5.info   v5.2.5.ziphash  v5.3.1.zip      v5.3.2.mod
v5.2.5.lock   v5.3.1.info     v5.3.1.ziphash  v5.3.2.zip
v5.2.5.mod    v5.3.1.lock     v5.3.2.info     v5.3.2.ziphash
```

1 バージョンにつき `.info` / `.lock` / `.mod` / `.zip` / `.ziphash` の 5 ファイルがあり、手元にあるバージョンの数だけ並びます。中身はテキストなので `cat` で読めます。

```sh
cd $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/

cat list
v5.2.5
v5.3.1
v5.3.2

cat v5.3.2.info
{"Version":"v5.3.2","Time":"2026-08-20T09:37:52Z","Origin":{"VCS":"git","URL":"https://github.com/go-chi/chi","Hash":"38939062c5df4d3e8814aad1a488983112627ced","Ref":"refs/tags/v5.3.2"}}

cat v5.3.2.mod
module github.com/go-chi/chi/v5
...
go 1.23

cat v5.3.2.ziphash
h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
```

| ファイル | 役割 |
|---|---|
| `.info` | バージョンと取得元のメタデータ |
| `.mod` | そのバージョンの `go.mod` 単体 |
| `.zip` | ソースアーカイブ。そのバージョンのソース一式を 1 つにまとめた zip |
| `.ziphash` | 検証済みハッシュのキャッシュ（ステップ5で効く） |
| `.lock` | 並行する `go` プロセス同士の排他用（0 バイト） |

パスには 2 つ仕掛けがあります。

ひとつは**大文字のエスケープ**です。大文字小文字を区別しないファイルシステム（macOS がまさにそれ）で衝突しないよう、大文字は `!` + 小文字に変換されます。

```txt
github.com/!burnt!sushi/   ← github.com/BurntSushi/
```

もうひとつは**メジャーバージョンサフィックス**です。`/v5` はモジュールパスの一部なので、ディレクトリ階層にそのまま現れます（`chi/v5/@v/`）。

層1のディレクトリ構成は、モジュールプロキシの URL 構成（ステップ2 の `-x` 出力に出た `.../@v/v5.3.2.info` など）とそのまま同じです。そのため、層1はそのままプロキシとして使えます。

```sh
# 空のキャッシュに、層1だけをプロキシにして取得する（ネットワークに出ない）
P=$(go env GOMODCACHE)/cache/download
GOMODCACHE=$(mktemp -d) GOFLAGS=-modcacherw GOSUMDB=off GOPROXY=file://$P \
  go mod download -json github.com/go-chi/chi/v5@v5.3.2
# → "Sum": "h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=" で取得できる
```

#### 層2: `<module>@<version>`（展開済みソース）

層2は、バージョンごとに別のディレクトリになっています。

```sh
ls $(go env GOMODCACHE)/github.com/go-chi/chi/
v5@v5.2.5  v5@v5.3.1  v5@v5.3.2
```

中を見ると、chi のリポジトリのソースがそのまま並んでいます。`go build` が読むのはここです。

```sh
ls $(go env GOMODCACHE)/github.com/go-chi/chi/v5@v5.3.2
CHANGELOG.md      chain.go          mux_test.go
CONTRIBUTING.md   chi.go            path_value_test.go
LICENSE           context.go        pattern_test.go
Makefile          context_test.go   testdata
README.md         go.mod            tree.go
SECURITY.md       middleware        tree_test.go
_examples         mux.go
```

層2の中身は、層1の zip を展開したものです。zip の中のパスが層2のパスと一致すること、`.mod` と展開後の `go.mod` が同じファイルであることで確かめられます。

```sh
unzip -l $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.zip | head -4
Archive:  .../cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
      723  00-00-1980 00:00   github.com/go-chi/chi/v5@v5.3.2/.github/FUNDING.yml

diff $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.mod \
     $(go env GOMODCACHE)/github.com/go-chi/chi/v5@v5.3.2/go.mod && echo same
same
```

| | 層1（`cache/download`） | 層2（`<module>@<version>`） |
|---|---|---|
| 置いてあるもの | `.info` `.mod` `.zip` `.ziphash` | 展開済みのソース |
| 使う処理 | `go get` / `go mod download` / 検証 / プロキシ配信 | `go build` / `go test` |
| 大きさ（chi v5.3.2） | zip 132K | 588K・86 ファイル |
| 権限 | 通常 | 読み取り専用 |

層2は**読み取り専用**で作られます。

```sh
ls -ld $(go env GOMODCACHE)/github.com/go-chi/chi/v5@v5.3.2
dr-xr-xr-x@ 24 ... github.com/go-chi/chi/v5@v5.3.2
```

`r-x` です。依存先をうっかり編集してハッシュ検証を壊すことを、パーミッションで物理的に防いでいます。消すときに `go clean -modcache` が必要なのもこれが理由です。

### ステップ4: 暗号学的ハッシュを計算する

`h1:` は、SHA-256 を使うハッシュ方式「Hash1」の識別子です。`golang.org/x/mod/sumdb/dirhash` のソースにも `Hash1 is the "h1:" directory hash function, using SHA-256.` とあります。値の先頭に方式名を書いておくことで、将来別の方式が増えても接頭辞で区別できます。

ハッシュをかける対象は**ソースアーカイブ**です。ソースアーカイブとは、ステップ3 で層1に保存した `v5.3.2.zip` のことで、そのバージョンのソース一式（chi v5.3.2 なら 86 ファイル）を 1 つにまとめた zip ファイルです。

中身で大事なのは、**zip のバイト列の SHA-256 ではない**という点です。手順は次の 4 段です。

```txt
1. ソースアーカイブ内の全ファイルについて sha256 を計算する
2. "<sha256 の hex>  <ファイル名>\n" という行を作る（スペースは 2 個）
3. ファイル名でソートして全行を連結する → これが「リスト」
4. そのリスト全体の sha256 を取り、その 32 バイトを base64 して "h1:" を付ける
```

#### 手順どおりに計算してみる

chi v5.3.2 のソースアーカイブを使い、4 段を 1 つずつシェルで再現します。

**1. ソースアーカイブ内の全ファイルについて sha256 を計算する**

```sh
cd $(mktemp -d)
unzip -q $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.zip
find github.com -type f | wc -l
      86
# ファイル数 86 個

shasum -a 256 github.com/go-chi/chi/v5@v5.3.2/.gitignore
785f18e2ac99c66b81d8d185aa1a65cfe4f39b12bef24c1f8037d85840490929  github.com/go-chi/chi/v5@v5.3.2/.gitignore
# ↑ これが "<sha256 の hex>  <ファイル名>\n" という行になる
```

ファイル名は zip 内のフルパスで、`github.com/go-chi/chi/v5@v5.3.2/` から始まります。展開先で `github.com` から相対パスを取れば、そのまま同じ名前になります。

**2〜3. 行を作り、ファイル名でソートして連結する**

`shasum` の出力形式は `<sha256 の hex>  <ファイル名>`（スペース 2 個）で、手順2 の行の形式と同じです。そのため、ソートしたファイル名の順に `shasum` を並べるだけでリストができます。

```sh
find github.com -type f | LC_ALL=C sort \
  | while read f; do shasum -a 256 "$f"; done > list.txt

wc -l < list.txt; wc -c < list.txt
      86
   10452
# 事前に確認したファイル数 86 個があった

cat list.txt
0e236eacfd30c201e85ff368730f13e6dfb24d21a72c57c8bed0a099bebb3c7b  github.com/go-chi/chi/v5@v5.3.2/.github/FUNDING.yml
4332e758142ca72fd505b77fa58518e5586caae8a6a1260ae2758d4aedeeccd3  github.com/go-chi/chi/v5@v5.3.2/.github/workflows/ci.yml
785f18e2ac99c66b81d8d185aa1a65cfe4f39b12bef24c1f8037d85840490929  github.com/go-chi/chi/v5@v5.3.2/.gitignore
65c419049d2e6efc04b110426e4b5efc628a4ee6b8f080092f25a0dbbb35f071  github.com/go-chi/chi/v5@v5.3.2/CHANGELOG.md
961da16c4fcec58d600fd49cc6ab0dd632d143646d4fcb2f5b1818fb0197fbe7  github.com/go-chi/chi/v5@v5.3.2/CONTRIBUTING.md
a2d51b7515acfaff2f7a88688650f2fc4fd99561383e72bba2305e3db59a1647  github.com/go-chi/chi/v5@v5.3.2/LICENSE
b343163c9d40108d37c7e5fe73e3c1d8f39c2acf26de5d1d227a1be2770205ad  github.com/go-chi/chi/v5@v5.3.2/Makefile
5e10ec94a4afe399d9ed4d9918d77a541fa30ea3bbc96086b6379d29e7af3a1d  github.com/go-chi/chi/v5@v5.3.2/README.md
c5f181d947996aebd888182fb8dccf6e7026690e701b5e5ca50b288b4a2cf5c3  github.com/go-chi/chi/v5@v5.3.2/SECURITY.md
6cf07102bf746267aeb610fb44c1a591ca2b4133e020b4c3428b02598cc089d8  github.com/go-chi/chi/v5@v5.3.2/_examples/README.md
1a7e8131f1ed48303dc42d93bc50ea330f08f0875cba64d0d03eb96754984127  github.com/go-chi/chi/v5@v5.3.2/_examples/chi.svg
e5a1f419488c233784f1f7951089305915b377af7462f137cd9bb5fd2a792924  github.com/go-chi/chi/v5@v5.3.2/_examples/custom-handler/main.go
de61c6c6cad7f28a9b8b506d32bba52485f6e902ae25954913e8480fa5cfd8a5  github.com/go-chi/chi/v5@v5.3.2/_examples/custom-method/main.go
9e343fdfdcb2892c5eb6906d405f9ac1f7cb1d73cd8b7f31cf6329719af96576  github.com/go-chi/chi/v5@v5.3.2/_examples/fileserver/data/notes.txt
0c0b96e8104689bc86312805d70a97aea7c870801f09ca03dd47c8edde31f620  github.com/go-chi/chi/v5@v5.3.2/_examples/fileserver/main.go
ad28151bf414776f787655d76fa42063a4174f95eeef6887a1724d9b3bfc32b5  github.com/go-chi/chi/v5@v5.3.2/_examples/graceful/main.go
e8104f868c9a673d864c489347970200808acc0652035a0376694bf5b7b07c86  github.com/go-chi/chi/v5@v5.3.2/_examples/hello-world/main.go
590c6d2f322481fed8abe1c2ca34eb0868f0d0a07d1c6aaccc0d6589b59a5d0d  github.com/go-chi/chi/v5@v5.3.2/_examples/limits/main.go
7b3cf7fb608cc030569d4bf75d30f0807ba0849543ee469ad5adadd13a3b8667  github.com/go-chi/chi/v5@v5.3.2/_examples/logging/main.go
11dd0741795c3c2f140b35b280a032a56241dea20c0e20b7cf191a0aedcdcac0  github.com/go-chi/chi/v5@v5.3.2/_examples/pathvalue/main.go
48c27bf0d7668fccfe7a9f6dac9199ee6240fe7a43855677c1f509e16732e2b6  github.com/go-chi/chi/v5@v5.3.2/_examples/router-walk/main.go
a0a551b6d6ca81210b129cd2a70d76e38370b2461be78813e6f117cfef6da020  github.com/go-chi/chi/v5@v5.3.2/_examples/todos-resource/main.go
5979a0124dbb9febcc1c6846ed843590d40f5ae440d2694e57c0860bf7bc395c  github.com/go-chi/chi/v5@v5.3.2/_examples/todos-resource/todos.go
212c60277ace9afb95ef02729be959c2c39472a14854d2af95c711c265e3778f  github.com/go-chi/chi/v5@v5.3.2/_examples/todos-resource/users.go
8c22d7bbc23f4b4d46ded5ef721a9c3a173031fa2c2c9a7b48c46c2b07e8bd80  github.com/go-chi/chi/v5@v5.3.2/chain.go
b659c130bc881c5cbaf03080497e3a25e6acbbed7c100ff7d19bf25f9bc64538  github.com/go-chi/chi/v5@v5.3.2/chi.go
bf1f093c84b276cc4790a23014f37abc92f85c0e4f5f6ed4830a233a29628916  github.com/go-chi/chi/v5@v5.3.2/context.go
74348e8085391a4f0a2af2f2e680debe8ffd7a5105d760451e7c9448dfd9e4a8  github.com/go-chi/chi/v5@v5.3.2/context_test.go
4d88d6917853bbb427146ba26661b51a21f8780f524162c9797a5769464e1482  github.com/go-chi/chi/v5@v5.3.2/go.mod
d1cfa0a32e9871cfba478590d1c1010297df20e53bbf83f34b09a4ec76e8b278  github.com/go-chi/chi/v5@v5.3.2/middleware/basic_auth.go
00cf47dd47cf4b6a1fa29dd3d5cdff4f7de8adbea410801128afe8aeaed35f0d  github.com/go-chi/chi/v5@v5.3.2/middleware/clean_path.go
68b1ba991406b3e0e8a86253143ce432554d2134e66a4a87417b19ecebe2ff8f  github.com/go-chi/chi/v5@v5.3.2/middleware/client_ip.go
6d432c0cb67290969057d899738970fb59c73b96bb9812f75d0a37dcc2a64946  github.com/go-chi/chi/v5@v5.3.2/middleware/client_ip_bench_test.go
4a00c5e73326b5be7985b99e433ef87ae90ec415dbb2a0a9cb366b81fb57b3b6  github.com/go-chi/chi/v5@v5.3.2/middleware/client_ip_example_test.go
bfc0d690534b833f73117d8a4c25aaedef6923c3145b6aadd3a0c149e2a427de  github.com/go-chi/chi/v5@v5.3.2/middleware/client_ip_test.go
493386222123e44826bb85ae88ed9c52b87d6d272ce8aacd949e13c42f7c552b  github.com/go-chi/chi/v5@v5.3.2/middleware/compress.go
49c8a68651c0542fde332dd035990226dd344d05b980e6591e46f0777a6571bb  github.com/go-chi/chi/v5@v5.3.2/middleware/compress_test.go
97fa5e37fc74e5bd2be57cb73d5a2abed89b198309056ffca6c3e819c5dde5d4  github.com/go-chi/chi/v5@v5.3.2/middleware/content_charset.go
39f7d5333c1f47a11796952ed1ec13aa6dd880c940069a2c49552c2936e2c1a1  github.com/go-chi/chi/v5@v5.3.2/middleware/content_charset_test.go
a1b35d124400a423d385eca9d066593c790e867db165af75e25b31b649d37fd8  github.com/go-chi/chi/v5@v5.3.2/middleware/content_encoding.go
98bf93fb59e1613843f8f838ddb17dc160cf88dd95135b1b2637e188c69c4644  github.com/go-chi/chi/v5@v5.3.2/middleware/content_encoding_test.go
5cedd359e3390edc613557681dda3cbc60faa2e4efb2a8974464876377b3a797  github.com/go-chi/chi/v5@v5.3.2/middleware/content_type.go
233f0c2f27f3d37719c6abc393c86804424a47d1ccebe17b018629b42f05a050  github.com/go-chi/chi/v5@v5.3.2/middleware/content_type_test.go
76af2014d2080f597a0205f51265badcad7af64d09652a79ed1a214074f5c582  github.com/go-chi/chi/v5@v5.3.2/middleware/get_head.go
a660972a6e389bbea07db1414dd72c0593809c644a941b431ce43c902ec65099  github.com/go-chi/chi/v5@v5.3.2/middleware/get_head_test.go
fc242f8c645cfc7640a36bccccb490cad9136ed07a751b0998f554412c29bca4  github.com/go-chi/chi/v5@v5.3.2/middleware/heartbeat.go
8e727c3d7630e6915f92bbf7e635bb5232e99bff73709151e7f17e9b5129b302  github.com/go-chi/chi/v5@v5.3.2/middleware/logger.go
f15a9900122e075c3b2d146aeae6128dfe17d57a45ef40b06f7a9e74cc4fb35f  github.com/go-chi/chi/v5@v5.3.2/middleware/logger_test.go
d600d7ea5b3184ea8ea952a145b1fca21e458b94601eae4a6cb215b3a04a3dee  github.com/go-chi/chi/v5@v5.3.2/middleware/maybe.go
8789f84814d815185cb5a7bc6c247969b67d0ff9b8045ee1c95f4c9ad59ff4b1  github.com/go-chi/chi/v5@v5.3.2/middleware/middleware.go
81744b82c9ea90b57a384c4a44cd889f7f5ab446b35eeb3d34983d0419a860bc  github.com/go-chi/chi/v5@v5.3.2/middleware/middleware_test.go
389eb1e308194561aea139b8d56d49f46b9fb30d84f8da146f5deaea66c5f5cc  github.com/go-chi/chi/v5@v5.3.2/middleware/nocache.go
03f3b50c4a55c345b52ce8f4b6d5da4b977647dd66e2e59b5f3635bda93e779c  github.com/go-chi/chi/v5@v5.3.2/middleware/page_route.go
0f3c01a40e320494a093d804376f3be8519b86285b5e24e8937bcee2c13c9279  github.com/go-chi/chi/v5@v5.3.2/middleware/path_rewrite.go
5b21bca97731b2df9aff87767ea1325b77a401c823fa6cdcf52ad38616d48faa  github.com/go-chi/chi/v5@v5.3.2/middleware/profiler.go
2f837f2213de7074c5c6915dad70f85c4588077291a8563d1cbad6698a757220  github.com/go-chi/chi/v5@v5.3.2/middleware/realip.go
6cab78b3ebc7fd581261afae6740423edc2199d9ba7ae3e0ffda99354f90d534  github.com/go-chi/chi/v5@v5.3.2/middleware/realip_test.go
472571093397d19475fd36ac1816a5c8e8de7706a9be9b155eb360657ca5f076  github.com/go-chi/chi/v5@v5.3.2/middleware/recoverer.go
1d7dd1b9f50370cd38ef5ff4738f12b811e304e1881175a1b1fea3d84fc4389f  github.com/go-chi/chi/v5@v5.3.2/middleware/recoverer_test.go
31b21034dd5cd6393fa9abc9ce1207acc78fea1e162d866edf6bb2a7badd288b  github.com/go-chi/chi/v5@v5.3.2/middleware/request_id.go
519e920890f2a33b77d0140a4a8ffcc933752bc04b3efb0af0793e11722b639e  github.com/go-chi/chi/v5@v5.3.2/middleware/request_id_test.go
81a98f72a81442a01a4353be9833fda6803f585bf768cebde5a34ebc1c44c390  github.com/go-chi/chi/v5@v5.3.2/middleware/request_size.go
cfd320f1a996ac95e29e25ab8a2ffa9606577996d6e001fb3835ede2795a3f69  github.com/go-chi/chi/v5@v5.3.2/middleware/route_headers.go
dd48799f1da69155927452801354ef413a97b1a83109fc8ad5f0560ae1f66906  github.com/go-chi/chi/v5@v5.3.2/middleware/route_headers_test.go
d5a0929e147b8b2a58ee78f0166a67f0949255b698bd5d26983a62643fc8fe10  github.com/go-chi/chi/v5@v5.3.2/middleware/strip.go
1002bd43aea540f4661806adcfc6ad21525e1e2c1e9fb92f47dda16fdaae0a4d  github.com/go-chi/chi/v5@v5.3.2/middleware/strip_test.go
16d99c4da576c4b7c024a12da139421d5659290250bed8f3d3ff6e0ed5d93b0e  github.com/go-chi/chi/v5@v5.3.2/middleware/sunset.go
ea640c109f86cc5d1fbceea558e3ab983d8f0de642c16fd2b4ad750efada0101  github.com/go-chi/chi/v5@v5.3.2/middleware/sunset_test.go
b3c6cb935d27aaddef0dcb5d165447babf5e9bb42cf28b4ee4e5f4693810888f  github.com/go-chi/chi/v5@v5.3.2/middleware/supress_notfound.go
da7433dde376ab7deef82ff5acad6bbbc64beca4dad95372468a5747e23dfaa3  github.com/go-chi/chi/v5@v5.3.2/middleware/terminal.go
05569797a14131ba177a71fccccd509363ca2016fa4c45a31cff033ae44301e4  github.com/go-chi/chi/v5@v5.3.2/middleware/throttle.go
a042c50fcbd5b72039329189aa3c1927636eb0ba3be3494dcf23e5ad4bd06ee2  github.com/go-chi/chi/v5@v5.3.2/middleware/throttle_test.go
3e83d007ec300b4ea605d97a336351e562ecd84b32c2b26206791d11ce7501c0  github.com/go-chi/chi/v5@v5.3.2/middleware/timeout.go
a1ea4ec9039b5a704f8aea7754a9b00d8026b72d18ba98208cd04b182e1b5800  github.com/go-chi/chi/v5@v5.3.2/middleware/url_format.go
04efa8cfbda922ff671ab3bdbc00f6fa06bf776c300dc8605815de1def33aaa1  github.com/go-chi/chi/v5@v5.3.2/middleware/url_format_test.go
2ccb524ee5ccd68885ff0d1de052068026f5c3298705fec930175e2ec06a1f91  github.com/go-chi/chi/v5@v5.3.2/middleware/value.go
bd6e76151cacaee5fdb11af48cf3c417f4587967cac8e7193a3cc0f7febe2f5b  github.com/go-chi/chi/v5@v5.3.2/middleware/wrap_writer.go
16392ba2eb524e838e2c3295dfa7a229339a478f85ac99466683074464435bce  github.com/go-chi/chi/v5@v5.3.2/middleware/wrap_writer_test.go
3255394e5c4f305b47b81e2e3912caaf58455223fe7a126cec5c494c663adb24  github.com/go-chi/chi/v5@v5.3.2/mux.go
c5240fec77d33920527c3676aeade48346fa1a0a79ccf2acc456cb673ebea012  github.com/go-chi/chi/v5@v5.3.2/mux_test.go
2c6f68e8540ceb9953869fbf57721ca6b03943c1903011bb0dcea202aa5d9fd8  github.com/go-chi/chi/v5@v5.3.2/path_value_test.go
2cb5166ea7db6a38ac3b67bf40ea2fd82393311967888a71e3b296592119d06e  github.com/go-chi/chi/v5@v5.3.2/pattern_test.go
8fb47f5e036bcb4daaeed180fa7164504a8ca43bdc010fd418618da48fdd340c  github.com/go-chi/chi/v5@v5.3.2/testdata/cert.pem
fecae623b5a0ae0a1110285c66e3e0e7f6019f49fd540ed3d9b85f65b9054f6c  github.com/go-chi/chi/v5@v5.3.2/testdata/key.pem
790a26ead9f2f9bb93e61f2c7ee588096068a26afcf56e115801753bd7ed230c  github.com/go-chi/chi/v5@v5.3.2/tree.go
47488f281c1d81ed2a0a68c29f9524cd48eba23e827cb5a7e2584d6eab8b7042  github.com/go-chi/chi/v5@v5.3.2/tree_test.go
```

「連結」といっても、1 ファイル 1 行のテキストファイル（86 行・10452 バイト）ができるだけです。これが「リスト」です。`LC_ALL=C` を付けるのは、Go の `slices.Sort` と同じバイト順で並べるためです。ロケール依存の順序で並べると、リストが変わってハッシュも変わります。`find .` ではなく `find github.com` にしているのは、書き出し中の `list.txt` 自身を拾わないためです。

**4. リスト全体の sha256 を取り、base64 して `h1:` を付ける**

```sh
shasum -a 256 list.txt
e58424202bd3092676e61a11b3225accdd2c723cca1a2bb854051cec7d68d676  list.txt
# リスト全体の sha256 を取得

shasum -a 256 list.txt | cut -d' ' -f1 | xxd -r -p | base64
5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
# リスト全体の sha256 の  「e58424202bd3092676e61a11b3225accdd2c723cca1a2bb854051cec7d68d676」 を取得
# xxd -r -p で 文字列に変換、それを base64 化する
```

base64 にかけるのは 64 文字の hex(16進数) 文字列ではなく、 hex が表す 32 バイトのバイト列です（`xxd -r -p` で hex からバイト列に戻しています）。頭に `h1:` を付けた `h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=` は、`go.sum` の 1 行目の値と一致します。

つまり**ファイル名と内容のリストのハッシュ**です。zip の圧縮レベル・タイムスタンプ・エントリ順序には依存しません。だから zip 内の日付が全部 `1980-00-00` に潰されていても問題になりません。

Go から同じ計算をするなら、`golang.org/x/mod/sumdb/dirhash` の 1 行で済みます。

```go
dirhash.HashZip(".../v5.3.2.zip", dirhash.Hash1)
// h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=  → go.sum と一致
```

`/go.mod` 行も同じ H1 ですが、対象はソースアーカイブではなく `v5.3.2.mod` 1 ファイルです。リストに書くファイル名は literal の `go.mod` です。

```sh
D=$(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v
printf '%s  go.mod\n' "$(shasum -a 256 $D/v5.3.2.mod | cut -d' ' -f1)" > modlist.txt

cat modlist.txt
4d88d6917853bbb427146ba26661b51a21f8780f524162c9797a5769464e1482  go.mod

shasum -a 256 modlist.txt | cut -d' ' -f1 | xxd -r -p | base64
R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
```

`h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=` は `go.sum` の 2 行目（`/go.mod h1:`）と一致します。名前を変えると値が変わるので、ここは実際に試すと分かりやすいです。

```txt
name = "github.com/go-chi/chi/v5@v5.3.2/go.mod" → h1:ZsvrnqFWiQ1WrgDxRZ...  一致しない
name = "go.mod"                                 → h1:R+tYY2hNuVUUjxoPtq...  go.sum と一致
```

### ステップ5: go.sum または checksum database と照合する

ここが一番中身のあるステップで、実際は 3 段構えです。

**(a) `go.sum` にすでに行がある場合**

それが絶対的な正解です。ローカルで計算した H1 と文字列比較して、違えば取得結果を捨ててエラーにします。

```txt
SECURITY ERROR: checksum mismatch
```

このとき checksum database には問い合わせません。`go.sum` が最優先です。

**(b) `go.sum` に無い場合**

`GOSUMDB`（デフォルト `sum.golang.org`）に問い合わせます。結果はキャッシュされていて、中身を直接読めます。

```sh
cat $(go env GOMODCACHE)/cache/download/sumdb/sum.golang.org/lookup/github.com/go-chi/chi/v5@v5.3.2
```

```txt
60295460
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=

go.sum database tree
63096908
8lHrQPRYe4KT4WQ9/Qk8/taKnD6IRHbFgfN7vQE2Y54=

— sum.golang.org Az3gri/oVBN1jAhNZ/iS+W31psFAWRT20/GyrhpU6taFwulQGetBkH5S7wpGXdS935jiIGiUUAXBVHzT/P6fYhqAmAk=
```

行ごとの意味はこうです。

| 内容 | 意味 |
|---|---|
| `60295460` | このモジュールのレコード番号（追記専用ログの何番目か） |
| `h1:` の 2 行 | checksum database が記録している正解のハッシュ |
| `63096908` | 署名時点のツリー全体のレコード数 |
| `8lHrQPR...` | Merkle ツリーのルートハッシュ |
| `— sum.golang.org ...` | ed25519 署名（signed note 形式） |

重要なのは、**この署名の検証鍵が go コマンドのバイナリに埋め込まれている**ことです。だから `proxy.golang.org` を信用する必要がありません。プロキシ経由で checksum database の応答を中継してもらっても（`GET /sumdb/sum.golang.org/lookup/...`）、署名が合わなければ弾かれます。

さらに「レコード `60295460` が本当にルートハッシュ `8lHrQPR...` のツリーに含まれるか」を Merkle の包含証明で確認します。そのために取得する断片が tile です。

```sh
ls $(go env GOMODCACHE)/cache/download/sumdb/sum.golang.org/tile/8/0/
013  x002  x040  x076  x111  x119  x123  ...
```

`8` は tile の高さ（2^8 = 256 リーフ単位で切り出す）、`x` 付きは未完成の部分 tile です。この仕組みのおかげで、6000 万件あるログ全体をダウンロードせずに包含を検証できます。Certificate Transparency と同じ発想の透明性ログで、「checksum database が特定の人にだけ嘘のハッシュを返す」ことが検出可能になっています。

**(c) 2 回目以降**

`.ziphash` があるので、ネットワークも再ハッシュも不要です。ローカルのファイル比較で終わります。

```sh
cat $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.ziphash
h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
```

なお `GOPRIVATE` / `GONOSUMDB` に一致するモジュールや `GOSUMDB=off` の場合は (b) をスキップします。プライベートリポジトリのパスを公開 checksum database に送らないための逃げ道です。

### ステップ6: go.sum に記録する

検証を通ったハッシュを `go.sum` に追記し、モジュールパスとバージョンの順にソートして書き出します。書き込みタイミングは「解決が完了して `go.mod` を更新するとき」で、`-mod=readonly`（Go 1.16 以降のデフォルト）だと勝手には書かず、代わりにエラーで `go mod tidy` を促します。

記録される行数にはルールがあります。

- **ビルドに実際にソースが必要** → `h1:`（zip）と `/go.mod` の**2 行**
- **モジュールグラフの構築にしか使わない** → `/go.mod` の**1 行**だけ

そして今回のケースで一番面白いのはここです。

```sh
cat go.sum
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
```

ステップ1 で取得された `github.com/go-chi/chi v1.5.5` は、**1 行も載っていません**。キャッシュには zip も展開済みツリーも残っているのに、です。理由は、あれが import path 解決のための「探り」であって、最終的なモジュールグラフの一員ではないからです。

```txt
モジュールキャッシュ ⊃ go.sum
```

キャッシュは「見に行ったもの全部」、`go.sum` は「ビルドに関与するものだけ」という非対称があります。

### 6 ステップの実体まとめ

```txt
1. MVS + import path の各接頭辞への最新バージョン要求（→ 余分な chi v1.5.5 の取得）
2. proxy の 5 エンドポイントを HTTP GET、または git を直叩き
3. cache/download に生データ、<mod>@<ver>/ に読み取り専用の展開ツリー
4. H1 = sha256(「sha256 hex + ファイル名」の行をソートして連結したリスト)
5. go.sum が最優先。無ければ checksum database に問い合わせ、
   ed25519 署名 + Merkle 包含証明を検証（検証鍵は go バイナリ内）
6. ソースが要るものは 2 行、グラフ用途だけなら /go.mod の 1 行
```

`go mod verify` は展開済みツリーを再ハッシュして `go.sum` と突き合わせる、つまりステップ4・5 の再実行なので、キャッシュ汚染の確認に使えます。

### go.sum の 2 行はそれぞれ何か

あらためて、今回生成された go.sum を見ます。

```txt
github.com/go-chi/chi/v5 v5.3.2 h1:...
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:...
```

それぞれ意味が違います。

```txt
github.com/go-chi/chi/v5 v5.3.2 h1:...
```

これは **モジュール本体の内容**のハッシュです。

```txt
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:...
```

これは **そのモジュールの go.mod ファイル**のハッシュです。

大事なのは、`go.sum` は `package-lock.json` みたいな「依存バージョンを固定するロックファイル」ではないことです。Go で実際に使うバージョンを決める中心は `go.mod` です。`go.sum` は「そのバージョンの中身が前と同じか」を検証するための証拠ファイルです。

なので一言でいうと、

> `go.sum` は、Go コマンドが外部モジュールを取得したときに、その中身が改ざんされていないことを後で確認できるよう自動生成するチェックサム台帳。

です。


## 疑問5

> ダウンロードされた実体はどこにある？ なぜ `node_modules` のようにプロジェクト直下に来ないのか？
>
> 追加の疑問: なぜ git 管理下のリポジトリの中にパッケージを置かず、`GOMODCACHE` の中に置くようにしたのか？ 別のリポジトリで違うバージョンを使っていたら、管理が大変にならないのか？

### まず実体の場所を確認する

```bash
go env GOMODCACHE
```

```txt
/Users/shiraki_takuma/go/pkg/mod
```

プロジェクト直下 (`GoBash/ex-Go/`) ではなく、ホーム配下の共有ディレクトリです。`go list` に実際のビルド対象ディレクトリを聞くと、はっきりします。

```bash
go list -f '{{.ImportPath}} {{.Dir}}' github.com/go-chi/chi/v5
```

```txt
github.com/go-chi/chi/v5 /Users/shiraki_takuma/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2
```

つまり `main.go` の `import "github.com/go-chi/chi/v5"` がコンパイル時に読んでいるのは、リポジトリの外にあるこのディレクトリです。

公式もこの 2 段構えをそのまま書いています。

> If needed, it downloads module source code so you can compile packages that depend on them. It can download modules from a module proxy like proxy.golang.org or directly from version control repositories. **The source is cached locally.**
> — [Managing dependencies](https://go.dev/doc/modules/managing-dependencies)

> The module cache is the directory where the `go` command stores downloaded module files. (...) The default location of the module cache is `$GOPATH/pkg/mod`. To use a different location, set the `GOMODCACHE` environment variable.
> — [Go Modules Reference / Module cache](https://go.dev/ref/mod#module-cache)

`$GOPATH/pkg/mod` がデフォルトで、`GOMODCACHE` で上書きできる。手元の `GOPATH` は `/Users/shiraki_takuma/go` なので、既定値そのままです。

```bash
go env GOPATH GOMODCACHE
```

```txt
/Users/shiraki_takuma/go
/Users/shiraki_takuma/go/pkg/mod
```

キャッシュが `cache/download/`（生データ）と `<module>@<version>/`（展開済み）の 2 層になっている点は疑問4のステップ3で見た通りです。公式のファイル一覧でも `$module@$version/` は「Directory containing extracted contents of a module `.zip` file」と定義されています。

### なぜリポジトリの中ではなく共有キャッシュなのか

公式の記述を根拠に分解すると、理由は 4 つに整理できます。

#### 1. 「共有してよい」と明示されている

> **The cache may be shared by multiple Go projects developed on the same machine. The `go` command will use the same cache regardless of the location of the main module.** Multiple instances of the `go` command may safely access the same module cache at the same time.
> — [Module cache](https://go.dev/ref/mod#module-cache)

「main module がどこにあってもキャッシュは同じ」と言い切っています。プロジェクトごとにコピーを持つ設計ではない、というのが仕様レベルの宣言です。並行実行しても安全、とも書かれています（内部的にはロックファイル。疑問4で見た `v5.3.2.lock` がそれ）。

#### 2. バージョンが不変（immutable）だから共有できる

> A version identifies **an immutable snapshot of a module**, which may be either a release or a pre-release.
> — [Go Modules Reference / Versions](https://go.dev/ref/mod#versions)

`github.com/go-chi/chi/v5@v5.3.2` は、世界中のどのマシンでも同じ中身であることが前提になっています。中身が同じと保証できるものを、プロジェクトの数だけコピーする理由がありません。

さらにそれを裏で担保しているのが checksum database です。

> It also ensures that **the bits associated with a specific version do not change from one day to the next, even if the module's author subsequently alters the tags in their repository.**
> — [Checksum database](https://go.dev/ref/mod#checksum-database)

「作者が後からタグを付け替えても、そのバージョンのビット列は変わらない」。npm の `unpublish` 騒動（left-pad）と対照的で、ここが**共有キャッシュを成立させている土台**です。

#### 3. ディレクトリ名にバージョンが入るので、別バージョンは衝突しない

「別のリポジトリで違うバージョンを使っていたら大変では？」という疑問への直接の答えがここです。キャッシュのパスは `$module@$version/` なので、**同じモジュールの別バージョンは別ディレクトリとして共存します**。

実際に手元のキャッシュを見ると、すでに共存していました。

```bash
ls -d $(go env GOMODCACHE)/golang.org/x/mod@* $(go env GOMODCACHE)/golang.org/x/tools@*
```

```txt
/Users/shiraki_takuma/go/pkg/mod/golang.org/x/mod@v0.22.0
/Users/shiraki_takuma/go/pkg/mod/golang.org/x/mod@v0.37.0
/Users/shiraki_takuma/go/pkg/mod/golang.org/x/tools@v0.47.1-0.20260707181000-a299dadba899
/Users/shiraki_takuma/go/pkg/mod/golang.org/x/tools@v0.49.0
```

`golang.org/x/mod` が v0.22.0 と v0.37.0、`golang.org/x/tools` が 2 バージョン。別々のプロジェクトが別々のバージョンを要求した結果が、そのまま並んでいます。

```txt
リポジトリA の go.mod → x/mod v0.22.0 ─┐
                                       ├→ 同じキャッシュの別ディレクトリを読む
リポジトリB の go.mod → x/mod v0.37.0 ─┘
```

なので「管理が大変になる」方向には働きません。**どのバージョンを使うかを決めるのは各リポジトリの `go.mod`** であって、キャッシュは「決まったバージョンの実体を置いておく場所」でしかない。キャッシュ側に競合という概念がありません（疑問9 の「go.mod は宣言、go.sum は証明」と同じ分業）。

ついでに言うと、同じバージョンを 10 個のリポジトリで使っていてもディスク上の実体は 1 つです。手元のキャッシュ全体はこのサイズでした。

```bash
du -sh $(go env GOMODCACHE)
```

```txt
166M	/Users/shiraki_takuma/go/pkg/mod
```

これは**マシン全体の全 Go プロジェクト分の合計**です。npm はプロジェクトごとに `node_modules` を持つので、リポジトリ 10 個なら 10 回展開されます。

#### 4. 共有しても安全な理由 —— 検証は各プロジェクトの go.sum が握っている

ここが一番効く根拠です。共有キャッシュの最大の懸念は「他のプロジェクトが汚したものを掴まされないか」ですが、公式はそれを名指しで潰しています。

> **The module cache is usually shared by all Go projects on a system, and each module may have its own `go.sum` file with potentially different hashes. To avoid the need to trust other modules, the `go` command verifies hashes using the main module's `go.sum` whenever it accesses a file in the module cache.** Zip file hashes are expensive to compute, so the `go` command checks pre-computed hashes stored alongside zip files instead of re-hashing the files.
> — [Authenticating modules](https://go.dev/ref/mod#authenticating)

読み解くとこうです。

```txt
キャッシュ  … マシンで 1 つ。誰が入れたか分からないものが入っている
go.sum      … リポジトリごと。git 管理下。「自分が認めたハッシュ」の台帳
              ↓
キャッシュ内のファイルにアクセスする「たび」に
main module の go.sum と突き合わせる
              ↓
他のプロジェクトを信頼する必要がない (avoid the need to trust other modules)
```

つまり **共有しているのは「実体」だけで、「信頼」は共有していない**。信頼の単位は git 管理下の `go.sum` のままです。疑問4のステップ5で見た `.ziphash` は、この「毎回の照合」を毎回 zip を再ハッシュせずに済ませるための事前計算値でした。

そして書き換え防止は、ハッシュだけでなくパーミッションでも二重にかかっています。

> The `go` command creates module source files and directories in the cache with **read-only permissions** to prevent accidental changes to modules after they're downloaded. This has the unfortunate side-effect of making the cache difficult to delete with commands like `rm -rf`. The cache may instead be deleted with `go clean -modcache`.
> — [Module cache](https://go.dev/ref/mod#module-cache)

実際に書き換えを試すと弾かれます。

```bash
M=$(go env GOMODCACHE)/github.com/go-chi/chi/v5@v5.3.2
ls -l $M/mux.go
echo "// hack" >> $M/mux.go
touch $M/evil.go
```

```txt
-r--r--r--@ 1 shiraki_takuma  staff  16981 ... mux.go

zsh: permission denied: .../chi/v5@v5.3.2/mux.go
touch: .../chi/v5@v5.3.2/evil.go: Permission denied
```

```bash
go mod verify
```

```txt
all modules verified
```

`r--r--r--`。共有ディレクトリなのに「あるプロジェクトが直接書き換えて他のプロジェクトを壊す」ことが、そもそもできない設計です。`node_modules` を手で書き換えてデバッグする、という運用が Go に無いのはこれが理由。

### 「それでもリポジトリに同梱したい」場合の公式の逃げ道

Go はその選択肢を捨てたわけではなく、opt-in として残しています。

> Vendoring may be used to allow interoperation with older versions of Go, or **to ensure that all files used for a build are stored in a single file tree.** (...) When vendoring is enabled, build commands like `go build` and `go test` load packages from the `vendor` directory **instead of accessing the network or the local module cache**.
> — [Vendoring](https://go.dev/ref/mod#vendoring)

`go mod vendor` を打てば `vendor/` に実体がコピーされ、リポジトリに同梱できる（= `node_modules` 相当の配置になる）。デフォルトではない、というだけです。

```txt
デフォルト   : go.mod（宣言） + go.sum（証明） を git 管理 / 実体は共有キャッシュ
vendor モード: 実体も git 管理（単一ファイルツリーに閉じる）
```

### npm との対比

| | Go (module cache) | npm (node_modules) |
| --- | --- | --- |
| 実体の場所 | `$GOMODCACHE` = マシンで 1 つ | プロジェクト直下に毎回展開 |
| 同一バージョンの重複 | 無い（`<mod>@<ver>` で 1 つ） | プロジェクト数だけコピー |
| 別バージョンの共存 | ディレクトリ名にバージョンが入るので自然に共存 | ネストや hoisting で解決 |
| 書き込み | 読み取り専用 (`r--r--r--`) | 書き換え可能 |
| git 管理 | `go.mod` / `go.sum` だけ | 原則 `.gitignore`（同梱する流派もある） |
| 安全性の担保 | アクセスのたびに main module の `go.sum` と照合 | install 時中心 |
| 削除 | `go clean -modcache`（**全プロジェクト共有**なので注意） | `rm -rf node_modules`（そのプロジェクトのみ） |

### まとめ

```txt
なぜリポジトリの中に置かないのか
  ├ バージョンは immutable → 同じ物を人数分コピーする意味がない
  ├ パスが <module>@<version> → 別バージョンは自然に共存、競合しない
  ├ 使うバージョンを決めるのは各リポジトリの go.mod → キャッシュは実体の置き場でしかない
  └ 共有しても安全 → アクセスのたびに main module の go.sum で検証する
                    ＝ 実体は共有、信頼は共有しない
```

「共有キャッシュにしたから管理が楽になる」のではなく、**バージョンを immutable にして go.sum で毎回検証する設計にしたから、共有キャッシュにできた**、という順序で捉えるのが正確です。

次に湧く疑問: ではその「実体」は最初にどこから、どうやって取ってきたのか？ `git clone` しているのか？（→ 疑問6）
