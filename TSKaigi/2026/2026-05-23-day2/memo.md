# Day2

## [400超のデータポイントを型で制す —— UPSIDERの与信審査エンジンを支えるTypeScriptの「柔軟性」と「結合力」（泉雄介）](https://2026.tskaigi.org/talks/39)

- x
  - https://x.com/yizumi
- 資料
  - 

- 与信のルールが変わるため、堅牢でアジリティのある言語である必要があった
  - typescriptとzodがよかった
    - zodバリデーションできるのが嬉しい
  - javaだと複数のファイルを使いdtoを作ったりする必要がある
- 動的マージと静的型の両立
- Strategyパターン
  - スタートアップ や 上場企業 で審査パターンが違う
  - switch 文で対応する
- openapi生成モデルとの訣別
  - tRPC + zod

## [React の props は値の集合ではない — UI の状態を宣言するコンポーネント設計（nabeliwo）](https://2026.tskaigi.org/talks/41) 🔥🔥

- 資料
  - https://nabeliwo.github.io/slides/talks/20260523_tskaigi-2026_react-props

メモ

- propsがUIの状態を宣言する
  - 判別可能な union を明確にする
  - error , loading, success の３つのunionで表現する
  - 責務の分離をする
    - jsx で成功、 ロード、失敗の状態管理
    - userlist でdataフェッチ

## [型プラグインシステムの実装に使われるテクニック（elecdeer）](https://2026.tskaigi.org/talks/42)

https://t.co/skWaXzyyuD

## [TypeScriptでドット絵エディタ実装録: 状態設計と実装判断（つくだに）](https://2026.tskaigi.org/talks/49)

## [LLM時代のリファクタリング戦略：AIエージェントによる段階的・安全なTS移行方法（市川 賢）](https://2026.tskaigi.org/talks/45) ☑️

## [AI Agent に“攻略本”を渡したら、150フォームの移行が回り始めた話（高橋悟生）](https://2026.tskaigi.org/talks/80)

- 判断基準を教える

## [AI活用の格差をなくす：チーム全体のAI開発生産性を底上げする方法（中津川篤司）](https://2026.tskaigi.org/talks/78) 🔥

## [AIコーディングエージェントの活用で、コードは静かに肥大化した —— 型・Lint・Skillsで挽回する10分（篠田 陽介）](https://2026.tskaigi.org/talks/82)


## [「バイトル」のTypeScriptリニューアル — 積み上がったレガシーとパフォーマンスに挑む現在地（横山 隼）](https://2026.tskaigi.org/talks/86)

- 資料
  - 

メモ
- よくあるベストプラクティス的なパフォーマンス改善でも効果があった

## [TypeScriptで実現する既存APIを活用したリモートMCPサーバー構築（鈴木翔大）](https://2026.tskaigi.org/talks/84) 🔥

- 資料
  - https://speakerdeck.com/soarteclab/tskaigi-2026

## [Next.js × OpenAPIで型安全なデータ境界を設計するアーキテクチャ改善（土本祐介）](https://2026.tskaigi.org/talks/79) ☑️

- 資料
  - 

メモ
- プロダクトで使っている技術なので気になる

## [TypeScriptでPlatform SDKを作る技術（樋口 彰）](https://2026.tskaigi.org/talks/52) ☑️

## [Real World Effect-TS: 堅牢なプロダクトを型で組み上げる（asa1984）](https://2026.tskaigi.org/talks/50)

- 資料
  - 

メモ

- 関数型ドメインモデリング
  - 型でドメインモデリングする
  - 関数で状態遷移する
  - エラーを値として扱う
  - ドメインロジックを純粋にする
  - ワークフローを合成
- 言語機能の制約を使う
  - 自由度が高いのが辛い
  - Go, Rustは静的かたづけ
  - Effect-TSというフレームワークの進め？
  - アプリ開発に必要なオールインワンのもの
- Effect-TS
  - ユーティリティー
  - effect runtime
    - 成功の値, err型, 依存型
