# 構成

## 登壇情報

```text
イベント    Go Bash vol.3（共催4社トーク）
日時        2026-09-30(水) 19:35-20:25 が4社トーク枠
持ち時間    50分 ÷ 4人 = 実質 10〜11分（転換込み）
会場        Eviry 六本木 38F / オンサイト / 40枠
タイトル    npm との違いを手がかりに Go の依存管理を理解したい
この後      20:40-21:10 ディスカッション(30分) → 21:15 懇親会
```

他3社は「作った話 / やった話」（AI Agent・災害データ・探検部レポート）。
内部構造を掘るのはこの枠だけなので、**深さで勝負する**。

## 方針

- **柱は1本**。疑問5 を「なぜ他のプロジェクトを信用しなくて済むのか」として最初に問い、最後に回収する
  - 「実体の置き場が違う」は**事実**であって問いではない。置き場が生む**帰結**まで降ろして初めて問いになる
- MVS は**定義1文と公式の引用だけ**入れる（スライド4 ステップ2）。グラフの例や `go mod graph` は出さない
- dep 史 / go 1.27 は**今回入れない**（未検証、かつ尺に入らない）
- 10分で完結させない。**ディスカッション30分に投げる球を残して終わる**
- 1枚1メッセージ。11枚前後、1枚あたり約50秒

## タイムテーブル（10分）

<!-- 下記は案 -->
| 時間 | 内容 | ネタ元 |
| --- | --- | --- |
| 0:00 | 自己紹介 + OPTiM（**30秒・1枚**。ここで2分使うと本題が8分になる） | - |
| 0:30 | 【問い】`ls` を2つ並べる / `go env GOMODCACHE` で即答 / **2プロジェクト共有の絵** / 本当の問いを提示 | 疑問5 |
| 2:00 | 【前提】package.json は幅、go.mod はピンポイント（go.mod は「宣言」） | 疑問1（圧縮） |
| 3:00 | 【go get の中身】出力4行を伏線に、探す → 決める（MVS）→ 置く → 書く → 照合する | 疑問4 / 7 / 9 |
| 6:30 | 【答え合わせ】p12 の4つの問いに回答 + read-only の実演 | 疑問5 |
| 8:15 | 【まとめ】npm / pnpm / Go の対比 → 「実体は共有、信頼は共有しない」→ 議論への問い | 疑問5 |

## スライド構成

### 1. 自己紹介 / OPTiM（30秒）

1枚。所属と一言だけ。会社紹介パートは作らない。

### 2. 問い 「npm と Golang の一番の違いは何だと思いますか？」

90秒の配分。

<!-- 下記は案 -->

| 秒 | 内容 |
| --- | --- |
| 0:30 | `ls` を2つ並べる。**3秒黙る** |
| 0:50 | 「chi を import できているのに実体がない。どこにある？」 |
| 1:05 | `go env GOMODCACHE` で即答 |
| 1:15 | **2プロジェクト共有の絵**（＋npm は別コピー） |
| 1:40 | 本当の問いを出す |

#### 「Hello World !!」をプリントするサーバーで比較

**サーバーのコードを提示する**

どちらも「HTTP サーバー + ルーター」。やっていることは同じ。

- Golang
  - ./../ex-Go/main.go
- npm
  - ./../ex-npm/index.js

#### パッケージが入っているpathを調べる

npmのimageファイル作成
<!-- https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=vscode&wt=none&l=application%2Fx-sh&width=680&ds=true&dsyoff=0px&dsblur=0px&wc=true&wa=true&pv=0px&ph=0px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=%2523%2520%25E3%2582%25AD%25E3%2583%25A3%25E3%2583%2583%25E3%2582%25B7%25E3%2583%25A5%25E3%2581%25A8%25E3%2581%2597%25E3%2581%25A6%25E4%25BF%259D%25E5%25AD%2598%25E3%2581%2595%25E3%2582%258C%25E3%2582%258B%25E3%2583%2587%25E3%2582%25A3%25E3%2583%25AC%25E3%2582%25AF%25E3%2583%2588%25E3%2583%25AA%25E3%2581%25AE%25E3%2583%2591%25E3%2582%25B9%250Ago%2520env%2520GOMODCACHE%250A%252FUsers%252Fopm008296%252Fgo%252Fpkg%252Fmod%250A%250A%2523%2520%25E3%2583%2580%25E3%2582%25A6%25E3%2583%25B3%25E3%2583%25AD%25E3%2583%25BC%25E3%2583%2589%25E3%2581%2597%25E3%2581%259F%25E3%2583%2590%25E3%2583%25BC%25E3%2582%25B8%25E3%2583%25A7%25E3%2583%25B3%25E3%2581%25AE%25E3%2583%2591%25E3%2583%2583%25E3%2582%25B1%25E3%2583%25BC%25E3%2582%25B8%25E3%2581%25AEpath%25E3%2582%2592%25E8%25A1%25A8%25E7%25A4%25BA%250Ago%2520list%2520-m%2520-f%2520%27%257B%257B.Dir%257D%257D%27%2520github.com%252Fgo-chi%252Fchi%252Fv5%250A%252FUsers%252Fopm008296%252Fgo%252Fpkg%252Fmod%252Fgithub.com%252Fgo-chi%252Fchi%252Fv5%2540v5.3.2 -->

