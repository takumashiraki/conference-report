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
- MVS / dep 史 / go 1.27 は**今回入れない**（未検証、かつ尺に入らない）
- 10分で完結させない。**ディスカッション30分に投げる球を残して終わる**
- 1枚1メッセージ。11枚前後、1枚あたり約50秒

## タイムテーブル（10分）

<!-- 下記は案 -->
| 時間 | 内容 | ネタ元 |
| --- | --- | --- |
| 0:00 | 自己紹介 + OPTiM（**30秒・1枚**。ここで2分使うと本題が8分になる） | - |
| 0:30 | 【問い】`ls` を2つ並べる / `go env GOMODCACHE` で即答 / **2プロジェクト共有の絵** / 本当の問いを提示 | 疑問5 |
| 2:00 | 【前提】package.json は幅、go.mod はピンポイント（go.mod は「宣言」） | 疑問1（圧縮） |
| 3:00 | 【go.sum】生まれる瞬間 → 2行の意味 → **消してもバージョンは変わらない** | 疑問4 / 7 / 9 |
| 6:00 | 【回収】共有キャッシュが成立する3条件 + read-only の実演 | 疑問5 |
| 8:30 | 【対比】npm / pnpm / Go の1枚 | 疑問5 |
| 9:30 | 【オチ】「実体は共有、信頼は共有しない」+ 導線 + 議論への問い | - |

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

- `require` は仕様上 **minimum** required version。実際に選ばれる版の決め方（MVS）は今回入れない方針なので、「どの版から要るか」と言うに留める。
  突っ込まれたらディスカッションに回す（ex-Go は chi 1個・chi は依存なしなので選ばれるのは v5.3.2 そのもの）
- 「だから lock が要らない」とは**このスライドでは言わない**。言い切るには MVS の決定性が要る（疑問12・未検証）。
  スライド4の「go.sum を消しても変わらない」という実測で見せる
- `module` 行が回収（スライド5 条件2）にもつながる: chi 自身の go.mod は `module github.com/go-chi/chi/v5` と宣言しており、
  これがそのままキャッシュのパス `github.com/go-chi/chi/v5@v5.3.2` になる。キャッシュには宣言だけの `v5.3.2.mod` も別に置かれている
  （= go.sum 2行目 `/go.mod h1:` のハッシュ対象。スライド4(b)）。尺があれば1行で触れる

※ 標準パッケージの話（疑問2・3）は**全カット**。回収に寄与しない。

#### go.sum は lock ファイルか？

**(a) go getする**

```bash
$ go get github.com/go-chi/chi/v5
go: downloading github.com/go-chi/chi/v5 v5.3.2
go: downloading github.com/go-chi/chi v1.5.5
go: added github.com/go-chi/chi/v5 v5.3.2

$ cat go.sum
github.com/go-chi/chi/v5 v5.3.2 h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
```

**(b) 2行ある理由**

- 1行目 = モジュール本体（zip）のハッシュ
- 2行目 = その go.mod 単体のハッシュ
- 依存解決では go.mod だけ先に読むから、別々に要る

**(c) 転回点 ★この枠の山**

2段構えで見せる。**先に「無いと止まる」を見せてから「再生成しても同じ」を見せる**。

第1段: go.sum が無いとビルドが止まる（実測済み）

```bash
mv go.sum /tmp/go.sum.bak
go build ./...
```

```text
main.go:8:2: missing go.sum entry for module providing package
github.com/go-chi/chi/v5 (imported by .../GoBash/ex-Go); to add:
        go get github.com/takumashiraki/conference-report/GoBash/ex-Go
```

※ Go 1.16 以降は `-mod=readonly` がデフォルトなので、`go build` は
go.sum を勝手に書かずエラーで止まる。古い記事の「再生成される」は現行と違う。

（余裕があれば1行）`to add:` が chi ではなく**メインモジュール自身**を指しているのが面白い。

第2段: 再生成しても中身は変わらない（実測済み）

```bash
go mod tidy
diff /tmp/go.sum.bak go.sum   # → 差分なし
```

※ 再生成は **`go mod tidy`** を使う。`go mod download`（引数なし）だと
`/go.mod` 行しか書かれず、ビルドはまだ通らない（下記）。

> **go.sum が無いとビルドは通らない。でも再生成すると 1 バイトも変わらない。
> つまり go.sum は必須だが、**バージョン決定には関与していない**。
> バージョンを決めているのは go.mod。go.sum はロックではなく検証。**

この二段構えが `package-lock.json` との差を一番はっきり出す。
（「無くてもいい記録」ではない。しかし解決結果を変える力も持たない）

