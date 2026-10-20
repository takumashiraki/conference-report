前回の記事「npm との違いを手がかりに Go の依存管理を理解したい」の最後に、登壇までに調べきれなかった疑問を 4 つ挙げました。
この記事では、その 4 つを実際に調べた結果をまとめます。

<!-- TODO: 前回記事の URL を貼る -->

## この記事で扱うこと

前回の記事では、`go get` の 6 ステップを追って「`go.sum` は lock ファイルではなく、取ってきた中身が本物かを確かめる台帳だ」という結論を出しました。残った疑問は次の 4 つです。

1. 最後は Google（`sum.golang.org` の運営者）を信用するしかないのか
2. Go 1.27.0 で依存管理がどう変わったのか
3. dep（`Gopkg.lock`）の時代の依存管理と、なぜ lock ファイルをやめたのか
4. 標準パッケージが `go.mod` と `go.sum` に記録されない理由（どうやってインストールし、どうやって検証しているのか）

この記事では順番を入れ替え、3 → 4 → 1 → 2 の順に進みます。lock ファイルをやめた経緯を知ると、`go.sum` と checksum database の役割がはっきりします。標準パッケージの検証を追うと、最後は checksum database にたどり着きます。そこで「Google を信用するしかないのか」を考え、その途中で見つかった Go 1.27 の修正につなげます。

確認した環境は次のとおりです。

| 対象 | バージョン |
| --- | --- |
| Go | 1.27.0（比較用に 1.26.8 / 1.26.2 / 1.26.0 / 1.25.3 も使用） |
| 依存に使うモジュール | `github.com/go-chi/chi/v5` v5.3.2 / `golang.org/x/text` v0.37.0 |

## 疑問 3：dep の時代と、lock ファイルをやめた理由

### GOPATH の時代は「いつも最新」だった

モジュールが入る前の `go get` は、バージョンという考え方を持っていませんでした。Russ Cox は vgo の連載の 1 本目で、当時の `go get` をこう説明しています。