Golangのimageファイル作成
<!-- https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=vscode&wt=none&l=application%2Fx-sh&width=680&ds=true&dsyoff=0px&dsblur=0px&wc=true&wa=true&pv=0px&ph=0px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=%2523%2520%25E3%2582%25AD%25E3%2583%25A3%25E3%2583%2583%25E3%2582%25B7%25E3%2583%25A5%25E3%2581%25A8%25E3%2581%2597%25E3%2581%25A6%25E4%25BF%259D%25E5%25AD%2598%25E3%2581%2595%25E3%2582%258C%25E3%2582%258B%25E3%2583%2587%25E3%2582%25A3%25E3%2583%25AC%25E3%2582%25AF%25E3%2583%2588%25E3%2583%25AA%25E3%2581%25AE%25E3%2583%2591%25E3%2582%25B9%250Ago%2520env%2520GOMODCACHE%250A%252FUsers%252Fopm008296%252Fgo%252Fpkg%252Fmod%250A%250A%2523%2520chi%25E3%2582%2592%25E5%2585%25A5%25E3%2582%258C%25E3%2581%25A6%25E3%2581%2584%25E3%2582%258Bpath%25E3%2582%2592%25E8%25A1%25A8%25E7%25A4%25BA%250Ago%2520list%2520-m%2520-f%2520%27%257B%257B.Dir%257D%257D%27%2520github.com%252Fgo-chi%252Fchi%252Fv5%250A%252FUsers%252Fopm008296%252Fgo%252Fpkg%252Fmod%252Fgithub.com%252Fgo-chi%252Fchi%252Fv5%2540v5.3.2 -->

実測済み（2026-09-23 / node v26.2.0 / express 5.2.1）。採取元は `ex-npm/README.md`。

```bash
go env GOMODCACHE
/Users/shiraki_takuma/go/pkg/mod

cd GoBash/ex-Go
go list -m -f '{{.Dir}}' github.com/go-chi/chi/v5
/Users/opm008296/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2

cd GoBash/ex-npm
npm ls express --parseable
/Users/opm008296/git/github.com-takumashiraki/takumashiraki/conference-report/GoBash/ex-npm/node_modules/express
```

「chi を import できているのに、リポジトリに実体がない。どこにある？」

```bash
ls GoBash/ex-npm/node_modules
accepts                 encodeurl               gopd                    ms                      send
body-parser             es-define-property      has-symbols             negotiator              serve-static
bytes                   es-errors               hasown                  object-inspect          setprototypeof
call-bind-apply-helpers es-object-atoms         http-errors             on-finished             side-channel
call-bound              escape-html             iconv-lite              once                    side-channel-list
content-disposition     etag                    inherits                parseurl                side-channel-map
content-type            express                 ipaddr.js               path-to-regexp          side-channel-weakmap
cookie                  finalhandler            is-promise              proxy-addr              statuses
cookie-signature        forwarded               math-intrinsics         qs                      toidentifier
debug                   fresh                   media-typer             range-parser            type-is
depd                    function-bind           merge-descriptors       raw-body                unpipe
dunder-proto            get-intrinsic           mime-db                 router                  vary
ee-first                get-proto               mime-types              safer-buffer            wrappy

# ディレクトリ数
ls -1 GoBash/ex-npm/node_modules | wc -l
      65

# ファイル数
find GoBash/ex-npm/node_modules -type f | wc -l
     601
```

※ `node_modules` は **65 ディレクトリ / 601 ファイル / 3.8M**（express 1個指定で 68 パッケージ）。
**「数万ファイル」とは言わない**。この構成では嘘になる。実測どおりに言う。

#### Goは別のプロジェクトが、同じパッケージを読んでいる

