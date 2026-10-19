2026 年 9 月 30 日の Go Bash vol.3 で「npm との違いを手がかりに Go の依存管理を理解したい」というタイトルで登壇しました。この記事はその内容を、登壇では時間の都合で省いた実測ログも足して書き直したものです。

@[speakerdeck](264d5427836c418685b6b783876da256)


## この記事で扱うこと

普段 Go を書いていると、`go get` を打てば依存が入り、`go.sum` も勝手に増えていきます。それで困ったこともありません。ただ、次の質問に即答できるでしょうか。

- `go.sum` は `package-lock.json` と同じものですか？
- `go get` したモジュールの実体は、どこに置かれていますか？
- そのディレクトリを、別のプロジェクトが書き換えたらどうなりますか？

私は答えられませんでした。この記事では、`go get` が内部でやっていることを 6 つのステップに分けて追いかけ、上の問いに答えていきます。

確認した環境は次のとおりです。

| 対象 | バージョン |
| --- | --- |
| Go | 1.26.8 |
| 依存に使うモジュール | `github.com/go-chi/chi/v5` v5.3.2 |
| 比較対象の Node.js / npm | Node.js v26 系 / express 5.2.1 |

最初に、登壇の冒頭で会場に投げかけた問いを紹介します。

> **npm と Go の一番の違いは何だと思いますか？**

答えは記事の最後で出します。少し考えてから読み進めてみてください。

## 同じサーバーを Go と npm で書いて比べる

「Hello」を返すだけの HTTP サーバーを、Go（chi）と Node.js（express）で書きました。やっていることは同じです。

```go:main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Hello, Go server!")
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
```

```js:index.js
const express = require("express");

const app = express();
app.get("/", (req, res) => {
  res.send("Hello, Node server!\n");
});
app.listen(8080);
```

依存の実体がどこにあるかを調べます。npm の場合は、プロジェクト直下の `node_modules` にあります。

```bash
$ npm ls express --parseable
/path/to/ex-npm/node_modules/express

$ ls -1 node_modules | wc -l
      65
$ find node_modules -type f | wc -l
     601
```

express を 1 つ指定しただけで、`node_modules` には 65 ディレクトリ・601 ファイル（3.8MB）が展開されていました。

一方、Go のプロジェクトを `ls` すると `main.go`・`go.mod`・`go.sum` の 3 つしかありません。chi を import してビルドできているのに、chi のソースはプロジェクトの中に無いのです。どこにあるのかを `go` コマンドに聞きます。

```bash
# ホームディレクトリは ~ に置き換えています
$ go env GOMODCACHE
~/go/pkg/mod

$ go list -m -f '{{.Dir}}' github.com/go-chi/chi/v5
~/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2
```

実体は、ホームディレクトリの下の **モジュールキャッシュ**（`GOMODCACHE`、既定は `~/go/pkg/mod`）にありました。このキャッシュはマシンに 1 つで、どのプロジェクトからも同じ場所が読まれます。

```text
Go
  プロジェクトA ─┐
                ├─→ ~/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2
  プロジェクトB ─┘     マシンに 1 つ

npm
  プロジェクトA/node_modules/express   ← プロジェクトごとに別のコピー
  プロジェクトB/node_modules/express
```

## 共有して、困らないのか？

npm に慣れていると、ここで不安になります。1 つの実体を全プロジェクトで共有して、本当に大丈夫なのでしょうか。私が感じた疑問は 4 つです。

1. なぜ同じディレクトリに置いているのか？
2. 別のプロジェクトが別のバージョン（たとえば v5.2.0）を使いたくなったら？
3. 別のプロジェクトがキャッシュの中身を書き換えたら？
4. モジュールの作者が v5.3.2 のタグを別のコミットに付け替えたら？

この 4 つの答えは、どれも `go get` の処理の中にあります。順番に追っていく前に、用語を 2 つだけ揃えておきます。