- Effect型にすると嬉しいこと
  - エラー型が合成される
    - エラーのログを運用監視で見れるようにする
    - エラーを値にすることが必ず正しいとは限らない
    - 握りつぶすこともいい
- DI
  - 依存のバケツリレーが発生する
    - Reactのpropsのバケツリレーみたいな感じ useContext が解決するやつ
  - DIコンテナを使おう
    - Effect-TSでは依存関係が解消されるまでneverでエラーになるようにしてる
  - 

## [AIエージェントと協働するCLI開発 — BunとOpenClawで学んだこと（yoshikouki）](https://2026.tskaigi.org/talks/59)

- 資料
  - https://speakerdeck.com/yoshikouki/aiezientotoxie-dong-suruclikai-fa-buntoopenclawdexue-ndakoto

メモ

- コマンドを使えるようにする
  - help
  - Log で何が起こっているかをAIが理解を使えるようにしておく
- Bun は必要な機能が豊富
- 1年後
  - ifや機能はどうなっているかわからないが
  - 利用者とデータ・品質は残っているはず
- AIエージェントの効用
  - 自分の書き方を言語化する
    - 再学習できる
  - 変わりにくいもの自身の関心を再学習できる

## [自動レビューエンジンの実装と運用 ~レビューのない世界へ~（Kanato）](https://2026.tskaigi.org/talks/60)

- 資料
  - https://speakerdeck.com/kurukuru1999/zi-dong-rebiyuenzinnoshi-zhuang-toyun-yong-rebiyunonaishi-jie-he

メモ

- 実装速度が上がり、同時に何個かタスクを持つことになった
  - 対策
    - AIによるレビューをする
  - codexにレビューさせてcodexに修正させる
    - レビュー内容を記載
    - 修正内容を記載
    - prコメントを記載
  - 観点
    - pr内の責務ないの実装か
    - 重複コードがないか
    - セッションは新しくする
  - コメントが残るので、プロンプトの良し悪しが振り返れる
    - 15往復くらいやり取りが起こる
  - code-policeというものを作った
    - コード品質を見てくれるものを作った
      - 本質的なバグを見てくれるようになった
    - 実際の内容については紹介されてなかった

## [Hono RPCとDrizzle ORMで実現する、AIにも優しいTypeScriptファーストな開発（Yudai Shinnoki）](https://2026.tskaigi.org/talks/61) 🔥

- 資料
  - 

メモ

- コード生成を挟まずにフルTSで実装できるので、AIと相性がいい
  - 設計思想やライブラリの引き出しを増やすことで、現状を超えた改善の余地を知る
  - 経験に基づくtipsやはまりどころを小顔
- HonoRPC
  - 実装コードを書くと型がつく
- Drizzle
  - スキーマをtsで定義できる
  - SQLっぽい描き恥になる
  - Drizzle kit でマイグレーション管理も可能
- Next.jsから直接DBを接続していた
  - 脆弱性がある前提で考慮が必要になった
  - Aurora Data APIを使うとなんかできる？ 🤷‍♂️
    - コストを下げられた？
    - 使わなかった
    - VPC超えることができる
      - 外から直接アクセスができる
- PGliteを入れた
  - 何が嬉しいのかわからなかった 🤷‍♂️
  - PostgreSQL を WebAssembly(WASM) 化して、ブラウザやNode.js内で動かせるようにしたもの
    - SQLite 的な手軽さ
    - PostgreSQL の機能
    - TypeScript エコシステム

## [TypeScriptバックエンドのオブザーバビリティ戦略 — Datadog × NestJSの実践（山本大星）](https://2026.tskaigi.org/talks/64) 🔥

- 資料
  - https://speakerdeck.com/taiseiyamamotoan/tskaigi-2026-typescriptbatukuendonoobuzababiriteizhan-lue-datadog-x-nestjsnoshi-jian

メモ