> When go get needs a package, it always fetches the latest copy, delegating the download and update operations to version control systems like Git or Mercurial.
> — [Go += Package Versioning](https://research.swtch.com/vgo-intro)

同じコードをビルドしても、取得した日によって依存の中身が変わります。この問題を解決するために、コミュニティの道具と公式の仕組みが順に登場しました。

| 時期 | 出来事 |
| --- | --- |
| 2013 年 9 月 | godep（依存を固定するツール）が登場 |
| 2015 年 8 月 | Go 1.5 で vendor experiment（`GO15VENDOREXPERIMENT`） |
| 2016 年 2 月 | Go 1.6 で vendor が既定で有効に |
| 2017 年 1 月 | dep が「official experiment」として公開 |
| 2018 年 2 月 | Russ Cox が vgo（のちの Go Modules）を提案 |
| 2018 年 8 月 | Go 1.11 で modules を試験導入 |
| 2019 年 9 月 | Go 1.13 で `proxy.golang.org` と `sum.golang.org` が既定に |
| 2020 年 | dep がアーカイブされる |
| 2021 年 2 月 | Go 1.16 で `GO111MODULE=on` が既定に |

### dep は「範囲」と「lock」の 2 ファイルだった

dep の構成は、前回の記事で見た npm とよく似ています。手で書く制約ファイル `Gopkg.toml` と、ツールが書く `Gopkg.lock` の 2 つです。

```toml:Gopkg.toml
[[constraint]]
  name = "github.com/Masterminds/vcs"
  version = "1.11.0"

[[constraint]]
  name = "github.com/Masterminds/semver"
  branch = "2.x"
```

`version = "1.11.0"` は、ちょうど 1.11.0 という意味ではありません。dep のドキュメントには、演算子を書かなければ `^` が付くとあります。つまり「1.11.0 以上 2.0.0 未満」という **範囲** です。npm の `^5.2.1` と同じ考え方です。

> When you specify a version without an operator, dep automatically uses the ^ operator by default.
> — [Gopkg.toml](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md)

範囲のどれを選んだかは `Gopkg.lock` に残ります。

```toml:Gopkg.lock
[[projects]]
  digest = "1:442020d2..."
  name = "github.com/Masterminds/vcs"
  packages = ["."]
  pruneopts = "NUT"
  revision = "3084677c2c188840777bff30054f2b553729d329"
  version = "v1.11.1"
```

ドキュメントは `Gopkg.lock` を "a transitively complete snapshot of a project's dependency graph" と説明しています。間接依存まで含めて、使うコミットをすべて書き留めたファイルです。

範囲から 1 つを選ぶ作業には、dep の中の `gps` パッケージにある CDCL 型のソルバー（SAT ソルバーと同じ系統の探索）を使っていました。

### lock をやめた理由 1：「許される中で最新のバージョンを選ぶ」方針をやめた

Go Modules が lock ファイルを持たない理由は、前回の記事で見た MVS にあります。Russ Cox は、選び方の既定そのものが間違っていたと書いています。

> Nearly all package managers today, including dep and cargo, use the newest allowed version of packages involved in the build. I believe this is the wrong default
> — [Go += Package Versioning](https://research.swtch.com/vgo-intro)

「許される中で最新」を選ぶ限り、結果はその日に公開されているバージョンで変わります。だから、選んだ結果を lock に残す必要がありました。MVS は逆に「要求された中で最も古いもの」を選ぶので、新しいバージョンが公開されても結果が変わりません。

> Minimal version selection delivers reproducible builds by default, without a lock file.

### lock をやめた理由 2：ライブラリの lock は読まれない

lock ファイルが読まれるのは、最終的にビルドするプログラムの lock ファイルだけです。

> But lock files only apply to whole programs, to package main.
> (中略)
> So Dep and the others must ignore lock files associated with libraries and fall back to the default time-based decisions.
> — [The Principles of Versioning in Go](https://research.swtch.com/vgo-principles)

ライブラリ A に `Gopkg.lock` があっても、A を使うアプリをビルドするときには A の lock は読まれません。
A の依存は、アプリをビルドした日の「許される中で最新のバージョン」で選び直されます。**A の作者がテストした組み合わせと、利用者がビルドする組み合わせが一致しません。**

MVS は各モジュールの `go.mod` に書いた最低バージョンだけを使って選びます。そのため、ライブラリでもアプリでも同じ結果を再現できます。

> Minimal version selection provides repeatability, for whole programs and for libraries, always, without any lock files.

### lock をやめた理由 3：バージョン選択は NP 完全だった

3 つ目は計算量です。範囲で制約を書くと、バージョン選択は NP 完全な問題になります。

> Unfortunately, the version selection problem is NP-complete
> — [Version SAT](https://research.swtch.com/version-sat)

:::details **NP 完全** とは

ざっくり言うと「答えが正しいかどうかはすぐ確かめられるのに、答えを速く見つける方法が知られていない問題」のことです。
組み合わせを片っ端から試すのに近いやり方しかなく、候補が増えると時間が爆発的に増えます。

バージョン選択に当てはめると、次のようになります。

```text
アプリ ─→ A（v1 か v2 を選べる）
      └→ B（v1 か v2 を選べる）

A v1 は C v1 が必要    A v2 は C v2 が必要
B v1 は C v2 が必要    B v2 は C v1 が必要
C は 1 つのバージョンしか入れられない
```

「A v1・B v2・C v1」という組み合わせを渡されれば、条件を満たしているかは 1 つずつ見ればすぐ分かります。ところが、条件を満たす組み合わせを最初から探す場合は、A と B の選び方を順に試すしかありません。選べるパッケージが 2 つなら 4 通りですが、20 個なら約 100 万通り、30 個なら約 10 億通りです。dep がこの探索に、論理式を解く SAT ソルバーと同じ系統の手法を使っていたのはこのためです。

Russ Cox の記事は、この難しさが次の前提から生まれると示しています。

- 依存に特定のバージョン（または範囲）を指定できる
- 同じパッケージの 2 つのバージョンは同時に入れられない

:::

Go Modules は、この 2 つの前提を両方とも崩しました。`require` に書けるのは「最低バージョン」の 1 種類だけです。互換性のないメジャーバージョンは `/v5` のように別のモジュールパス（前回の記事）になるので、同時に入れられます。

制約が「〜以上」だけなら、各モジュールについて要求された中で最も高いバージョンを取れば、すべての要求を同時に満たせます。試行錯誤は要りません。依存の数に比例した時間（線形時間）で答えが出て、manifest と lock を分ける必要もなくなりました。Go 公式ブログはこう書いています。

> The constraints of Cargo and Dep make version selection equivalent to solving Boolean satisfiability
> (中略)
> a trivial, linear-time algorithm … in turn eliminates the need for separate lock and manifest files.
> — [A Proposal for Package Versioning in Go](https://go.dev/blog/versioning-proposal)

:::details 反対意見もあった

この決定は、すんなり受け入れられたわけではありません。dep の中心にいた Sam Boyer は、提案の issue に強い反対を書いています。

> i believe that MVS is unfit for purpose, and should not make it to production.
> — [golang/go#24301](https://github.com/golang/go/issues/24301#issuecomment-388732254)

同じコメントには「間接依存まで含めた lock ファイルを別に足すなら、このまま進めてもよい」という妥協案もありました。ほかに、「古いバージョンを選ぶので、セキュリティ修正が自動では入らない」という批判も出ています。Russ Cox は「更新は開発者が求めたときだけ行う」のが正しいと答え、提案は 2018 年 5 月に採択されました。

:::

### Go 1.17 で、go.mod は lock に近づいた

ここで 1 つ補足があります。Go 1.17 以降、`go.mod` には間接依存も `// indirect` 付きですべて並ぶようになりました（[Module graph pruning](https://go.dev/ref/mod#graph-pruning)）。

> At go 1.17 or higher: The go.mod file includes an explicit require directive for each module that provides any package transitively imported by a package or test in the main module.

この変更の設計文書でも、`go.mod` が lock ファイルに似てくることを認めています。

> The set of recorded modules would more closely resemble a 'lock' file as used in other dependency-management systems. (However, the go command would still not require a separate 'manifest' file, and unlike a lock file, the go.mod file would still be updated automatically to reflect new requirements discovered during package loading.)
> — [Lazy Module Loading](https://go.googlesource.com/proposal/+/master/design/36460-lazy-module-loading.md)

### 疑問 3 の回答

- dep は npm と同じ「範囲 + lock」の 2 ファイル構成だった
- Go Modules は選び方を「最新」から「最低」に変えたので、結果を lock に残す必要がなくなった。ライブラリでも同じ結果を再現でき、計算も線形時間で済む
- いまの `go.mod` は間接依存まで並ぶので、見た目は lock に近い。ただし手で書く宣言と自動で書く結果が 1 ファイルにまとまっていて、中身の検証は `go.sum` が受け持つ

## 疑問 4：標準パッケージは、なぜ go.mod と go.sum に載らないのか

前回の記事の `main.go` は `fmt` と `net/http` も import していました。それなのに、`go.mod` にも `go.sum` にも出てきません。

### 標準パッケージはモジュールに属さない

`go list` で、標準パッケージがどのモジュールに属するかを見てみます。

```bash
$ go list -f '{{.ImportPath}} std={{.Standard}} module={{.Module}} dir={{.Dir}}' net/http fmt
net/http std=true module=<nil> dir=~/.local/share/mise/installs/go/1.27.0/src/net/http
fmt std=true module=<nil> dir=~/.local/share/mise/installs/go/1.27.0/src/fmt
```

`module=<nil>` です。標準パッケージは利用者から見るとどのモジュールにも属さず、実体は `GOROOT` の下にあります。import path を見ただけで標準パッケージだと分かるように、ルールも決まっています。

> The standard library reserves all import paths without a dot in the first element for its packages.
> — `go help importpath`

`github.com/...` のように最初の要素にドットがあればモジュール、`net/http` のように無ければ標準パッケージです。前回の記事で見た「import path の末尾から削ってプロキシに問い合わせる」処理は、標準パッケージでは行われません。

### 標準ライブラリにも go.mod と go.sum がある

「標準ライブラリには `go.mod` が無い」と思っていたのですが、`GOROOT` の中を見るとありました。

```bash
$ cat $(go env GOROOT)/src/go.mod
module std

go 1.27

require (
	golang.org/x/crypto v0.52.1-0.20260526024921-9beb694f9766
	golang.org/x/net v0.55.1-0.20260731170536-c1d18010be90
)

require (
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

$ head -2 $(go env GOROOT)/src/go.sum
golang.org/x/crypto v0.52.1-0.20260526024921-9beb694f9766 h1:ABD+jVg0H4Hwu2sGcUtKeb3T8mlS+jS3uWrkTAPcXjs=
golang.org/x/crypto v0.52.1-0.20260526024921-9beb694f9766/go.mod h1:1QgfPxDqh0T2M/elOJtp9RvuR95kVjir0e6/BvEmGbc=
```

標準ライブラリは `std` という名前のモジュールで、TLS や HTTP/2 の実装に `golang.org/x/crypto` や `golang.org/x/net` を使っています。それらは `src/vendor` にコピーされています。`$(go env GOROOT)/src/README.vendor` には、これがほかのモジュールに影響しないと書かれています。

> The module requirements of std and cmd do not influence version selection in other modules. They are only considered when running module commands like 'go get' and 'go mod vendor' from a directory in GOROOT/src.

標準ライブラリの `go.mod` と `go.sum` は、Go 本体を開発する人のためのファイルです。利用者のビルドでは、vendor 済みのソースがそのまま使われます。そのため、利用者の `go.mod` に `golang.org/x/net` が出てくることはありません。

### どうやってインストールしているのか

標準パッケージのソースは **toolchain**（`go` コマンドと一緒に配られる Go 本体）の中にあります。手でインストールした Go では、配布物の中の `src/` が標準ライブラリです。

Go 1.21 以降は、もう 1 つ経路があります。`go.mod` の `go` 行や `toolchain` 行が手元の Go より新しいと、`go` コマンドが必要な toolchain を自分でダウンロードします（`GOTOOLCHAIN=auto`、既定）。わざと古いバージョンを指定して、その様子を見ます。

```bash
$ GOMODCACHE=$(mktemp -d) GOFLAGS=-modcacherw \
    GOTOOLCHAIN=go1.26.8 GODEBUG=toolchaintrace=1 go version
go: default toolchain set to go1.26.8 from GOTOOLCHAIN=go1.26.8
go: downloading go1.26.8 (darwin/arm64)
go: using go1.26.8 toolchain from cache located at .../golang.org/toolchain@v0.0.1-go1.26.8.darwin-arm64/bin/go
go version go1.26.8 darwin/arm64
```

ダウンロードされた toolchain の置き場所は `golang.org/toolchain@v0.0.1-go1.26.8.darwin-arm64` です。前回の記事で見た `<モジュールパス>@<バージョン>` の形と同じです。**toolchain は、`golang.org/toolchain` という特別なモジュールとして配られています。**

### どうやって検証しているのか

モジュールとして配られるなら、検証も同じ仕組みのはずです。キャッシュに残った `.ziphash` と、checksum database の値を比べます。

```bash
$ cat $GOMODCACHE/cache/download/golang.org/toolchain/@v/v0.0.1-go1.26.8.darwin-arm64.ziphash
h1:dFjGhY6f8DTsY8/rLchUTDRRfjJP4pCrGCwxvHFT2Fg=

$ curl -s https://sum.golang.org/lookup/golang.org/toolchain@v0.0.1-go1.26.8.darwin-arm64
61744041
golang.org/toolchain v0.0.1-go1.26.8.darwin-arm64 h1:dFjGhY6f8DTsY8/rLchUTDRRfjJP4pCrGCwxvHFT2Fg=
golang.org/toolchain v0.0.1-go1.26.8.darwin-arm64/go.mod h1:8wlg68NqwW7eMnI1aABk/C2pDYXj8mrMY4TyRfiLeS0=
...
```

一致しました。toolchain、つまり標準ライブラリを含む Go 本体は、`chi` と同じく checksum database で検証されています。

ただし `go.sum` には書かれません。その理由は公式ドキュメントに書かれています。

> Toolchains are downloaded like any other module, meaning that toolchain downloads can be proxied by setting `GOPROXY` and have their checksums checked by the Go checksum database. Because the specific toolchain used depends on the system's own default toolchain as well as the local operating system and architecture (GOOS and GOARCH), it is not practical to write toolchain module checksums to `go.sum`. Instead, toolchain downloads fail for lack of verification if `GOSUMDB=off`. `GOPRIVATE` and `GONOSUMDB` patterns do not apply to the toolchain downloads.
> — [Go Toolchains](https://go.dev/doc/toolchain#download)

使う toolchain は OS とアーキテクチャで変わるので、リポジトリの `go.sum` に 1 つの値を書いても意味がありません。そこで toolchain だけは、毎回 checksum database に問い合わせます。`GOSUMDB=off` にすると、検証できないのでダウンロード自体が失敗します。普通のモジュールよりも厳しい扱いです。

手でインストールする場合は、配布ページ（[go.dev/dl](https://go.dev/dl/)）にファイルごとの SHA256 が載っています。

```bash
$ curl -s 'https://go.dev/dl/?mode=json&include=all' |
    jq -c '.[] | select(.version=="go1.27.0") | .files[] | select(.os=="darwin" and .arch=="arm64") | {filename,sha256}'
{"filename":"go1.27.0.darwin-arm64.tar.gz","sha256":"90493b3bbd5e10f91d12153198bf1994fd756399b4fec93b49b0c6e2acdeeb3e"}
{"filename":"go1.27.0.darwin-arm64.pkg","sha256":"59dba31a035ddd56268904b353177053dbf5c64e893b15902733e5128ad46cbd"}
```

※ mise / Homebrew などのインストーラがこの SHA256 を照合しているかは未確認です。

### 検証は取得したときの 1 回だけ

前回の記事では、モジュールキャッシュを読むたびに `.ziphash` と `go.sum` を照合すると書きました。標準ライブラリはどうでしょうか。使い捨てのキャッシュにダウンロードした go1.26.8 の toolchain で、`fmt` パッケージを書き換えて試しました。

```bash
$ TC=$GOMODCACHE/golang.org/toolchain@v0.0.1-go1.26.8.darwin-arm64
$ ls -ld $TC/src/fmt
dr-xr-xr-x  ... src/fmt                  # モジュールと同じく読み取り専用
$ chmod u+w $TC/src/fmt $TC/src/fmt/print.go
$ printf '\nfunc init() { println("tampered fmt") }\n' >> $TC/src/fmt/print.go

$ cat main.go
package main

import "fmt"

func main() {
	fmt.Println("hello")
}

$ GOTOOLCHAIN=go1.26.8 go run .
tampered fmt
hello
```

エラーは出ず、書き換えた `fmt` がそのままビルドされました。標準ライブラリは利用者の `go.sum` に載らないので、ビルドのときに照合する相手がありません。検証されるのは toolchain をダウンロードした瞬間だけです。その後に改ざんを防いでいるのは、ファイルが読み取り専用になっていることだけです。

なお、mise で入れた手元の `GOROOT` は、ディレクトリが `drwxr-xr-x`（書き込み可）でした。手でインストールした Go は、インストール先の権限の設定によっては読み取り専用ですらありません。

:::message
この実験は、`GOMODCACHE` に指定した使い捨てディレクトリの toolchain だけを書き換えています。普段使っている `GOROOT` や `~/go/pkg/mod` で試さないでください。
:::

### 疑問 4 の回答

- 標準パッケージは import path の先頭にドットが無く、どのモジュールにも属さない。実体は toolchain の中にあるので、利用者の `go.mod` / `go.sum` には載らない
- 標準ライブラリ自身は `std` というモジュールで、外部の依存は vendor してある。その `go.mod` / `go.sum` は利用者のバージョン選択に影響しない
- 自動でダウンロードされる toolchain は `golang.org/toolchain` モジュールとして checksum database で検証される。`go.sum` には書けないので毎回 checksum database に問い合わせ、`GOSUMDB=off` では取得できない
- 検証は取得時の 1 回だけ。取得後の `GOROOT` を書き換えても、ビルドは気付かない

## 疑問 1：最後は Google を信用するしかないのか

疑問 4 で、Go 本体まで checksum database で検証していると分かりました。そうなると、`sum.golang.org` を運営する Google を信用しなければならない範囲は、さらに広がります。

### 設計の前提は「サーバーを信用しない」

設計文書（[Proposal: Secure the Public Go Module Ecosystem](https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md)）は、checksum database のサーバーを信用しない前提で書かれています。

> the log server is not trusted to store the log properly, nor is it trusted to put the right records into the log.

目標は、改ざんを **防ぐ** ことではなく、**見つけて証拠を残す** ことです。運営者が特定の利用者にだけ偽のハッシュを見せる攻撃について、設計文書はこう書いています。

> if the compromised go.sum entry is served in a forked, victim-specific log, the server must always serve that forked log to the victim, and only to the victim, or else the go command's consistency checks will fail loudly, and with enough information to cryptographically prove the compromise of the server.

偽の記録を見せた相手には、その後もずっと偽の台帳を見せ続けなければなりません。一度でも本物の台帳を見せれば、`go` コマンドが矛盾に気付きます。

### go コマンドは何を確かめているのか

`go` コマンドは checksum database から記録を受け取るたびに、2 種類の証明を確かめます（[Checksum database](https://go.dev/ref/mod#checksum-database)）。

> The go command performs "inclusion" proofs (that a specific record exists in the log) and "consistency" proofs (that the tree hasn't been tampered with) before adding new go.sum lines to the main module's go.sum file.

- inclusion proof：受け取った記録が、署名付きの台帳に本当に入っているか
- consistency proof：いま見ている台帳は、前に見た台帳へ追記しただけのものか（過去の記録を消したり書き換えたりしていないか）

consistency proof を確かめるには、「前に見た台帳」を覚えておく必要があります。その記録は `$GOPATH/pkg/sumdb` にありました（`cmd/go/internal/cfg/cfg.go` の `SumdbDir`）。

```bash
$ cat ~/go/pkg/sumdb/sum.golang.org/latest
go.sum database tree
69442998
25pvxZ0/bvFjoZt2L9oDz6caOdAOdP6jxLJ6dNWi3cw=

— sum.golang.org Az3grqcITjhAilBHAr0UD6S/W9ksRbdBM3Tak+btX9thORugQm4V0W4i9Qf1lycD+Fl+JgUcx4artG23jBULWu9/rQM=
```

「69442998 件目までの台帳で、ルートのハッシュはこれ」という、署名付きの記録です。このファイルはモジュールキャッシュの外にあるので、`go clean -modcache` では消えません。新しく受け取った台帳がこの記録と矛盾すると、`go` コマンドは次のメッセージを出して止まります（`golang.org/x/mod/sumdb/client.go`）。

```text
SECURITY ERROR
go.sum database server misbehavior detected!

old database:
	...
new database:
	...
proof of misbehavior:
	...
```

新旧 2 つの署名付き台帳と矛盾の証明を、そのまま出力します。サーバーが落ちていても、第三者が後から不正を確かめられるようにするためです。

### 見つけられること、見つけられないこと

ここまでを踏まえると、手元の `go` コマンドが見つけられる範囲は次のようになります。

| 攻撃 | 手元の go コマンドで見つかるか |
| --- | --- |
| 途中のプロキシが偽のハッシュを返す | 見つかる（署名が合わない） |
| 運営者が過去の記録を書き換える・消す | 見つかる（前に見た台帳と矛盾する） |
| 運営者が自分にだけ別の台帳を見せ、途中で本物に戻す | 見つかる（同上） |
| 運営者が自分にだけ別の台帳を見せ続ける | 見つからない（ほかの利用者と比べる手段が無い） |
| 運営者が全員に同じ偽のハッシュを見せ、台帳にも記録する | 見つからない（台帳としては正しい） |
| `latest` を持たない環境（新しい CI コンテナなど）で初めて問い合わせる | 前に見た台帳が無いので、矛盾の検出はできない |

「全員に同じ偽のハッシュ」は、台帳の中身を元のリポジトリと突き合わせる **監査者** にしか見つけられません。設計文書もそう書いています。

> If the compromised go.sum entry is stored in the actual log, an auditor can find it.

### 運営者への信用を減らす仕組みは、まだ go コマンドに入っていない

「自分にだけ別の台帳」を見つける方法も、設計文書に書かれています。利用者どうしが見ている台帳を比べ合う **gossip** です。

> the usual way to further improve detection of fork attacks is to add gossip, so that different users can check whether they are seeing different logs.

:::details gossip とは
手元の `go` コマンドが確かめられるのは、「自分が前に見た台帳」と「いま見ている台帳」が矛盾しないことだけです。運営者が被害者にだけ偽の台帳を見せ続けると、被害者から見た台帳は最初から最後まで矛盾しません。1 人で確かめている限り、偽物だと気付く材料がありません。

gossip は、利用者どうしで「自分が見ている台帳」を見せ合う仕組みです。見せ合うのは、`~/go/pkg/sumdb/sum.golang.org/latest` に入っている署名付きのツリーヘッド（何件目までの台帳で、ルートのハッシュは何か）です。

```text
運営者 ─ 本物の台帳 → 利用者 A   「69442998 件目までのハッシュは 25pv...」
      └ 偽物の台帳 → 利用者 B   「69442998 件目までのハッシュは Xq9a...」

A と B がツリーヘッドを見せ合う
→ 同じ件数なのにハッシュが違う
→ どちらにも運営者の署名が付いているので、そのまま不正の証拠になる
```

運営者は、被害者がほかの誰とも台帳を見せ合わないように、ずっと隔離し続けなければなりません。それは現実には難しいので、偽の台帳を 1 人にだけ見せる攻撃は割に合わなくなります。

設計文書は、別のプロキシ経由で最新のツリーヘッドを取り直せば gossip の代わりになると書いています。例に挙がっている `go fetch-latest-checksum-log-from-goproxy` は仮の名前で、Go 1.27.0 の `go` コマンドにそのようなコマンドはありません。
:::

ただし、この仕組みは実装されていません。Go 1.27.0 の `cmd/go` と、vendor 済みの `golang.org/x/mod/sumdb` を検索しました。`gossip`・`witness`・`cosign` のどれも、該当するコードはありませんでした。

ほかの透明性ログでは、**witness**（第三者が台帳の一貫性を確かめて連署する仕組み、[C2SP tlog-witness](https://c2sp.org/tlog-witness)）が使われ始めています。`sum.golang.org` の現在の署名を見ると、署名は運営者自身の 1 行だけでした。

```bash
$ curl -s https://sum.golang.org/latest | tail -1
— sum.golang.org Az3gr...
```

transparency.dev のブログによると、第三者の witness ネットワークが Go の checksum database を外から観測していた時期があります（[Can I get a witness network?](https://blog.transparency.dev/can-i-get-a-witness-network)）。ただし `go` コマンドはその連署を確かめていません。

:::details witness とは

gossip は、利用者が自分で見せ合う相手を探さなければなりません。witness は、その役を専門の第三者に任せる仕組みです。

witness は台帳の「立会人」です。運営者は新しいツリーヘッドを作るたびに、consistency proof を添えて witness に送ります。witness は、届いたツリーヘッドが前に見たものへ追記しただけかを確かめます。問題がなければ、自分の鍵で連署（cosignature）します。

> Witnesses verify that the checkpoint is consistent with their previously recorded state of the log (if any)
> — [C2SP tlog-witness](https://c2sp.org/tlog-witness)

```text
運営者 ─ 新しいツリーヘッド + consistency proof → witness W1, W2, W3 ...
                                               │ 前に見たものと矛盾しないか確認
                                               ↓
         運営者の署名 + W1 の連署 + W2 の連署 ... が付いたツリーヘッド
                                               ↓
                                             利用者
```

利用者が「信頼する witness の連署が付いたツリーヘッドしか受け入れない」と決めておけば、運営者が被害者にだけ偽の台帳を見せるには、witness にも偽の台帳へ連署させる必要があります。witness は前に本物の台帳を見ているので、矛盾する偽の台帳には連署しません。運営者 1 社を信じる代わりに、「運営者と witness が全員で口裏を合わせることはない」と信じればよくなります。

gossip と同じく、witness も「全員に同じ偽のハッシュを見せる」攻撃は見つけられません。台帳としては矛盾が無いからです。こちらは、台帳の中身を元のリポジトリと突き合わせる監査者の仕事です。

`sum.golang.org` の署名は運営者自身の 1 行だけで、`go` コマンドにも連署を確かめるコードはありません。Go の checksum database は、今のところ witness の仕組みを使っていません。

:::

### 2026 年 8 月に見つかった脆弱性

検証する側の `go` コマンドにも、2026 年 8 月に脆弱性が見つかっています。どちらも Go 1.27 のマイルストーンで修正されました。

- [golang/go#80744](https://go.dev/issue/80744)（CVE-2026-56865）：悪意ある `GOPROXY` が checksum database のタイル（台帳の部分木）を偽造し、検証をすり抜けられた。報告者は Filippo Valsorda（設計文書の共著者）
- [golang/go#80745](https://go.dev/issue/80745)（CVE-2026-56864）：悪意ある `GOSUMDB` が、台帳に入っていない内容を返せた

どちらの issue にも、影響を受けたかどうかを確かめる手順が書かれています。`go.sum` を作り直して、checksum database にもう一度問い合わせる方法です。

```bash
rm -r go.sum go.work.sum vendor/ && go mod tidy
```

設計が正しくても、実装に脆弱性があれば守れません。checksum database による検証を信頼できる状態に保つには、Go 本体を最新のパッチに上げておく必要があります。

### Google を外したいときの選択肢

Google を信用したくない場合は、次の設定で照合先を変えられます。

| 設定 | 効果 | 代わりに信用するもの |
| --- | --- | --- |
| `GOSUMDB="<名前>+<公開鍵> <URL>"` | 別の checksum database を使う | その database の運営者 |
| `GOPRIVATE` / `GONOSUMDB` | 指定したパスのモジュールだけ照合しない | 取得元（社内の Git など） |
| `GOSUMDB=off` | すべての未知のモジュールを照合しない | 何も確かめない（toolchain は取得できなくなる） |

どれを選んでも、信用の相手が移るだけで、信用そのものは無くなりません。設計文書は、Google の checksum database と別の組織のプロキシを組み合わせれば、攻撃者は 2 つの組織を同時に破らなければならなくなる、とも書いています。

### 疑問 1 の回答

**最初の 1 回は、Google を信用するしかありません。** ただし、Google が嘘をついても気付けない場面は限られています。台帳の書き換えや、途中で本物に戻す偽装は手元で見つかります。
見つからないのは「全員に同じ嘘をつく」場合と「1 人に嘘をつき続ける」場合です。前者は監査者が、後者は gossip や witness が見つける想定ですが、今の `go` コマンドには後者の仕組みがまだ入っていません。

もう 1 つ大事なことがあります。checksum database に問い合わせるのは、`go.sum` に行が無いときだけです。いったん `go.sum` に書いてコミットすれば、それ以降はリポジトリの `go.sum` が照合の基準になります。前回の記事の結論「信頼の単位は、リポジトリごとに git で管理している `go.sum`」は、Google を信用する場面を **最初の 1 回** に限る仕組みでもありました。

## 疑問 2：Go 1.27.0 で依存管理はどう変わったか

[Go 1.27 のリリースノート](https://go.dev/doc/go1.27)を読むと、依存管理の仕組みそのもの（MVS・`go.sum`・`GOPROXY`・`GOSUMDB`・`GOTOOLCHAIN` の既定値）は変わっていません。変わったのは、`go.mod` の書かれ方と、古い設定を残したモジュールの扱いです。

### go mod tidy が require ブロックをまとめる

`go` 行が 1.27 以上なら、`go mod tidy` が `require` ブロックを「直接依存」と「間接依存」の最大 2 つにまとめます。

> For modules specifying `go 1.27` or later in their `go.mod` file, `go mod tidy` now automatically merges duplicate require blocks.

`require` が 2 つに分かれた `go.mod` で試しました。

```go:go.mod（tidy の前）
module example.com/x

go 1.27.0

require github.com/go-chi/chi/v5 v5.3.2

require golang.org/x/text v0.37.0
```

```bash
$ go mod tidy
```

```go:go.mod（tidy の後）
module example.com/x

go 1.27.0

require (
	github.com/go-chi/chi/v5 v5.3.2
	golang.org/x/text v0.37.0
)
```

`go 1.26.0` の `go.mod` で同じことをすると、ブロックは分かれたままでした。マージのコンフリクトを手で解くたびに `require` ブロックが増えていく、という小さな悩みが解消されます。

### go mod init が書く go 行が元に戻った

リリースノートには載っていませんが、`go mod init` が書く `go` 行の既定値も変わっています。Go 1.26 では「toolchain より 1 つ古いマイナーバージョン」を書くように変わりましたが（[#74748](https://go.dev/issue/74748)）、その後「toolchain 自身のバージョン」に戻されました（[#77653](https://go.dev/issue/77653)）。

| toolchain | `go mod init` が書く go 行 |
| --- | --- |
| 1.25.3 | `go 1.25.3` |
| 1.26.0 | `go 1.25.0` |
| 1.26.2 | `go 1.26.2` |
| 1.26.8 | `go 1.26.8` |
| 1.27.0 | `go 1.27.0` |

<!-- TODO: 1.26 系のどのパッチで戻ったか（1.26.1 か 1.26.2 か）を確認する -->

1.26 のリリースノートだけを読んで「`go mod init` は 1 つ古いバージョンを書く」と覚えていると、1.27 では食い違います。

### 削除済みの GODEBUG を書いたままだと止まる

Go 1.27 では `asynctimerchan` など、いくつかの `GODEBUG` 設定が削除されました（[一覧](https://go.dev/doc/godebug)）。`go.mod` の `godebug` 行に削除済みの設定が残っている場合、最後の既定値なら受け入れられ、古い値ならビルドが止まります。

```bash
$ go build ./...
go.mod:5: removed GODEBUG "asynctimerchan" set to old value "1" (https://go.dev/doc/godebug#go-127)
```

### go test が stdversion を既定で実行する

`go test` が、vet の `stdversion` を既定で実行するようになりました。`go.mod` の `go` 行より新しい標準ライブラリの関数を使っていると、テストがビルドの段階で落ちます。

```bash
$ go test ./...
./x_test.go:...: slices.Sorted requires go1.23 or later (module is go1.21)
FAIL	example.com/x [build failed]
```

`go` 行を古いまま放置しているモジュールでは、1.27 に上げたタイミングでテストの失敗が出る可能性があります。`go` 行は「最低限必要な Go のバージョン」の宣言です。その宣言を、テストでも確かめるようになりました。

### 依存管理の安全性：checksum database クライアントの修正

疑問 1 で見た CVE-2026-56865 と CVE-2026-56864 の修正も、Go 1.27 のマイルストーンで入りました。リリースノートの機能一覧には出てきませんが、依存管理にとってはこれが一番大きな変更だと私は考えています。

### 疑問 2 の回答

- 仕組み（MVS・`go.sum`・`GOPROXY`・`GOSUMDB`・`GOTOOLCHAIN`）は変わっていない
- `go mod tidy` が `require` ブロックをまとめる（`go 1.27` 以上のとき）。`go mod init` の `go` 行は toolchain 自身のバージョンに戻った
- 古い `GODEBUG` や古い `go` 行を残したモジュールは、ビルドやテストで止まるようになった
- checksum database を検証するクライアント側の穴が 2 つ塞がれた

## 疑問の答え合わせ

| 疑問 | 答え |
| --- | --- |
| 3. dep の時代と、なぜ lock をやめたのか | dep は npm と同じ「範囲 + lock」。Go Modules は「最低バージョン」を選ぶ MVS に変えたので、結果を残す lock が要らない。lock が読まれなかったライブラリでも同じ結果を再現でき、計算も線形時間で済む |
| 4. 標準パッケージが go.mod / go.sum に無い理由 | どのモジュールにも属さず toolchain の中にある。toolchain は `golang.org/toolchain` モジュールとして checksum database で検証されるが、OS ごとに違うので `go.sum` には書かない。検証は取得時だけ |
| 1. 最後は Google を信用するしかないのか | 最初の 1 回は信用する。ただし台帳の書き換えは手元で見つかる。全員への同じ嘘は監査者頼みで、1 人への嘘の検出（gossip / witness）は go コマンドに未実装 |
| 2. Go 1.27.0 で何が変わったか | 仕組みは同じ。`go mod tidy` の require 統合、`go mod init` の go 行、古い設定への厳格化、checksum database クライアントの修正 |

## まとめ

- Go Modules に lock ファイルが無いのは、設計として決めたことだ。依存を「その日の最新」でなく「各 `go.mod` に書いた最低バージョン」から決めるので、時間が経っても結果が変わらない
- 標準ライブラリも、自動で取得する toolchain としては checksum database で検証されている。ただし `go.sum` には載らず、取得後の改ざんには気付かない
- checksum database は「運営者を信用しない」前提の設計だが、最初の 1 回の正しさは運営者を信じるしかなく、1 人だけを騙し続ける攻撃は今の `go` コマンドでは見つけられない。だからこそ、`go.sum` をコミットしてレビューすることと、Go 本体を最新のパッチに保つことが大切になる

前回の記事と合わせると、Go の依存管理は「中身の変わらない実体を共有し、信頼はリポジトリの `go.sum` に置く」設計でした。今回調べた範囲では、その `go.sum` の最初の 1 行を書くときだけ、Google の checksum database を信じています。

## 参考

- [Go Modules Reference](https://go.dev/ref/mod)
  - [Module graph pruning](https://go.dev/ref/mod#graph-pruning)
  - [Checksum database](https://go.dev/ref/mod#checksum-database)
- [Go Toolchains](https://go.dev/doc/toolchain)
- [Go 1.27 Release Notes](https://go.dev/doc/go1.27)
- [Go, Backwards Compatibility, and GODEBUG](https://go.dev/doc/godebug)
- [Proposal: Secure the Public Go Module Ecosystem](https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md)
- [Proposal: Lazy Module Loading](https://go.googlesource.com/proposal/+/master/design/36460-lazy-module-loading.md)
- [A Proposal for Package Versioning in Go](https://go.dev/blog/versioning-proposal)
- Russ Cox の vgo 連載
  - [Go += Package Versioning](https://research.swtch.com/vgo-intro)
  - [Minimal Version Selection](https://research.swtch.com/vgo-mvs)
  - [The Principles of Versioning in Go](https://research.swtch.com/vgo-principles)
  - [Version SAT](https://research.swtch.com/version-sat)
- [golang/go#24301 proposal: cmd/go: add package version support to Go toolchain](https://github.com/golang/go/issues/24301)
- [golang/dep](https://github.com/golang/dep)（[Gopkg.toml](https://github.com/golang/dep/blob/master/docs/Gopkg.toml.md) / [Gopkg.lock](https://github.com/golang/dep/blob/master/docs/Gopkg.lock.md)）
- [golang/go#80744](https://go.dev/issue/80744) / [golang/go#80745](https://go.dev/issue/80745)
- [C2SP tlog-witness](https://c2sp.org/tlog-witness)
- [Can I get a witness network? | transparency.dev](https://blog.transparency.dev/can-i-get-a-witness-network)