### Go の「モジュール」は npm の「パッケージ」にあたる

同じ「パッケージ」という言葉でも、Go と npm では指すものが逆です。

| | Go | npm |
| --- | --- | --- |
| バージョンを付けて配布・取得される単位 | module（`go.mod` があるもの） | package（`package.json` があるもの） |
| コードを読み込む単位 | package（同じディレクトリにあり、まとめてコンパイルされる `.go` ファイル群） | module（`require` / `import` で読み込めるファイルやディレクトリ） |

この記事で npm の package と比べるのは、Go の **モジュール** のほうです。

### package.json は範囲、go.mod は最低バージョン

バージョンの書き方も違います。

```json:package.json
{
  "dependencies": {
    "express": "^5.2.1"
  }
}
```

```go:go.mod
module github.com/shirokuma/conference-report/GoBash/ex-Go

go 1.26.8

require github.com/go-chi/chi/v5 v5.3.2
```

`^5.2.1` は「5.2.1 以上 6.0.0 未満」という **範囲** です。範囲のうちどのバージョンが入るかはインストールした時点の最新版で変わるので、npm は選んだ結果を `package-lock.json` に記録します。範囲の宣言と結果の記録、2 つのファイルで 1 組です。

`go.mod` の `require` には範囲を書きません。公式リファレンスでは、`require` は依存モジュールの **最低バージョン** を宣言するものだと定義されています。

> A require directive declares a minimum required version of a given module dependency.

ここで次の疑問が湧きます。範囲を書かないなら、選んだ結果を記録する必要もなさそうです。**では、`go.sum` は何を記録しているのでしょうか。** 冒頭の「`go.sum` は lock ファイルか」という問いも、ここにつながります。

## go get の出力に残る 3 つの謎

モジュールキャッシュが空の状態で `go get` を 1 回打つと、出力はこれだけです。

```bash
$ go get github.com/go-chi/chi/v5
go: downloading github.com/go-chi/chi/v5 v5.3.2
go: downloading github.com/go-chi/chi v1.5.5      # ① 頼んだのは 1 つなのに 2 つ取得している
go: added github.com/go-chi/chi/v5 v5.3.2         # ② v5.3.2 はどう決まった？

$ cat go.sum
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
# ③ なぜ 2 行？ h1: とは何？
```

短い出力ですが、1 行ずつ読むと分からないところが 3 つあります。特に ① は、`/v5` を頼んだのに v1 系の `chi v1.5.5` まで落としてきているので不思議です。