- 運用の辛さ
  - エラーやアラートの原因がすぐにわからないのが辛い
    - logで三秒で原因がわかったら楽しい
    - 実例
      - アラートみて、logみて、のように原因を調査する
  - 原因特定のための仕込み = オブザーバービリティー
  - ベテランと新人の違い
    - 新人
      - 昔からあるからスルーしてオッケー
      - 深夜バッチでいつも出るからオッケー
    - ベテラン
      - 何これ
      - 対応した方がいいの？
      - 緊急度わからない
  - Datadogを使う前
    - AWS最小構成で
  - Datadogを使う(現在)
    - monitorやAPM、error_tracking、log, syntheticsをdatadogが提供しているものを使えるようになった
    - apm
      - アプリのトレースが見れる
    - error_tracking は エラーの集約と分析ができる
      - エラーをignoreできる
    - dash_borad
  - ツールを入れるだけではだめ
    - アプリケーションの改修をしないといけない
  - 気負わないことを目標にする
    - 前半は入門
    - 後半はDatadogのTips
- 入門
  - ログの構造化とは
    - 人なら単純な文章をアウトプット
    - json構造化する
      - 機械向けに整形する
    - jsonのkey:valueを解釈して可視化してくれる
  - 設定例
    - main.tsに設定を入れるだけでいいらしい
      - localとproductionでだしわけする
    - NestJSの例ではあるが参考できそう
  - エラーのプロパティー
    - エラークラス名、コンテキスト、causeを追加する
      - どこで、何をしてる時にエラーが出たのかを記載する
      - causeは詳しく書いてあるもの
    - コンテキスト(context)にエラーの内容を保持する
- DatadogのTips
  - Datadogのライブラリに巻き込まれた話
    - ライブラリを入れてapmを動かす
    - アラートがなる
    - リバートしたら治った
    - apmのプラグインが悪さしていた
      - 各middleware/ライブラリの詳細をトレーシングできる
      - 再帰処理
    - プラグイン有効化は慎重にする
  - dd-traceをちょうどよく依存する
    - 依存しすぎないようにして、今後もっといい製品が出たら乗り換えること検討
    - dd-traceが使えない箇所を独自実装する必要があった
      - 自作ライブラリを挟んで、
        - SQS受信やメッセージ処理を抽象化
        - spanの生成やタグづけ
  - **ログとトレースの紐付けが大事**
    - ログにspan_idとtrace_idを入れて、トレースの紐付ける
      - 自作のライブラリで span_id と trace_id を付与する

## [いつテストを書くか？―ソフトウェア開発における安心と不安について考える（lacolaco）](https://2026.tskaigi.org/talks/63)

- 資料
  - https://docs.google.com/presentation/d/e/2PACX-1vSuTOMT21i3AsVT0xiy9wZV1h2ZytSbFKjEFzcslLvmv2sHQ2VNgPRGXiEQid32557D1wxlB_XmlIpT/pub?slide=id.SLIDES_API1416270684_0

## [Polymorphic Components パターンで作る、型安全でセマンティックな UI コンポーネント（ryo）](https://2026.tskaigi.org/talks/55) ☑️

- 資料
  - https://t.co/Ab1Azz8fJs

## [TypeScript7 - 非推奨設定から読む責務の変化（Ayu）](https://2026.tskaigi.org/talks/71)

- 資料
  - 

メモ

- source code
- parser

## [TypeScript Compiler はどのように未使用変数を検出しているのか？（つねみ@tocomi）](https://2026.tskaigi.org/talks/72)


## [string地獄を脱出する — ValueObject + Zod 実践パターン（高田 理功）](https://2026.tskaigi.org/talks/70) 🔥

- 資料
  - https://sansan.box.com/s/j4xkoqhxvrcskder3r08t02madbaxcdn

メモ

- ValueObject
  - 同一性に基づかない等価性を持つ
- バリデーションを別のクラスに責務を委譲する
  - スッキリ描ける。しかし、実装が大変
  - バリデーションライブラリーを導入する
    - zod
  - 静的解析ジェネレター
    - Plop
      - ValueObject を自動生成できる
