# Day1

## tscからtsgoへ ── DenoのTypeScript基盤はどう変わったか

## [TypeScriptの「型」をAIのスキルに昇華させてみた件について（higak9）](https://2026.tskaigi.org/talks/7)

- AIの型知識チェック
  - 基本型、複合型、構造型、ジェネリクス、条件型、Utilityなど
  - 知識面はある程度知ってそう
- **書き方によってAIの解釈が変わるか？**
  - ストリングとリテラルの使い分け
    - ストリングだと文字列としか予測できない
    - リテラルなら予測できる
    - 型ガードなしでアサーションするのは危険
      - 型チェックを放棄するような書き方は危険
- 登壇者目線のプラクティス
  - 型の使い方の使い所と使い分けがわからない 🤷‍♂️
  - **判別可能なユニオンとneverで完全性を検証** 🤷‍♂️
  - **Mapped TypesによるDRYな書き方** 🤷‍♂️
  - **ROROパターン** 🤷‍♂️
  - **narrowingによる型の絞り込み** 🤷‍♂️
  - **ブランド方で構造体を区別する** 🤷‍♂️
- 実験
  - 観点
    - 型の拡張でどのような機能追加
    - 型定義のリスク
  - 高度な型の方が具体的な形や仕様や機能拡張できた

## [TypeScriptだけでAIエージェントを作る ― フロント・エージェント・インフラのフルスタック実践（福地開）](https://2026.tskaigi.org/talks/20) ☑️

- 資料
  - https://t.co/X7zz8QolYg

TODO: 登壇を見れていないので、資料を見る
- AI向けには登壇資料をもとに説明をして欲しい 🤷‍♂️

## [TanStack StartのcreateServerFnで作る、型が通るAPI（Yuki Terashima）](https://2026.tskaigi.org/talks/25) 🔥

型がとっているが、壊れやすい設計の対策

serverFnは薄くする
出力もSchema化する

- 資料
  - https://t.co/wBVDRhM2iK

- TanStack Start がよくわかってないので、Next.js(useSWR)と比較して教えて欲しい🤷‍♂️

## [実践TanStack Start: 新規プロダクトを開発して確立した、サーバーとクライアント境界の設計パターン（Shimmy）](https://2026.tskaigi.org/talks/26) 🔥

- next と TanStart の棲み分け
  - loader/SeverFunciton を自分でおく
- 境界
  - loader
    - クライアントからデータを受け取る境界
    - 全てのデータを受け取るのではなく、最低限の内容を受け取る
  - SeverFunciton
    - 処理リクエストを投げる境界
    - 横断的な関心ごとをmiddlewareに任せる
    - 境界として薄くする
    - 下記に移譲する
      - i/o(リポジトリ)
      - バリデーション
      - ビジネスロジク
      - ミドルウェア

- 資料
  - https://t.co/yVmwDAdniV

## [TanStack Router の型定義を読み解く（IORI）](https://2026.tskaigi.org/talks/27) 🔥

- https://x.com/Yz_Iori

- ルート定義そのものが型を生み出す
  - ルート定義の型が理解できてない 🤷‍♂️

## [Zod v4 Codec でスキーマに型変換を埋め込む REST API 設計（Ryutaro Yako）](https://2026.tskaigi.org/talks/31) ☑️

- 資料
  - https://speakerdeck.com/ryutaro_yako/zod-v4-codec-desukimanixing-bian-huan-womai-meip-mu-rest-api-she-ji-number-tskaigi2026

- 登壇を見れていないので、資料を見る
- AI向けには登壇資料をもとに説明をして欲しい 🤷‍♂️

## [アンチパターンを避ける型駆動React最適化（Kazuya Serizawa）](https://2026.tskaigi.org/talks/17) ☑️

- 資料
  - https://t.co/Bh4nfqfJls

## ハンズオン 🔥🔥

https://typescript-jpc.connpass.com/event/392953/

- 資料
  - https://tskaigi2026-handson.berlysia.workers.dev/#inference-limit-handler