```text
Go
  プロジェクトA ─┐
               ├─→  ~/go/pkg/mod/github.com/go-chi/chi/v5@v5.3.2
  プロジェクトB ─┘     マシンに1つ・読み取り専用

npm
  プロジェクトA/node_modules/express   ← それぞれ別のコピー
  プロジェクトB/node_modules/express
```

- npm はプロジェクトごとにパッケージを入れている
- Go は別のプロジェクトでもパッケージを共有している

#### パッケージを共有して、困らないのか？

- なぜ同じディレクトリに配置しているのか？
- 別のプロジェクトが別のバージョン(例: v5.2.0)を使いたくなったら？
- 別のプロジェクトが中身を書き換えたら？
- コミッターが v5.3.2 のタグを別のコミットに付け替えたら？

### 3. 方針: 「go.mod, go.sum を手掛かり探る」

#### package.json は範囲、go.mod はピンポイント

Golangのimageファイル作成
<!-- https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=vscode&wt=none&l=text%2Fx-go&width=561&ds=true&dsyoff=0px&dsblur=0px&wc=true&wa=false&pv=0px&ph=0px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=%252F%252F%2520https%253A%252F%252Fgo.dev%252Fref%252Fmod%2523go-mod-file-module%250A%252F%252F%2520module%2520%25E3%2583%2587%25E3%2582%25A3%25E3%2583%25AC%25E3%2582%25AF%25E3%2583%2586%25E3%2582%25A3%25E3%2583%2596%25E3%2581%25AF%25E3%2580%2581%25E3%2581%259D%25E3%2581%25AE%25E3%2583%25A2%25E3%2582%25B8%25E3%2583%25A5%25E3%2583%25BC%25E3%2583%25AB%25E8%2587%25AA%25E8%25BA%25AB%25E3%2581%25AE%25E3%2583%25A2%25E3%2582%25B8%25E3%2583%25A5%25E3%2583%25BC%25E3%2583%25AB%25E3%2583%2591%25E3%2582%25B9%25E3%2582%2592%25E5%25AE%259A%25E7%25BE%25A9%25E3%2581%2599%25E3%2582%258B%250Amodule%2520github.com%252Ftakumashiraki%252Fconference-report%252FGoBash%252Fex-Go%250A%250A%252F%252F%2520https%253A%252F%252Fgo.dev%252Fref%252Fmod%2523go-mod-file-go%250A%252F%252F%2520go%2520%25E3%2583%2587%25E3%2582%25A3%25E3%2583%25AC%25E3%2582%25AF%25E3%2583%2586%25E3%2582%25A3%25E3%2583%2596%25E3%2581%25AF%25E3%2580%2581%25E3%2581%259D%25E3%2581%25AE%25E3%2583%25A2%25E3%2582%25B8%25E3%2583%25A5%25E3%2583%25BC%25E3%2583%25AB%25E3%2581%258C%25E3%2581%25A9%25E3%2581%25AE%2520Go%2520%25E3%2583%2590%25E3%2583%25BC%25E3%2582%25B8%25E3%2583%25A7%25E3%2583%25B3%25E3%2581%25AE%25E4%25BB%2595%25E6%25A7%2598%25E3%2583%25BB%25E6%258C%2599%25E5%258B%2595%25E3%2582%2592%25E5%2589%258D%25E6%258F%2590%25E3%2581%25A8%25E3%2581%2597%25E3%2581%25A6%25E6%259B%25B8%25E3%2581%258B%25E3%2582%258C%25E3%2581%25A6%25E3%2581%2584%25E3%2582%258B%25E3%2581%258B%25E3%2582%2592%25E7%25A4%25BA%25E3%2581%2599%250Ago%25201.26.8%250A%250A%252F%252F%2520https%253A%252F%252Fgo.dev%252Fref%252Fmod%2523go-mod-file-require%250A%252F%252F%2520require%2520%25E3%2583%2587%25E3%2582%25A3%25E3%2583%25AC%25E3%2582%25AF%25E3%2583%2586%25E3%2582%25A3%25E3%2583%2596%25E3%2581%25AF%25E3%2580%2581%25E4%25BE%259D%25E5%25AD%2598%25E3%2581%2599%25E3%2582%258B%25E3%2583%25A2%25E3%2582%25B8%25E3%2583%25A5%25E3%2583%25BC%25E3%2583%25AB%25E3%2581%25AB%25E3%2581%25A4%25E3%2581%2584%25E3%2581%25A6%25E5%25BF%2585%25E8%25A6%2581%25E3%2581%25A8%25E3%2581%25AA%25E3%2582%258B%25E6%259C%2580%25E4%25BD%258E%25E3%2583%2590%25E3%2583%25BC%25E3%2582%25B8%25E3%2583%25A7%25E3%2583%25B3%25E3%2582%2592%25E6%258C%2587%25E5%25AE%259A%25E3%2581%2599%25E3%2582%258B%250Arequire%2520github.com%252Fgo-chi%252Fchi%252Fv5%2520v5.3.2%250A -->

