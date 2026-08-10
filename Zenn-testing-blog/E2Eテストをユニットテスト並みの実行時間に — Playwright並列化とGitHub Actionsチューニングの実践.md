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

> 