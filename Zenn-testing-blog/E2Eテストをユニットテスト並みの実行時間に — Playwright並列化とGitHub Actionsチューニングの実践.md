# [E2Eテストをユニットテスト並みの実行時間に — Playwright並列化とGitHub Actionsチューニングの実践](https://zenn.dev/berry_blog/articles/39392e1da7ca71#7.-e2e-%E9%96%8B%E7%99%BA%E3%83%BB%E9%81%8B%E7%94%A8%E3%82%92%E6%94%AF%E3%81%88%E3%82%8B-ai-%E3%82%B9%E3%82%AD%E3%83%AB%E7%BE%A4)

E2Eテストの実行時間をかなり短縮したいケースで使えそう。
ただし手を動かすことや修正範囲は大きいので、導入前に検討したい。
AI(ClaudeやCodex)と会話しながら導入するもの、導入を見送るものを決める

## 従来型 E2E テストが抱える問題

> 実行時間が長く、開発フローに載らない: Google Testing Blog の [Just Say No to More End-to-End Tests](https://testing.googleblog.com/2015/04/just-say-no-to-more-end-to-end-tests.html)（2015）は、E2E 偏重の戦略では実行時間と flake が膨張し、フィードバックが遅く失敗時の原因特定も難しいと指摘し、ユニット 70 / インテグレーション 20 / E2E 10 のテストピラミッドを推奨しました。
> flaky（不安定）: Google は自社基盤で「全テスト実行の約 1.5% が flaky な結果を報告し、テスト全体の約 16% が何らかの flakiness を持つ」（[Flaky Tests at Google and How We Mitigate Them](https://testing.googleblog.com/2016/05/flaky-tests-at-google-and-how-we.html), 2016）、さらに「テストが大きい（E2E 的である）ほど flaky になる」（[Where do our flaky tests come from?](https://testing.googleblog.com/2017/04/where-do-our-flaky-tests-come-from.html), 2017）と報告しています。Microsoft の研究（[ISSTA 2019](https://www.microsoft.com/en-us/research/publication/root-causing-flaky-tests-in-a-large-scale-industrial-setting/)）でも社内 5 プロジェクトのテストケースの 4.6% が flaky で、flaky の再実行だけにテスト予算の数%〜十数%が消えるとされます。
> メンテナンスサイクルが回らない: [The Practical Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html)（martinfowler.com）は E2E を「悪名高いほど flaky で、失敗の多くが偽陽性」と評します。「失敗しても誰も見ない → 直らない → さらに信頼を失う」という悪循環で、E2E・手動テストに偏ってピラミッドが逆転した「テストのアイスクリームコーン」アンチパターンに陥りがちです。

## 並列化可能な E2E テストの構造設計

### 2.1 per-test fixture でデータ独立

> すべてのテストは、自分専用のデータを Factory 経由で Supabase に直接 INSERT して作り、テスト終了時に teardown が派生レコード（レビュー・承認・署名）から Storage のファイルまで回収
> 
> ポイントは 2 つ
> - shared state を一切持たない: 「テスト A が作った文書をテスト B が使う」を構造的に禁止。どのテストがどの worker でどの順に走っても互いに干渉しません。
> - UI でデータ準備をしない: UI 操作は「検証対象の操作」だけに使い、前提データは全て DB 直 INSERT。これは高速化（UI 経由のデータ準備は遅い）と安定化（前提部分の flaky 要因を排除）の両方に効きます。
> - UI 操作で生まれたデータも teardown に登録: アップロードなど UI 経由で生成されたエンティティは trackDocument(id) のような track 関数で明示的に回収対象へ登録します。

### 2.2 storageState でログインを省略

ログイン機能を提供しているAPIサーバーへの負荷が上昇するので、これを防ぐ狙いがある

### 2.3 レイヤー構造と POM（Page Object Model）

E2Eスイートは、下向きの依存だけを許すレイヤー構造を採用

```
spec → fixture → pages(POM) / helpers → repositories / precheckers / schemas → selectors / types
```

> spec から page.locator(...) を書くことを禁止。Locator と操作・待機ロジックは Page Object[[1]](https://zenn.dev/berry_blog/articles/39392e1da7ca71#fn-5d6a-1)に閉じ込め、spec は「S1: レビュー依頼する」「E1: ステータスが変わる」というシナリオ記述に徹します
> data-testid 文字列は selectors/ の定数に集約し、Page Object からのみ import
> DB 検証は repositories/（SELECT のみ）、事前データ検証は precheckers/ + zod schema

## GitHub Actions: Runner 設定とチューニング

※ GitLab なら GitLab CI/CD と置き換える

### 4.3 並列化が暴いたインフラの限界（1）: PostgREST プール枯渇

PostgREST を作るのではなく、テスト用のデータを注入するスクリプトを用意する方が良い？

## 疑問点整理

Q1. Playwright から「データ注入スクリプト」を呼んで実行してから E2E テストを走らせる、という構成は取れるか？

A. 取れる。Playwright の fixture からスクリプト（またはライブラリ関数）を呼び出す構成にすればよい。

典型パターンは次のいずれかになる。

- Playwright の fixture 内で ⁠node script.js⁠（または Go CLI）を ⁠child_process.exec⁠ で呼び出し、標準出力の JSON をパースしてテストに渡す

- 「スクリプト」を Node/TS の関数としてライブラリ化し、fixture から直接 import して呼び出す

- 共有データだけなら Playwright の ⁠globalSetup⁠ で一度だけ seed する（per-test 独立がいらない場合）

いずれも「テスト開始前に DB に必要なデータを用意し、その ID を spec で使う」という流れになる。

Q2. 自分は PostgREST を持っておらず、Golang 製 API サーバーしかない。この場合、データ注入は「API を叩くべきか」「DB 直 INSERT すべきか」どちらがよいか？

A. E2E 用のテストデータに限れば、HTTP API 直叩きより「DB 直 INSERT（またはそれに近い経路）」の方が現実的で扱いやすい。

理由は次の通り。

- HTTP API 直叩きの場合

 - 認証・認可の仕組みをテスト準備用にも一式通す必要がある

 - テストデータ作成のためだけに大量の HTTP リクエストを発生させることになる

 - そもそもその API 自体が E2E の検証対象なので、「前提条件のセットアップ」と「検証のための呼び出し」が同じレイヤーに混ざってしまう

- DB 直 INSERT の場合

 - HTTP/認証をバイパスできるため高速でシンプル

 - 「テスト専用の Factory / Repository」に DB 知識を閉じ込めれば、影響範囲を局所化できる

 - 記事で Supabase Factory がやっていること（Supabase クライアントで直 INSERT）と本質的に同じパターンになる

したがって、E2E の「前提データ作成」については、API サーバーは経由せず、DB 直 INSERT 相当の経路を持つ設計が望ましい。

Q3. 「DB 直 INSERT するスクリプト」が DB スキーマを知るのは、クリーンアーキテクチャ的に問題ではないか？

A. プロダクションコードとテストコードは役割が違うので、テスト側に「DB を知る薄い専用レイヤー」を持たせるのは実務上許容される。

ポイントは以下。

- 本番コードはクリーンアーキテクチャを維持する（UseCase / Repository / Entity を分離）

- テストコード側には「E2E 用のテストデータ Factory（あるいは test 専用 Repository）」を 1 箇所だけ作り、そこでだけ DB スキーマを扱う

- 「DB を知っているテスト用コード」は ⁠test/e2e/factories⁠ のような専用パッケージに閉じ込め、他のテストからは Factory の関数だけを呼ぶ

このように、「クリーンアーキの外側に、テスト専用インフラ層を薄く 1 枚追加する」と割り切れば、設計上の整合性は保てる。

Q4. HTTP API 直叩きと DB 直 INSERT の中間案として、どのような構成が取れるか？

A. Golang のユースケース／リポジトリを再利用するテスト用 CLI（またはライブラリ）を作り、そこから DB に書き込む構成が取り得る。

具体的なイメージは次の通り。

- 既存の Go アプリは、HTTP ハンドラの下に UseCase 層・Repository 層を持っている想定

- テスト用に ⁠cmd/seed⁠ のような CLI を追加し、その中で HTTP サーバーは起動せず ⁠InitContainer⁠ だけ呼んで UseCase を直接叩く

- CLI からは UseCase 経由で DB にデータを作成し、その結果（ID 等）を JSON で標準出力に出す

- Playwright 側の fixture は ⁠exec('go run ./cmd/seed ...')⁠ 的に CLI を叩き、JSON を読み取ってテストに渡す

この方式だと、

- HTTP / 認証はバイパスするので速くてシンプル

- しかし DB へのアクセスやビジネスルールは本番コードを経由するため、スキーマ変更にも比較的強い

「素の DB 直スクリプト」と「HTTP API 直叩き」の中間であり、クリーンアーキの構造を活かしたテストデータ注入パスになる。

Q5. Playwright との接続はどのように設計するのがよいか？

A. 「per-test fixture でテスト専用 CLI / Factory を呼び出してデータを作成し、その ID をテストに渡す」構成がよい。

基本パターンは次。

- Playwright 側に test.extend⁠ でカスタム fixture（例: ⁠seededDataId⁠）を定義

- その fixture の中で

 - Go 製 seed CLI を ⁠child_process.execFile⁠ で呼ぶ

 - あるいは Node/TS の seed 関数を直接呼ぶ

- CLI/関数の標準出力（JSON）や返り値から ID を取り出し、⁠use(id)⁠ で spec に渡す

- 必要に応じて fixture の teardown でクリーンアップ用の CLI/関数を呼ぶ

これにより、各テストは

1. fixture によって自分専用の前提データを DB に作成してもらい

2. その ID を使って UI を操作・検証する

という形になり、「並列実行可能で shared state を持たない E2E」というこの記事の前提とも整合する。