npmのimageファイル

package.json
<!-- https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=vscode&wt=none&l=application%2Fjson&width=320&ds=true&dsyoff=0px&dsblur=0px&wc=true&wa=false&pv=0px&ph=0px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=%257B%250A%2520%2520%252F%252F%25205.2.1%2520%25E4%25BB%25A5%25E4%25B8%258A%25206.0.0%2520%25E6%259C%25AA%25E6%25BA%2580%25E3%2582%2592%25E4%25BD%25BF%25E3%2581%2586%250A%2520%2520%2522dependencies%2522%253A%2520%257B%250A%2520%2520%2520%2520%2522express%2522%253A%2520%2522%255E5.2.1%2522%250A%2520%2520%257D%250A%257D -->

package-lock.json
<!-- https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=vscode&wt=none&l=application%2Fjson&width=266&ds=true&dsyoff=0px&dsblur=0px&wc=true&wa=false&pv=0px&ph=0px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=%257B%250A%2520%2520%252F%252F%25205.2.1%2520%25E3%2582%2592%25E5%2585%25A5%25E3%2582%258C%25E3%2581%259F%250A%2520%2520%2522node_modules%252Fexpress%2522%253A%2520%257B%250A%2520%2520%2520%2520%2522version%2522%253A%2520%25225.2.1%2522%252C%250A%2520%2520%257D%252C%250A%2520%2520%252F%252F%2520Node%252018%2520%25E4%25BB%25A5%25E4%25B8%258A%25E3%2582%2592%25E5%25AF%25BE%25E8%25B1%25A1%25E3%2581%25A8%25E3%2581%2599%25E3%2582%258B%250A%2520%2520%2522engines%2522%253A%2520%257B%250A%2520%2520%2520%2520%2522node%2522%253A%2520%2522%253E%253D%252018%2522%250A%2520%2520%257D%252C%2520%2520%250A%257D -->

