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

- propsがUIの状態を宣言する
  - 判別可能な union を明確にする
  - error , loading, success の３つのunionで表現する
  - 責務の分離をする
    - jsx で成功、 ロード、失敗の状態管理
    - userlist でdataフェッチ

## [型プラグインシステムの実装に使われるテクニック（elecdeer）](https://2026.tskaigi.org/talks/42)

https://t.co/skWaXzyyuD

## [TypeScriptでドット絵エディタ実装録: 状態設計と実装判断（つくだに）](https://2026.tskaigi.org/talks/49)

                                              

## [Hono RPCとDrizzle ORMで実現する、AIにも優しいTypeScriptファーストな開発（Yudai Shinnoki）](https://2026.tskaigi.org/talks/61)

## [Polymorphic Components パターンで作る、型安全でセマンティックな UI コンポーネント（ryo）](https://2026.tskaigi.org/talks/55)