**おまけ（余裕があれだ30秒。なければブログへ）**

`go mod download`（引数なし）で再生成すると、書かれるのは 1 行だけ。

```text
github.com/go-chi/chi/v5 v5.3.2/go.mod h1:R+tYY2...   ← /go.mod 行のみ
→ go build は依然として missing go.sum entry で止まる

go mod download github.com/go-chi/chi/v5  とモジュールを明示すると
h1: 行も書かれ、ビルドが通る
```

これはスライド 4(b) の「2 行ある理由」の**生の実演**になっている。

```text
モジュールグラフの構築にしか使わない → /go.mod の 1 行だけ
ビルドに実際にソースが必要   → h1: と /go.mod の 2 行
```

※ `go help mod download` にも go.dev/ref/mod にもこの差の明記は無い。
上記は実測結果と、疑問4 で確認した「2 行 / 1 行」ルールからの解釈。

ハッシュの正しさは `sum.golang.org` が保証している。
検証鍵は go コマンドのバイナリに埋め込まれているので、proxy を信用する必要がない。
（Merkle / tile の話は**1行だけ**。深入りしない）

### 5. 回収（6:00-8:30）

スライド2(d) の問い —— **なぜ他のプロジェクトを信用しなくて済むのか** —— に答える。
共有キャッシュが成立する3条件。

```text
1. バージョンが immutable
   → 作者が後からタグを付け替えてもビット列は変わらない（checksum database が担保）
   → 同じ物をプロジェクトの数だけコピーする意味がない

2. パスが <module>@<version>
   → 別バージョンは別ディレクトリとして自然に共存する
   → 「別リポジトリで違うバージョン」で困らない

3. アクセスのたびに main module の go.sum と照合する
   → 他のプロジェクトを信頼する必要がない
```

さらにパーミッションで二重に守られている（キャプチャを1枚で見せる）。

```text
-r--r--r--@ ... mux.go
zsh: permission denied: .../chi/v5@v5.3.2/mux.go
```

→ `node_modules` を手で書き換えてデバッグする運用が Go に無い理由。

### 6. npm との対比（8:30-9:30）

1枚。connpass のタイトルが「npm との対比」で公開済みなので**必ず入れる**。

| | Go | npm | pnpm |
| --- | --- | --- | --- |
| 実体の場所 | `$GOMODCACHE` マシンで1つ | `node_modules` にプロジェクトごと展開 | グローバルストア + ハードリンク |
| 同一バージョンの重複 | 無い | プロジェクト数だけコピー | 無い |
| 実測（本サンプル） | 依存1 / 86ファイル / 588K（**外**） | 依存68 / 601ファイル / 3.8M（**内**） | - |
| 書き込み | 読み取り専用 | 書き換え可能 | ストアは共有 |
| git 管理 | go.mod / go.sum だけ | 原則 `.gitignore` | 同左 |
| 検証 | アクセスのたびに go.sum と照合 | install 時中心 | install 時中心 |

npm も `~/.npm/_cacache` に共有キャッシュを持っている。
**後から同じ設計に寄ってきた**のであって、npm が雑なわけではない。
Go が最初からそうできたのは「バージョンが immutable」を仕様で決めたから。

### 7. オチ + 議論への問い（9:30-10:00）

```text
共有しているのは「実体」だけ。「信頼」は共有していない。
信頼の単位は、git 管理下の go.sum のまま。
```

→ だから **他のプロジェクトを信用する必要がない**。2(d) の問いはここで閉じる。

最後の1枚は**ディスカッション(20:40-)に投げる球**にする。

```text
みなさんに聞きたいこと
  - vendor/ 使ってますか？ どういう時に？
  - GOPRIVATE、社内モジュールでどう運用してますか？
  - node_modules を手で書き換えてデバッグした経験、ありますか？
```

ブログ導線（QR）もここ。10分の登壇より、その後の45分のほうが効果が大きい。

## 今回入れないもの（→ ブログへ）

| 内容 | 理由 |
| --- | --- |
| 疑問2・3（標準パッケージ / `IsStandardImportPath`） | 回収に寄与しない。面白いが本筋外 |
| 疑問6（proxy の5エンドポイント / VCS 直叩き） | 尺がない |
| 疑問4 ステップ5 の tile / Merkle 包含証明 | 1行の言及に留める |
| 疑問11・12（MVS / なぜ lock が要らないか） | **未検証**。柱を2本にすると10分に入らない |
| dep / `Gopkg.lock` の歴史 | 同上 |
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