公式仕様がそう書いている（[go.dev/ref/mod](https://go.dev/ref/mod#go-mod-file)）。

```text
module  … "A module directive defines the main module's path."
go      … "A go directive indicates that a module was written assuming the semantics of a given version of Go."
require … "A require directive declares a minimum required version of a given module dependency."
```

go.mod は「入れた物の一覧」ではなく、**このモジュールの自己申告**。書くのは人と `go get`、読むのは go コマンド。
依存を1つも書かないうちから `go mod init` で go.mod ができるのは、`module` 行（自分は誰か）が先に要るから（疑問1）。

```text
package.json       "express": "^5.2.1"    → >=5.2.1 <6.0.0-0 の範囲（npm 同梱 semver で確認）
package-lock.json  "version": "5.2.1"     → 範囲から選んだ結果の記録

go.mod             require github.com/go-chi/chi/v5 v5.3.2   → 幅を書かない。ピンポイント
```

- npm は「範囲の宣言（package.json）」と「結果の記録（lock）」の2枚で1組
- Go の go.mod は最初からピンポイントで版を宣言する。結果の記録を別に持たなくてよい
- → では go.sum は何を記録しているのか？ がスライド4への橋

##### スライド本文（案・8行以内）

```text
# go.mod（ex-Go）
module github.com/takumashiraki/…/ex-Go   ← 自分は誰か
go 1.26.8                                 ← どの Go の意味で書いたか
require github.com/go-chi/chi/v5 v5.3.2   ← 何を、どの版から要るか

# package.json（ex-npm）
"express": "^5.2.1"                       ← 幅。どれを入れたかは lock
```

##### トーク（約60秒）

> go.mod を開くと3行しかありません。仕様では、それぞれが defines / indicates / declares、つまり「自分は誰で、どの Go で書いて、何が要るか」の自己申告です。
> package.json も同じく自己申告ですが、`^5.2.1` は幅なので、実際にどれを入れたかは package-lock.json に記録するしかない。
> go.mod は幅を書かず、ピンポイントで版を書きます。じゃあ Go の go.sum は何を記録しているのか。ここから本題です。

##### 言わないこと・注意

- `require` は仕様上 **minimum** required version。実際に選ばれる版の決め方（MVS）はスライド4 ステップ2で扱うので、ここでは「どの版から要るか」と言うに留める。
  突っ込まれたらディスカッションに回す（ex-Go は chi 1個・chi は依存なしなので選ばれるのは v5.3.2 そのもの）
- 「だから lock が要らない」とは**このスライドでは言わない**。言い切るには MVS の決定性が要る。
  スライド4 ステップ2で公式の引用（"the build list is not saved in a "lock" file"）とともに言う
- `module` 行がスライド4 ステップ3（`<module>@<version>` のパス）とスライド5の答え合わせにもつながる: chi 自身の go.mod は `module github.com/go-chi/chi/v5` と宣言しており、
  これがそのままキャッシュのパス `github.com/go-chi/chi/v5@v5.3.2` になる。キャッシュには宣言だけの `v5.3.2.mod` も別に置かれている
  （= go.sum 2行目 `/go.mod h1:` のハッシュ対象。スライド4 ステップ4）。尺があれば1行で触れる

### 4. go get の中で何が起きているか（3:00-6:30）

##### このパートの組み立て

`go get` の出力4行と go.sum の2行を最初に見せ、分からない点を4つ（①〜④）残す。
その後 `go get` が中でやっている順（探す → 決める → 置く → 書く → 照合する）に1ステップ1枚で追い、各ステップで1つずつ回収する。

このパートで答える問い:

- どうやって依存関係を整理して `$GOMODCACHE` に入れるのか（ステップ1〜3）
- なぜ go.mod は範囲ではなくピンポイントなのか（ステップ2）
- go.sum の役割は何か（ステップ4・5）

#### go get の出力を1行ずつ読む

carbon.now.sh の画像1枚（8行以内。コマンド間の空行は詰める）

```bash
$ go get github.com/go-chi/chi/v5
go: downloading github.com/go-chi/chi/v5 v5.3.2   ← ① どこに置かれた？
go: downloading github.com/go-chi/chi v1.5.5      ← ② 頼んだのは1つなのに2つ
go: added github.com/go-chi/chi/v5 v5.3.2         ← ③ v5.3.2 はどう決まった？
$ cat go.sum
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
                                                  ← ④ なぜ2行？ h1: は何？
```

##### トーク（約30秒）

> `go get` を1回打つと、出力はこれだけです。でも1行ずつ読むと、分からないところが4つあります。
> `go get` が中で何をしているかを順番に追うと、この4つが全部解けます。

##### 注意

- 出力は question.md 疑問4 の実測。`downloading` の行はキャッシュが空のときしか出ないので、
  `GOMODCACHE=$(mktemp -d) go get github.com/go-chi/chi/v5` のように空のキャッシュで事前キャプチャする（ネットワーク依存のためライブではやらない）
- 「v5 を入れたのに？」のような煽りは入れない。出力を読み上げて、引っかかる点を指さすだけにする

#### ステップ1: モジュールを探す

```text
github.com/go-chi/chi/v5   ← この名前のモジュール？
github.com/go-chi/chi      ← このモジュールの中の v5/ ディレクトリ？
```

- import path のどこまでがモジュール名か、go は知らない
- 接頭辞ごとに proxy へ「最新版は？」と並列で問い合わせる
- → ② `chi v1.5.5` は候補として見に行ったもの

##### トーク（約30秒）

> `github.com/go-chi/chi/v5` と書いても、go には「`chi/v5` というモジュール」なのか「`chi` モジュールの中の `v5` ディレクトリ」なのか区別がつきません。
> なので両方を proxy に聞きに行きます。2行目の `chi v1.5.5` は、後者の候補を確かめに行った跡です。

##### 根拠

- go.dev/ref/mod「[Resolving a package to a module](https://go.dev/ref/mod#resolve-pkg-module)」:
  `golang.org/x/net/html` を探すとき `golang.org/x/net/html` / `golang.org/x/net` / `golang.org/x` / `golang.org` の latest を
  "in parallel" で要求する例が載っている
- **未確認**: 問い合わせだけで終わらず v1.5.5 の zip まで download される理由。
  「中に `v5/` パッケージがあるか確かめるため」は question.md の観測にもとづく解釈で、公式の明記は見つけていない。
  質疑で聞かれたら「観測からの解釈」と断る

#### ステップ2: バージョンを決める

```text
go get   → 最新版 v5.3.2 を選び、go.mod に1点で書く   ← ③ added
以降     → 各依存の go.mod を読み、要求された版のうち一番高い版を使う（MVS）
           新しい版が出ても結果は変わらない → 結果を lock に残す必要がない
```

- → ③ `added` は「go.mod の `require` に v5.3.2 を書いた」という意味
- 範囲で書かなくていいのは、選び方が go.mod だけで決まるから

##### トーク（約45秒）

> `added` は go.mod に書いたという意味です。では、なぜ範囲ではなく1点で書くのか。
> npm の `^5.2.1` は範囲なので、範囲のどれを選ぶかは実行した時点の最新版しだいです。だから選んだ結果を lock に残します。
> Go の `v5.3.2` は「この版以上が要る」という最低版の宣言で、実際に使う版は「要求された中で一番高い版」と決まっています。
> 材料が go.mod だけなので、新しい版が出ても結果は変わりません。結果を記録しておく必要がないから、1点で書けば足ります。

##### 根拠（スライドに英文を1つ引用してもよい）

- go.dev/ref/mod「[Minimal version selection (MVS)](https://go.dev/ref/mod#minimal-version-selection)」
  - "tracking the highest required version of each module"
  - "MVS is deterministic, and the build list doesn't change when new versions of dependencies are released"
  - "Unlike other dependency management systems, the build list is not saved in a "lock" file."
- go.dev/ref/mod `go get`: "required versions in `go.mod` files are _minimum versions_"

##### 言わないこと

- 依存グラフの図・`go mod graph`・「同じモジュールを違う版で2つ要求したら」の例は出さない（尺がない。ディスカッションで聞かれたら答える）
- `go get` が引数なしで最新版を選ぶ細かい規則（`@upgrade` 等）には触れない

#### ステップ3: 取得してキャッシュに置く

```text
$ ls ~/go/pkg/mod/github.com/go-chi/
chi  chi@v1.5.5          ← 2行目の downloading
$ ls ~/go/pkg/mod/github.com/go-chi/chi/
v5@v5.3.2                ← 1行目の downloading（読み取り専用で展開）
```

- → ① 答え: `~/go/pkg/mod`（= `go env GOMODCACHE`）の下の `<module>@<version>/`
- downloading の2行とも、ここに置かれている
- 展開前の go.mod と zip は `cache/download/.../@v/` に別々に保存（`v5.3.2.mod` / `v5.3.2.zip`）

##### トーク（約30秒）

> ①の downloading は、ここに置かれています。p10 で見た `~/go/pkg/mod` の下に、`モジュール名@バージョン` のディレクトリで、読み取り専用で展開されます。
> 2行目の `chi v1.5.5` も `chi@v1.5.5` として同じ場所にあります。
> 展開する前の go.mod と zip は別々のファイルで取ってきています。go.mod だけ別に取るのは、ステップ2で版を決めるのに go.mod しか読まないからです。

##### 補足

- 上の `ls` は空のキャッシュで `go get` し直して採取（2026-09-25 / go1.27.0。`GOMODCACHE` を一時ディレクトリに向けて実行）。
  `cache/download/github.com/go-chi/chi/v5/@v/` には `list` / `v5.3.2.info` / `.lock` / `.mod` / `.zip` / `.ziphash`、
  `cache/download/github.com/go-chi/chi/@v/` には `v1.5.5.*` が同じ構成で置かれた
- **普段使いのキャッシュでは再現しないことがある**: 手元の `~/go/pkg/mod/github.com/go-chi/` には `chi@v1.5.5` が無かった（`chi/v5@v5.2.5` / `v5.3.1` / `v5.3.2` のみ）。
  スライド用の `go get` と `ls` は、空のキャッシュで撮ったものを使う
- 取得元は `GOPROXY`（既定 `https://proxy.golang.org,direct`）。proxy の5エンドポイントの話は入れない（ブログへ）
- `@<version>` 入りのパスと読み取り専用は、スライド5の答え合わせで使う伏線

#### ステップ4: ハッシュを go.sum に書く

```text
github.com/go-chi/chi/v5 v5.3.2        h1:5YQk...   ← 本体（zip）のハッシュ
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tY...   ← go.mod 単体のハッシュ
```

- `h1:` = 中身の指紋（zip 内の全ファイルの SHA-256 を1つにまとめた値）
- go.mod が別の行なのは、ステップ2で go.mod だけ先に読むから
- → ④ 回収。`chi v1.5.5` は候補を見ただけなので go.sum に載らない（② も回収）

##### トーク（約45秒）

> 取ってきた中身からハッシュを計算して、go.sum に書きます。`h1:` はハッシュの方式の1番という意味で、中身の指紋だと思ってください。
> 2行あるのは、本体の zip と go.mod を別々に取ってきたからです。
> そして、ステップ1で見に行った `chi v1.5.5` は go.sum に1行もありません。キャッシュには残っているのに、です。
> go.sum に載るのは、実際にビルドに使うものだけです。

##### 補足

- `h1:` の計算: 各ファイルの sha256 と名前の行をソートして連結し、その全体の sha256 を base64 にしたもの（疑問4 ステップ4）。
  zip のバイト列そのもののハッシュではない。スライドでは「指紋」と言うに留める
- `キャッシュ ⊃ go.sum`: キャッシュは見に行ったもの全部、go.sum はビルドに関わるものだけ。オチの「実体は共有、信頼は共有しない」の伏線
- `go mod download`（引数なし）だと `/go.mod` 行しか書かれず、ビルドが止まる（疑問9 のおまけ・実測済み）。尺が余れば1行

#### ステップ5: go.sum で照合する

```text
次のダウンロードから   取ってきた中身のハッシュ ≠ go.sum → SECURITY ERROR
go.sum に無い初回      sum.golang.org に問い合わせて照合
                        （署名の検証鍵は go コマンドに内蔵）
```

- proxy も、他のプロジェクトも信用しなくていい

##### トーク（約30秒）

> go.sum に書いた値は、次に取ってきたときの答え合わせに使います。合わなければ `SECURITY ERROR` で止まります。
> 初めて取るモジュールは sum.golang.org に聞きます。この応答の署名を確かめる鍵は go コマンドに埋め込まれているので、proxy を信用する必要がありません。

##### 根拠・補足

- `go help module-auth`: "When the go command downloads a module zip file or go.mod file into the module cache, it computes a cryptographic hash and compares it with a known value ... Known hashes are stored in ... go.sum"
- Merkle / tile は言わない（ブログへ）。聞かれたら「Certificate Transparency と同じ透明性ログ」と1行で答える

#### go.sum は lock ファイルか？

- バージョンを決めるのは go.mod（ステップ2）
- go.sum は、取ってきた中身が本物かを確かめる（ステップ4・5）
- → go.sum はロックではなく、検証のための台帳

##### トーク（約20秒）

> 最初の問いに戻ります。go.sum は lock ファイルか。
> 版を決めているのは go.mod で、go.sum は中身を確かめているだけです。go.sum は lock ではなく、検証のための台帳です。

##### 質疑で聞かれたとき（実測済み・スライドには出さない）

- go.sum を退避して `go build` → `missing go.sum entry` で止まる（Go 1.16 以降 `-mod=readonly` が既定）
- `go mod tidy` で作り直す → `diff` で差分なし。「無いと困るが、版の決定には関与しない」の実測

### 5. パッケージを共有して、困らないのか？（6:30-8:15）

##### このパートの組み立て

p12 の4つの問いに、スライド4のステップで答える。柱の問い「なぜ他のプロジェクトを信用しなくて済むのか」はここで閉じる。

#### 答え合わせ

| p12 の問い | 答え | 根拠 |
| --- | --- | --- |
| なぜ同じディレクトリに置くのか | 同じ版の中身は変わらない。コピーする意味がない | ステップ4・5 |
| 別のプロジェクトが別の版を使いたくなったら | `<module>@<version>` なので別ディレクトリで共存 | ステップ3 |
| 別のプロジェクトが中身を書き換えたら | 読み取り専用。書き換わっても go.sum と合わない | ステップ3・4 |
| タグを別のコミットに付け替えたら | ハッシュが go.sum / sum.golang.org と合わず止まる | ステップ5 |

##### トーク（約60秒）

> 最初に出した4つの問いに戻ります。どれも、さっきのステップで答えが出ています。
> 版が同じなら中身は変わらないので、プロジェクトごとにコピーする意味がありません。版が違えばディレクトリが違うので共存できます。
> 書き換えやタグの付け替えは、go.sum のハッシュと合わなくなるので止まります。
> だから、同じキャッシュを読んでいる他のプロジェクトを信用する必要がありません。

##### 注意

- 「書き換えても go.sum と合わない」を検出するのは `go mod verify`（展開済みツリーを再ハッシュして go.sum と突き合わせる）。
  ビルドのたびに全ファイルを再ハッシュするわけではない（疑問4 ステップ5(c): 2回目以降は `.ziphash` との比較）。「アクセスのたびに照合」とは言わない
- タグ付け替えへの答えは、proxy / sum.golang.org が最初に記録した値を返し続けるという前提に立つ（疑問5）

#### 書き換えようとすると

キャプチャ1枚

```text
-r--r--r--@ ... mux.go
zsh: permission denied: .../chi/v5@v5.3.2/mux.go
```

##### トーク（約20秒）

> 実際に書き換えようとすると、そもそも書き込めません。`node_modules` を手で書き換えてデバッグする、という運用が Go に無い理由です。

### 6. まとめ（8:15-10:00）

#### npm との対比

connpass のタイトルが「npm との対比」で公開済みなので**必ず入れる**。

| | Go | npm | pnpm |
| --- | --- | --- | --- |
| 実体の場所 | `$GOMODCACHE` マシンで1つ | `node_modules` にプロジェクトごと展開 | グローバルストア + ハードリンク |
| 同一バージョンの重複 | 無い | プロジェクト数だけコピー | 無い |
| 実測（本サンプル） | 依存1 / 86ファイル / 588K（**外**） | 依存68 / 601ファイル / 3.8M（**内**） | - |
| 版の宣言 | go.mod にピンポイント（最低版） | package.json に範囲 + lock に結果 | 同左 |
| 書き込み | 読み取り専用 | 書き換え可能 | ストアは共有 |
| git 管理 | go.mod / go.sum だけ | 原則 `.gitignore` | 同左 |

##### トーク（約40秒）

> npm も `~/.npm/_cacache` に共有キャッシュを持っています。後から同じ設計に寄ってきたのであって、npm が雑なわけではありません。
> Go が最初からそうできたのは、版の中身は変わらないことと、版の選び方が go.mod だけで決まることを、仕様で決めたからです。

#### 共有しているのは「実体」だけ

```text
共有しているのは「実体」だけ。「信頼」は共有していない。
信頼の単位は、git 管理下の go.sum のまま。
```

##### トーク（約20秒）

> 共有しているのはキャッシュの実体だけで、信頼までは共有していません。信頼の単位は、それぞれのリポジトリの go.sum です。
> だから他のプロジェクトを信用する必要がありません。

#### みなさんに聞きたいこと

```text
- vendor/ 使ってますか？ どういう時に？
- GOPRIVATE、社内モジュールでどう運用してますか？
- node_modules を手で書き換えてデバッグした経験、ありますか？
```

ブログ導線（QR）もこの1枚に置く。

##### 狙い

- ディスカッション（20:40-）に投げる球にする。10分の登壇より、その後の45分のほうが効果が大きい

## 今回入れないもの（→ ブログへ）

| 内容 | 理由 |
| --- | --- |
| 疑問2・3（標準パッケージ / `IsStandardImportPath`） | 回収に寄与しない。面白いが本筋外 |
| 疑問6（proxy の5エンドポイント / VCS 直叩き） | 尺がない |
| 疑問4 ステップ5 の tile / Merkle 包含証明 | 1行の言及に留める |
| 疑問11（違う版を2つ要求したときの MVS / `go mod graph`） | 尺がない。MVS は定義1文だけステップ2で使う |
| dep / `Gopkg.lock` の歴史 | 未検証。柱を2本にすると10分に入らない |
| `go 1.27` の変更 | 未検証。ガイドラインの「推測で書かない」に反する |

## 会場条件からの制約

- **ターミナル出力は8行まで**。38F の部屋、40人規模。後列から読めるのは8行が限界
  - `.info` の JSON、sumdb の lookup 出力（7行＋署名）はそのままだと読めない → 該当行だけ抜いて拡大
  - go.sum の2行 と `-r--r--r--` の1行は、それぞれ1枚使う価値がある
- **`ls ex-Go` を出す前にビルド済みバイナリ `ex-Go/ex-Go`（7.9M）を消す**。
  `go build` の産物が `ls` に混ざると、3ファイルの絵が濁る。
  `question.md` / `README.md` も調査メモなので、スライドでは
  `main.go` / `go.mod` / `go.sum` の3行だけ抜いて見せる
- **`go clean -modcache` のライブ実行はしない**。166M の再取得が走る。キャプチャで
- ネットワーク依存のコマンド（`go get` / `curl proxy.golang.org`）も事前キャプチャ
- ライブでやる価値があるのは `ls` の2枚並べだけ

## 残り7日の進め方

```text
Day 1  疑問9 の実測 —— **完了**（第1段・第2段・おまけとも採取済み）
Day 2  ex-npm/ の実体を作る —— **完了**（express 5.2.1 / 68パッケージ / 実測は ex-npm/README.md）
Day 2-4 スライド作成（11枚）
Day 5  ターミナル出力の投影チェック（8行・フォントサイズ）
Day 6  通し。10分に収まらなければ「前提」パート（スライド3）から削る
Day 7  予備 / ブログ下書き（疑問3・6・10 の深掘りをここへ）
```
