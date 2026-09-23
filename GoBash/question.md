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

`github.com/go-chi/chi/v5` という import path を見ても、go コマンドは「どこまでがモジュール名で、どこからがパッケージのサブディレクトリか」を知りません。なので接頭辞を長い順に総当たりします。

```txt
github.com/go-chi/chi/v5   ← これがモジュール？
github.com/go-chi/chi      ← モジュールで、中に v5/ パッケージがある？
github.com/go-chi          ← ?
```

これが、先ほどのログに出てきた謎の 1 行の正体です。

```sh
go: downloading github.com/go-chi/chi/v5 v5.3.2
go: downloading github.com/go-chi/chi v1.5.5   ← これ
```

`github.com/go-chi/chi` の中に `v5` というパッケージディレクトリがあるかもしれないので、確認のために取得しています。実際、キャッシュには両方残っています。

```sh
ls $(go env GOMODCACHE)/github.com/go-chi/
chi   chi@v1.5.5
```

**(b) MVS（Minimal Version Selection）**

候補が確定したら、モジュールグラフを構築してバージョンを決めます。「各モジュールについて、要求された中で**最大**のバージョンを選ぶ」だけのアルゴリズムです。npm のような SAT ソルバ的な解決はしません。

ここで大事なのは、グラフの構築には**依存モジュールの `go.mod` だけあればよい**ということです。ソース全体は要りません。だから後述の proxy プロトコルは `.mod` と `.zip` を別エンドポイントに分けていて、`go.sum` にも 2 行載る構造になっています。

### ステップ2: GOPROXY または VCS から取得する

デフォルトの設定はこうなっています。

```sh
go env GOPROXY
https://proxy.golang.org,direct
```

カンマ区切りはフォールバックのチェーンです。`direct` は「プロキシを使わず直接 VCS を叩く」という特殊な値で、プロキシが 404 / 410 を返したら次に進みます。

プロキシへのリクエストは、この 5 種類の HTTP GET だけです。

```txt
GET /<module>/@v/list              バージョン一覧
GET /<module>/@v/<version>.info    メタデータ (JSON)
GET /<module>/@v/<version>.mod     そのバージョンの go.mod
GET /<module>/@v/<version>.zip     ソース本体
GET /<module>/@latest              最新版の解決
```

バージョン無指定の `go get` だと `@latest` → `.info` → `.mod` → `.zip` の順に叩かれます。

`.info` には取得元の証跡が残っています。

```sh
cat $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.info
```

```json
{"Version":"v5.3.2","Time":"2026-08-20T09:37:52Z",
 "Origin":{"VCS":"git","URL":"https://github.com/go-chi/chi",
 "Hash":"38939062c5df4d3e8814aad1a488983112627ced","Ref":"refs/tags/v5.3.2"}}
```

プロキシ自身が「git の `refs/tags/v5.3.2`、コミット `3893906` から作った zip だ」と申告しているわけです。`direct` の場合はここを go コマンド自身がやります（`git ls-remote` でタグ一覧を取り、該当コミットを取得して zip を自前で組み立てる）。

社内モジュールをプロキシに漏らしたくない場合の制御はこのあたりです。

```sh
GOPRIVATE   GONOPROXY   GONOSUMDB   GOVCS
```

### ステップ3: モジュールキャッシュに保存する

保存先は `GOMODCACHE` で、キャッシュは 2 層構造になっています。

```txt
$GOMODCACHE/cache/download/...   ← ダウンロードした生データ（zip のまま）
$GOMODCACHE/<module>@<version>/  ← 展開済みツリー（ビルドが読む）
```

生データ側を覗くとこうなっています。

```sh
ls $(go env GOMODCACHE)/cache/download/github.com/go-chi/chi/v5/@v/
list  v5.3.2.info  v5.3.2.lock  v5.3.2.mod  v5.3.2.zip  v5.3.2.ziphash
```

| ファイル | 役割 |
|---|---|
| `.info` | バージョンと取得元のメタデータ |
| `.mod` | そのバージョンの `go.mod` 単体 |
| `.zip` | ソース本体 |
| `.ziphash` | 検証済みハッシュのキャッシュ（ステップ5で効く） |
| `.lock` | 並行する `go` プロセス同士の排他用（0 バイト） |

パスには 2 つ仕掛けがあります。

ひとつは**大文字のエスケープ**です。大文字小文字を区別しないファイルシステム（macOS がまさにそれ）で衝突しないよう、大文字は `!` + 小文字に変換されます。

```txt
github.com/!burnt!sushi/   ← github.com/BurntSushi/
```

もうひとつは**メジャーバージョンサフィックス**です。`/v5` はモジュールパスの一部なので、ディレクトリ階層にそのまま現れます（`chi/v5/@v/`）。

展開済みツリー側は**読み取り専用**で作られます。

```sh
ls -ld $(go env GOMODCACHE)/github.com/go-chi/chi/v5@v5.3.2
dr-xr-xr-x@ 24 ... github.com/go-chi/chi/v5@v5.3.2
```

`r-x` です。依存先をうっかり編集してハッシュ検証を壊すことを、パーミッションで物理的に防いでいます。消すときに `go clean -modcache` が必要なのもこれが理由です。

### ステップ4: 暗号学的ハッシュを計算する

`h1:` の `1` はアルゴリズムのバージョン番号です。中身で大事なのは、**zip のバイト列の SHA-256 ではない**という点です。

```txt
1. アーカイブ内の全ファイルについて sha256 を計算する
2. "<sha256 の hex>  <ファイル名>\n" という行を作る（スペースは 2 個）
3. ファイル名でソートして全行を連結する → これが「リスト」
4. そのリスト全体の sha256 を取り、base64 して "h1:" を付ける
```

つまり**ファイル名と内容のリストのハッシュ**です。zip の圧縮レベル・タイムスタンプ・エントリ順序には依存しません。だから zip 内の日付が全部 `1980-00-00` に潰されていても問題になりません。

ハッシュ対象の名前は zip 内のフルパスです。

```txt
github.com/go-chi/chi/v5@v5.3.2/.gitignore
github.com/go-chi/chi/v5@v5.3.2/CHANGELOG.md
...
```

`golang.org/x/mod/sumdb/dirhash` で再計算すると、`go.sum` の値と一致します。

```go
dirhash.HashZip(".../v5.3.2.zip", dirhash.Hash1)
// h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=  → go.sum と一致
```

`/go.mod` 行も同じ H1 ですが、**ファイル名が literal `go.mod` 1 個だけ**のリストを対象にします。名前を変えると値が変わるので、ここは実際に試すと分かりやすいです。

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
1. MVS + import path の接頭辞の総当たり（→ 余分な chi v1.5.5 の取得）
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
