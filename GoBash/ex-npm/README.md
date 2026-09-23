# ex-npm

`ex-Go` と同じことをする Node 実装。登壇スライド 2(a)「`ls` を2つ並べる」の npm 側。

```text
ex-Go   … net/http + github.com/go-chi/chi/v5
ex-npm  … node:http 相当 + express
```

## 再現手順

```bash
cd GoBash/ex-npm
npm install express
```

## 環境

```text
node      v26.2.0
npm       11.13.0
express   5.2.1
採取日    2026-09-23
```

`node_modules/` は `.gitignore` 済み。リポジトリには `package.json` /
`package-lock.json` / `index.js` だけが入る。

---

## 実測（2026-09-23 採取）

### ls を2つ並べる

```bash
ls ex-npm
```

```text
README.md             ← このメモ。スライドでは省く
index.js
node_modules          ← これ
package-lock.json
package.json
```

```bash
ls ex-Go
```

```text
README.md             ← 調査メモ。スライドでは省く
go.mod
go.sum
main.go
question.md           ← 調査メモ。スライドでは省く
```

`ex-Go` 側に**実体が入るディレクトリが存在しない**。これがスライド 2(a) の絵。

※ 両側とも調査メモ（`README.md` / `question.md`）が混ざる。また `ls ex-Go` には
`ex-Go`（`go build` が吐いた 7.9M のバイナリ・`.gitignore` 済み）が出ることがある。
投影前に `rm ex-Go/ex-Go` し、スライドでは下の4行 / 3行だけ抜いて拡大する。

```text
ls ex-npm                    ls ex-Go
-----------------------      -----------------------
index.js                     main.go
node_modules          ←      go.mod
package-lock.json            go.sum
package.json
```

これがスライドに載せる最終形。

### node_modules の中身

```bash
ls ex-npm/node_modules | wc -l       # 65
find ex-npm/node_modules -type f | wc -l   # 601
du -sh ex-npm/node_modules           # 3.8M
```

`npm install express` は **1 個指定して 68 パッケージ**を持ってくる。

```bash
ls ex-npm/node_modules | head -10
```

```text
accepts
body-parser
bytes
call-bind-apply-helpers
call-bound
content-disposition
content-type
cookie
cookie-signature
debug
```

### chi の実体（リポジトリ外）

```bash
find "$(go env GOMODCACHE)/github.com/go-chi/chi/v5@v5.3.2" -type f | wc -l   # 86
du -sh "$(go env GOMODCACHE)/github.com/go-chi/chi/v5@v5.3.2"                 # 588K
```

```bash
cd ex-Go && go list -m all
```

```text
github.com/takumashiraki/conference-report/GoBash/ex-Go
github.com/go-chi/chi/v5 v5.3.2
```

### 対比表

| | ex-Go (chi) | ex-npm (express) |
| --- | --- | --- |
| リポジトリ内の実体 | **無い** | `node_modules/` |
| 依存パッケージ数 | 1 | 68 |
| 実体のファイル数 | 86（リポジトリ**外**） | 601（リポジトリ**内**） |
| 実体のサイズ | 588K（`$GOMODCACHE`） | 3.8M（`ex-npm/` 直下） |
| git 管理するもの | `go.mod` / `go.sum` | `package.json` / `package-lock.json` |

## 注意: 数字の使い方

- **「数万ファイル」は誇張**。この構成では 601 ファイル / 3.8M。実測どおりに言うこと
- **`du` の比較はしない**。`ex-Go` は `question.md` とビルド済みバイナリで 8.0M あり、
  `ex-npm`（3.8M）より大きく出る。サイズで殴ると自滅する
- この枠で効くのは**サイズではなく場所**。「リポジトリの中か、外か」だけを言う
- 「1 個入れたら 68 個来た」は面白いが**別の軸（依存グラフの深さ）**。
  1枚1メッセージの原則から、入れるなら 6（対比表）で1行。2(a) では言わない