`go get` は内部で、だいたい次の 6 ステップを踏んでいます（[Go Modules Reference](https://go.dev/ref/mod)）。

1. 必要なモジュールとバージョンを決める
2. モジュールプロキシ（または Git などの VCS）からモジュールを取得する
3. モジュールキャッシュに保存する
4. 暗号学的ハッシュを計算する
5. `go.sum` または checksum database（`sum.golang.org`）と照合する
6. `go.sum` に記録する

このステップを 1 つずつ追うと、①〜③ の謎と、さきほどの 4 つの疑問が順に解けていきます。

:::message
`go: downloading` の行はキャッシュが空のときしか出ません。手元のキャッシュを汚さずに再現するには、使い捨てのディレクトリを `GOMODCACHE` に指定します。キャッシュは読み取り専用で作られるため、後で消せるように `-modcacherw` も付けておきます。

```bash
GOMODCACHE=$(mktemp -d) GOFLAGS=-modcacherw go get -x github.com/go-chi/chi/v5
```
:::

## STEP1 必要なモジュールとバージョンを決める

### import path のどこまでがモジュールなのか

`import "github.com/go-chi/chi/v5"` と書いても、`go` コマンドには「どこまでがモジュールパスで、どこからがモジュール内のディレクトリか」が分かりません。

```text
github.com/go-chi/chi/v5   ← この名前のモジュール？
github.com/go-chi/chi      ← このモジュールの中の v5/ ディレクトリ？
```

そこで `go` コマンドは 2 段階で探します（[Resolving a package to a module](https://go.dev/ref/mod#resolve-pkg-mod)）。

まず手元の **ビルドリスト**（ビルドに使うモジュールとバージョンの一覧。`go list -m all` で見られます）から、モジュールパスが import path の接頭辞になっているものを探します。`go.mod` に `require` があればここで決まり、ネットワークは使いません。

```bash
$ GOPROXY=off go list -f '{{.ImportPath}} → {{.Module.Path}} {{.Module.Version}}' \
    github.com/go-chi/chi/v5 github.com/go-chi/chi/v5/middleware
github.com/go-chi/chi/v5 → github.com/go-chi/chi/v5 v5.3.2
github.com/go-chi/chi/v5/middleware → github.com/go-chi/chi/v5 v5.3.2
```

初回の `go get` ではビルドリストに chi が無いので、次の段階に進みます。import path の末尾から 1 要素ずつ削ったパスをすべてモジュールパスの候補にして、モジュールプロキシに「最新バージョンは？」と並列に問い合わせるのです。`-x` を付けると、その様子が見えます。

```bash
$ go get -x github.com/go-chi/chi/v5 2>&1 | grep -E '/@v/list: |go: downloading'
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/list: 200 OK
# get https://proxy.golang.org/github.com/go-chi/chi/@v/list: 200 OK
# get https://proxy.golang.org/github.com/go-chi/@v/list: 404 Not Found
# get https://proxy.golang.org/github.com/@v/list: 404 Not Found
go: downloading github.com/go-chi/chi/v5 v5.3.2
go: downloading github.com/go-chi/chi v1.5.5
```

200 が返ったのは `chi/v5` と `chi` の 2 つです。公式リファレンスによると、`go` コマンドは応答があったモジュールの最新バージョンを取得し、要求されたパッケージを含むかどうかを確かめます。含むモジュールが複数あれば、パスが最も長いものを使います。

**① の謎はこれで解けます。** `chi v1.5.5` は、「`chi` モジュールの中に `v5/` ディレクトリがあるのでは？」という候補を確かめるために取ってきたものでした。

### MVS でバージョンを決める

モジュールが決まったら、次はバージョンです。Go は **MVS（Minimal Version Selection）** という方法で決めます（[Minimal version selection (MVS)](https://go.dev/ref/mod#minimal-version-selection)）。

MVS はメインモジュールから各モジュールの `go.mod` の `require` をたどり、モジュールごとに「要求された中で最も高いバージョン」を記録します。たどり終えた時点の記録がビルドリストです。

```text
メインモジュール ─require→ 依存 A ── chi v5.0.0 以上が必要
                └require→ 依存 B ── chi v5.2.0 以上が必要
                                      ↓ MVS
                          ビルドリスト: chi v5.2.0
```

「最も高い」なのに「Minimal」と呼ぶのは、これが **すべての要求を満たす最小のバージョン** だからです。上の例で v5.3.2 が公開されていても、誰も要求していなければ選びません。

この決め方の材料は、各モジュールの `go.mod` だけです。新しいバージョンが公開されても結果は変わりません。公式リファレンスはこう書いています。

> MVS is deterministic, and the build list doesn't change when new versions of dependencies are released.
> (中略)
> Unlike other dependency management systems, the build list is not saved in a "lock" file.

npm は範囲の中から「いま一番新しいもの」を選ぶため、実行する時期によって結果が変わります。lock ファイルが必要なのはこのためです。Go は `go.mod` に書いた最低バージョンとほかのモジュールの `go.mod` から、毎回同じ結果を計算できます。だから結果を別ファイルに残す必要がありません。

## STEP2 モジュールプロキシから取得する

取得先は `GOPROXY` で決まります。既定値はこうです。

```bash
$ go env GOPROXY
https://proxy.golang.org,direct
```

`proxy.golang.org` が **モジュールプロキシ** です。Go チームが運営する、公開モジュールのキャッシュサーバーです。`,` はフォールバックの区切りで、プロキシが 404 / 410 を返したら次の `direct`（プロキシを使わず Git などから直接取る）に進みます。社内リポジトリのように `GOPRIVATE` に一致するモジュールも `direct` で取得されます。

モジュールプロキシは、npm の registry とは配布物の集め方が違います。npm では作者が `npm publish` で registry にアップロードします。Go のモジュールプロキシは、初めて要求されたバージョンを自分で元リポジトリから取ってきて zip にします。いったんキャッシュされると、作者が元リポジトリでソースを消しても取得できます（[proxy.golang.org](https://proxy.golang.org/) の FAQ）。

空のキャッシュで `go get -x` した出力から、`chi/v5` の取得に関わる行だけを抜き出すとこうなります（URL 以外の情報は省略しています）。

```bash
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/list         ← バージョン一覧
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/v5.3.2.info  ← メタデータ
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/v5.3.2.mod   ← go.mod 単体
# get https://sum.golang.org/lookup/github.com/go-chi/chi/v5@v5.3.2     ← 照合用のハッシュ（STEP5）
go: downloading github.com/go-chi/chi/v5 v5.3.2
# get https://proxy.golang.org/github.com/go-chi/chi/v5/@v/v5.3.2.zip   ← モジュール zip ファイル
go: added github.com/go-chi/chi/v5 v5.3.2
```

バージョンを指定しない `go get` は、`@v/list` で取った一覧からプレリリースを除き、SemVer で最も大きいバージョンを選びます。

```bash
$ curl -s https://proxy.golang.org/github.com/go-chi/chi/v5/@v/list | sort -V | tail -3
v5.3.0
v5.3.1
v5.3.2
```

**② の謎はここで解けます。** v5.3.2 は「その時点の最新リリース」として選ばれ、`go.mod` の `require` に書かれました。`added` は「`go.mod` に書いた」という意味です。いったん書かれたあとは、STEP1 で見た MVS がこの最低バージョンを材料にビルドリストを決めます。

`.mod` と `.zip` が別々に取得されている点も覚えておいてください。MVS は依存モジュールの `go.mod` しか読みません。そのため `go` コマンドは、ソース全体を落とさずに `go.mod` だけを先に取得できます。これが ③ の伏線になります。

## STEP3 モジュールキャッシュに保存する

取ってきたものは `GOMODCACHE` の下に、役割の違う 2 つの層に分けて置かれます。

```text
~/go/pkg/mod/
├── cache/download/github.com/go-chi/chi/v5/@v/   ← 層1: 取得したものをそのまま置く
│   ├── list
│   ├── v5.3.2.info      ← バージョンと取得元（git のコミット）
│   ├── v5.3.2.mod       ← go.mod 単体
│   ├── v5.3.2.zip       ← モジュール zip ファイル
│   ├── v5.3.2.ziphash   ← zip のハッシュ（STEP5 で使う）
│   └── v5.3.2.lock
│
└── github.com/go-chi/chi/v5@v5.3.2/              ← 層2: zip を展開したソース（ビルドが読む）
    ├── chi.go
    ├── mux.go
    └── ...
```

`go build` や `go test` が読むのは層2です。層2のディレクトリ名は `<モジュールパス>@<バージョン>` の形です。展開した時点で、中のディレクトリとファイルは **読み取り専用** になります。

```bash
$ ls -ld ~/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2
dr-xr-xr-x@ 24 ... github.com/go-chi/chi/v5@v5.3.2
```

`go get` の出力にあった `downloading` の 2 行は、どちらもこのキャッシュに置かれています。`chi v1.5.5` も `github.com/go-chi/chi@v1.5.5/` として同じ場所に残っています。

`.info` を開くと、プロキシがどこから作った zip なのかが分かります。

```json
{"Version":"v5.3.2","Time":"2026-08-20T09:37:52Z",
 "Origin":{"VCS":"git","URL":"https://github.com/go-chi/chi",
 "Hash":"38939062c5df4d3e8814aad1a488983112627ced","Ref":"refs/tags/v5.3.2"}}
```

「git の `refs/tags/v5.3.2`、コミット `3893906` から作った」と記録されています。疑問 4（タグの付け替え）を考えるときに、この記録が効いてきます。

## STEP4 暗号学的ハッシュを計算する

`go.sum` に並ぶ `h1:` は、SHA-256 を使うハッシュ方式「Hash1」の識別子です。対象はモジュール zip ファイルですが、**zip のバイト列をそのままハッシュしているわけではありません**。手順は次の 4 段です。

1. zip 内の全ファイルについて SHA-256 を計算する
2. `"<SHA-256 の hex>  <ファイル名>\n"` という行を作る（スペースは 2 個）
3. ファイル名でソートして全行を連結する（これが「リスト」）
4. リスト全体の SHA-256 を取り、その 32 バイトを base64 にして `h1:` を付ける

本当にそうなっているのか、シェルで再現してみます。

```bash
$ cd $(mktemp -d)
$ unzip -q ~/go/pkg/mod/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.zip
$ find github.com -type f | LC_ALL=C sort |  # 3. ファイル名でソート
    xargs shasum -a 256 |                    # 1〜3. 「<hex>  <ファイル名>」の行を順に出力
    shasum -a 256 | cut -d' ' -f1 |          # 4. リスト全体の SHA-256
    xxd -r -p | base64                       # 4. hex を 32 バイトに戻して base64
5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
```

先頭に `h1:` を付けると、`go.sum` の 1 行目と一致しました。`shasum` の出力形式が「hex + スペース 2 個 + ファイル名」なので、ソートしたファイル名の順に並べるだけでリストになります。`LC_ALL=C` を付けているのは、Go と同じバイト順でソートするためです。

ファイル名と中身だけから計算するので、zip の圧縮レベルやタイムスタンプ、エントリの順序には左右されません。誰が zip を組み立てても、中身が同じなら同じ値になります。

2 行目の `/go.mod` 行も同じ Hash1 です。対象は `v5.3.2.mod` の 1 ファイルだけで、リストに書くファイル名は `go.mod` です。

```bash
$ D=~/go/pkg/mod/cache/download/github.com/go-chi/chi/v5/@v
$ printf '%s  go.mod\n' "$(shasum -a 256 $D/v5.3.2.mod | cut -d' ' -f1)" |
    shasum -a 256 | cut -d' ' -f1 | xxd -r -p | base64
R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
```

Go のコードから計算するなら、`golang.org/x/mod/sumdb/dirhash` の `dirhash.HashZip(path, dirhash.Hash1)` で同じ値が得られます。

## STEP5 go.sum または checksum database と照合する

計算したハッシュを何と比べるかは、`go.sum` に行があるかどうかで変わります。

### go.sum に行がある場合

`go.sum` の値が正解です。手元で計算したハッシュと一致すれば、checksum database には問い合わせずに検証を終えます。一致しなければ、取得したものを捨ててエラーで止まります。

```text
SECURITY ERROR
This download does NOT match an earlier download recorded in go.sum.
```

### go.sum に行が無い場合

checksum database（既定は `sum.golang.org`）に正解のハッシュを問い合わせます。STEP2 の `-x` 出力にあった `lookup` がこれです。

```bash
$ curl https://sum.golang.org/lookup/github.com/go-chi/chi/v5@v5.3.2
60295460                                                  # ← 台帳の何番目の記録か
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=

go.sum database tree                                      # ← ここから下は署名付きのツリーヘッド
63096908
8lHrQPRYe4KT4WQ9/Qk8/taKnD6IRHbFgfN7vQE2Y54=

— sum.golang.org Az3gri/oVBN1jAhNZ/iS+W31psFAWRT20/GyrhpU6taFwulQGetBkH5S7wpGXdS935jiIGiUUAXBVHzT/P6fYhqAmAk=
```

STEP4 で計算した値と同じハッシュが返ってきています。checksum database は追記専用の台帳（Merkle tree を使った透明性ログ）で、一度記録したハッシュを後から書き換えられません。

最終行は `sum.golang.org` の電子署名（Ed25519）です。この署名を確かめる公開鍵は、`go` コマンドのソースに定数として書かれています。

```go:src/cmd/go/internal/modfetch/key.go
var knownGOSUMDB = map[string]string{
	"sum.golang.org": "sum.golang.org+033de0ae+Ac4zctda0e5eza+HJyk9SxEdh+s3Ux18htTTAD8OuAn8",
}
```

定数なので、ビルドすると `go` の実行ファイルにそのまま入ります。

```bash
$ strings $(which go) | grep -o 'sum.golang.org+033de0ae+.\{44\}'
sum.golang.org+033de0ae+Ac4zctda0e5eza+HJyk9SxEdh+s3Ux18htTTAD8OuAn8
```

途中のプロキシは秘密鍵を持っていないので、応答を書き換えると署名が合わなくなります。検証用の鍵は手元の `go` コマンドの中にあるので、鍵ごと差し替えることもできません。

### 2 回目以降

キャッシュ済みのモジュールを使うたびに、`go` コマンドはメインモジュールの `go.sum` と照合しています。毎回 zip を再計算すると重いので、STEP3 で保存した `.ziphash`（計算済みのハッシュ）と `go.sum` を比べます。

> To avoid the need to trust other modules, the `go` command verifies hashes using the main module's `go.sum` whenever it accesses a file in the module cache. Zip file hashes are expensive to compute, so the `go` command checks pre-computed hashes stored alongside zip files instead of re-hashing the files.
> — [Authenticating modules](https://go.dev/ref/mod#authenticating)

"To avoid the need to trust other modules" という一節を覚えておいてください。最後に「一番の違い」を考えるときに使います。

## STEP6 go.sum に記録する

検証を通ったハッシュが `go.sum` に追記されます。行数にはルールがあります（[go.sum files](https://go.dev/ref/mod#go-sum-files)）。

- ビルドにソースが必要なモジュールは、zip の `h1:` と `/go.mod` の `h1:` の **2 行**
- モジュールグラフをたどるためだけに使うモジュールは、`/go.mod` の **1 行**

**③ の謎の答えです。** 2 行あるのは、STEP2 で見たとおり zip と `go.mod` を別々に取得し、別々に検証しているからです。MVS は `go.mod` しか読まないので、ソースが要らないモジュールについては `go.mod` の行だけで足ります。

ここで改めて `go.sum` を見てください。

```text
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
```

STEP1 で取得した `github.com/go-chi/chi` v1.5.5 は、zip の取得・展開も checksum database での検証も済んでいます。それなのに `go.sum` には 1 行もありません。import path を解決するための「探り」で取っただけで、最終的なモジュールグラフには入っていないからです。

キャッシュには「見に行ったもの全部」が残り、`go.sum` には「ビルドに関わるものだけ」が載ります。この非対称が、最後の答えにつながります。

## go.sum は lock ファイルなのか

6 ステップを見終えたので、途中で保留した問いに戻ります。`go.sum` を消したらどうなるかを試しました。

```bash
$ mv go.sum /tmp/go.sum.bak
$ go build ./...
main.go:8:2: missing go.sum entry for module providing package github.com/go-chi/chi/v5 (imported by github.com/shirokuma/conference-report/GoBash/ex-Go); to add:
        go get github.com/shirokuma/conference-report/GoBash/ex-Go
```

ビルドは止まりました。Go 1.16 以降は `-mod=readonly` が既定なので、`go build` が `go.sum` を勝手に作り直すことはありません。`go mod tidy` で作り直して、元のファイルと比べます。

```bash
$ go mod tidy
$ diff /tmp/go.sum.bak go.sum
$ # 差分なし
```

1 バイトも変わりませんでした。`go.sum` が無いとビルドは止まりますが、作り直しても中身は同じです。つまり **`go.sum` はバージョンの決定に関わっていません**。バージョンを決めているのは `go.mod` と MVS です。

| ファイル | 役割 |
| --- | --- |
| `go.mod` | どのモジュールのどのバージョン以上が要るかの **宣言** |
| `go.sum` | 取ってきた中身が本物であることを確かめる **証明** の台帳 |

`go.sum` は lock ファイルではなく、検証のための台帳でした。

## 4 つの疑問の答え合わせ

最初に挙げた 4 つの疑問に、ここまでのステップで答えます。

| 疑問 | 答え | 根拠 |
| --- | --- | --- |
| 1. なぜ同じディレクトリに置いているのか？ | 同じ `<モジュール>@<バージョン>` なら、中身の `h1:` ハッシュは 1 つに決まる。1 つ置けばどのプロジェクトも同じものを使える | STEP4・5 |
| 2. 別のプロジェクトが別のバージョンを使いたくなったら？ | ディレクトリ名にバージョンが入るので、別ディレクトリとして共存する。どのバージョンを使うかは各プロジェクトの `go.mod` が決める | STEP1・3 |
| 3. 別のプロジェクトが中身を書き換えたら？ | 展開済みのソースは読み取り専用なので、そもそも書き込めない | STEP3 |
| 4. 作者がタグを別のコミットに付け替えたら？ | プロキシにキャッシュ済みなら、`.info` に記録された元のコミットから作った zip がそのまま配られる。別経路で付け替え後の中身を取っても、`go.sum` や checksum database に残る最初のハッシュと合わずエラーになる | STEP2・3・5 |

疑問 3 は実際に試せます。

```bash
$ M=~/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2
$ echo "// hack" >> $M/mux.go
zsh: permission denied: .../chi/v5@v5.3.2/mux.go
```

`chmod` で権限を外して無理やり書き換えた場合は、`go mod verify` で見つけられます。`go mod verify` はキャッシュの zip と展開済みディレクトリのハッシュを計算し直し、ダウンロード時に記録した `.ziphash` と比べるコマンドです。合わなければ `dir has been modified` などのエラーを出します。

:::message
ビルドのたびに展開済みディレクトリの全ファイルを計算し直しているわけではありません。STEP5 で見たとおり、普段の照合は計算済みの `.ziphash` と `go.sum` の比較です。手でキャッシュを触ってしまったときは `go mod verify` を実行してください。`go mod verify` の照合相手は、Go 1.27.0 の `cmd/go/internal/modcmd/verify.go` で確認しました。
:::

おまけの疑問にも 1 つ答えておきます。途中で checksum database の応答まで偽のハッシュに書き換えられた場合は、STEP5 で見た署名の検証で止まります。

## 最初の問いに戻る：npm と Go の一番の違い

冒頭の問いに戻ります。**npm と Go の一番の違いは何でしょうか。**

表にすると、違いはいくつも挙がります。

| 観点 | Go | npm |
| --- | --- | --- |
| 依存の置き場所 | `GOMODCACHE`（マシンに 1 つ） | プロジェクトごとの `node_modules` |
| 同一バージョンの重複 | 無い | プロジェクトの数だけコピーされる |
| バージョンの宣言 | `go.mod` に最低バージョンを 1 つ | `package.json` に範囲、lock に結果 |
| 依存の書き込み | 読み取り専用 | 書き換えられる |
| git 管理するもの | `go.mod` と `go.sum` | `package.json` と `package-lock.json` |

登壇の準備を始めたとき、私は一番の違いを「依存の置き場所」だと考えていました。`node_modules` が無いのは、Go を触ると最初に気付く違いだからです。

ただ、6 ステップを追った今は少し違う答えを出せます。npm も `~/.npm/_cacache` に共有キャッシュを持っているので、実体を共有すること自体は Go だけの特徴ではありません。違いは、共有しても安全だと言える仕組みが仕様で決まっている点にあります。

- バージョンの中身は変わらない（同じバージョンなら `h1:` は 1 つに決まり、checksum database がそれを記録し続ける）
- 使うバージョンは各プロジェクトの `go.mod` だけで決まる（MVS）
- キャッシュを読むたびに、各プロジェクトの `go.sum` と照合する

STEP5 で引用した一節 "To avoid the need to trust other modules" が言っているのはこのことです。キャッシュには、ほかのプロジェクトが取ってきたものも探りで取っただけのものも入っています。それでも各プロジェクトは自分の `go.sum` で照合するので、ほかのプロジェクトを信用する必要がありません。

**Go が共有しているのは依存の「実体」だけで、「信頼」は共有していません。信頼の単位は、リポジトリごとに git で管理している `go.sum` です。**

これが、私の考える「一番の違い」です。

裏を返すと、`go.sum` 自体を書き換えられたら、この仕組みでは防げません。`go.sum` に行があれば checksum database には問い合わせないので、偽のハッシュを書いた `go.sum` がリポジトリに入れば、`go` コマンドはそれを正解として扱います。`go.sum` の差分は、自動生成のファイルだからと読み飛ばさず、コードと同じようにレビューしてください。

## 調べきれなかったこと

登壇までに調べきれず、宿題として残っているものです。

- 最後は Google（`sum.golang.org` の運営者）を信用するしかないのか
- Go 1.27.0 で依存管理がどう変わったのか
- dep（`Gopkg.lock`）の時代の依存管理と、なぜ lock ファイルをやめたのか
- 標準パッケージが `go.mod` と `go.sum` に記録されない理由（どうやってインストールし、どうやって検証しているのか）

1 つ目については手がかりだけ見つけています。透明性ログには、運営者が特定の人にだけ違うハッシュを返すと、台帳の枝分かれとして検出できる仕組みがあります。ただ、それで運営者への信用をどこまで減らせているのかは、まだ整理できていません。

分かったら続きの記事を書く予定です。詳しい方がいたら、ぜひコメントで教えてください。

## まとめ

- `go get` は「決める → 取得する → キャッシュに置く → ハッシュを計算する → 照合する → 記録する」の 6 ステップで動いている
- `go.mod` は最低バージョンの宣言で、MVS により同じ結果が毎回計算できるので lock ファイルが要らない
- `go.sum` は lock ファイルではなく、取ってきた中身が本物かを確かめる台帳。依存の実体は全プロジェクトで共有し、信頼は各リポジトリの `go.sum` に置いている

次のコマンドを一度打ってみると、ここで追った 6 ステップがそのまま流れていくのが見えます。

```bash
GOMODCACHE=$(mktemp -d) GOFLAGS=-modcacherw go get -x <モジュール>
```

普段使っているモジュールで試してみてください。

## 参考

- [Go Modules Reference](https://go.dev/ref/mod)
  - [Resolving a package to a module](https://go.dev/ref/mod#resolve-pkg-mod)
  - [Minimal version selection (MVS)](https://go.dev/ref/mod#minimal-version-selection)
  - [Module cache](https://go.dev/ref/mod#module-cache)
  - [Authenticating modules](https://go.dev/ref/mod#authenticating)
  - [go.sum files](https://go.dev/ref/mod#go-sum-files)
- [Managing dependencies](https://go.dev/doc/modules/managing-dependencies)
- [proxy.golang.org](https://proxy.golang.org/)
- [About packages and modules | npm Docs](https://docs.npmjs.com/about-packages-and-modules)
