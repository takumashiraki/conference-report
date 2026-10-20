# パッケージ改ざん攻撃への対応メモ（Go / npm）

調査日: 2026-10-19〜20。記事 `zenn-article-step2.md` の調査と並行してまとめた。

- 出典 URL は調査時に取得したもの。「二次」は報道・第三者 DB 経由、「未確認」は一次情報で裏が取れなかったもの
- 記事に転用する前に、使う項目だけ一次情報を開き直すこと

## 1. 攻撃の類型と実例

| 類型 | 実例 | 出典 |
| --- | --- | --- |
| A. 公開済みバージョンの中身の差し替え・タグの付け替え | 2025-02 報告: Go の `github.com/boltdb-go/bolt`。悪意ある v1.3.1 が 2021-11 に proxy.golang.org にキャッシュされ、その後 GitHub のタグを無害なコミットに付け替えた。GitHub を見ても痕跡が無く、proxy からは悪意版が配られ続けた | https://socket.dev/blog/malicious-package-exploits-go-module-proxy-caching-for-persistence |
| | 2025-03: GitHub Actions `tj-actions/changed-files`（CVE-2025-30066）。複数のタグが悪意あるコミットに付け替えられた（Go / npm 以外の例） | https://github.com/advisories/GHSA-mrrh-fwg8-r2c3 |
| B. メンテナの乗っ取り・権限移譲・トークン漏えい | 2018-11: npm `event-stream`。保守を申し出た人物に権限を渡し、`flatmap-stream` 経由で暗号資産ウォレットを狙うコードが入った | https://github.com/dominictarr/event-stream/issues/116 |
| | 2021-10: npm `ua-parser-js` 0.7.29 / 0.8.0 / 1.0.0。アカウント乗っ取りで、preinstall でマイナーとパスワード窃取を投入 | https://github.com/advisories/GHSA-pjwm-rvh2-c87w |
| | 2025-09: npm `chalk` / `debug` など 18 パッケージ以上。メンテナが偽の 2FA リセットメール（`support@npmjs.help`）でフィッシングされ、ブラウザ内で暗号資産アドレスを差し替えるコードが入った | https://www.stepsecurity.io/blog/20-popular-npm-packages-compromised-chalk-debug-strip-ansi-color-convert-wrap-ansi |
| | 2025-09: Shai-Hulud ワーム。500 超のパッケージに感染した。盗んだトークンで被害者の他のパッケージを再公開し、自己増殖した | https://www.cisa.gov/news-events/alerts/2025/09/23/widespread-supply-chain-compromise-impacting-npm-ecosystem |
| | 2024-03: xz-utils（CVE-2024-3094）。約 2 年かけて共同メンテナの信頼を得た。リリース tarball にだけ git に無いビルドスクリプトを入れた | https://research.swtch.com/xz-timeline |
| C. typosquatting（似た名前） | Go: 上記 `boltdb-go/bolt` が正規の `boltdb/bolt` を装った | A と同じ |
| D. dependency confusion（社内名と同名の公開パッケージ） | 2021-02: Alex Birsan が Apple や Microsoft などのビルドに取り込ませた（原典の Medium は未取得、二次） | https://www.bleepingcomputer.com/news/security/researcher-hacks-microsoft-apple-more-in-novel-supply-chain-attack/ |
| E. インストール時スクリプトの悪用 | `ua-parser-js` の preinstall、Shai-Hulud の postinstall | B と同じ |
| F. 経路（プロキシ・checksum database・MITM）での改ざん | 2026-08: Go の CVE-2026-56865（悪意ある GOPROXY が sumdb タイルを偽造）、CVE-2026-56864（悪意ある GOSUMDB が台帳外の内容を返す）。Go1.27 milestone で修正 | https://go.dev/issue/80744 , https://go.dev/issue/80745 （issue 本文を直接確認済み） |
| G. lock ファイル・go.sum 自体の改ざん | npm: PR で lockfile の `resolved` を攻撃者のホストに差し替える手口。`lockfile-lint` がこれを検査する（実被害の例は未確認） | https://github.com/lirantal/lockfile-lint , https://snyk.io/blog/why-npm-lockfiles-can-be-a-security-blindspot-for-injecting-malicious-modules/ |

## 2. 類型ごとに Go と npm の仕組みが効くか

| 類型 | Go | npm |
| --- | --- | --- |
| A. 同じバージョンの中身の差し替え | **効く**。go.sum と checksum database で不一致になる。proxy もキャッシュを書き換えない | 効く（lock の `integrity` と不一致になる）。ただし registry は同じバージョンの再公開を許さない |
| A'. 最初から悪意あるバージョン（boltdb-go 型） | **効かない**。最初のハッシュが悪意版のものになる。proxy の不変キャッシュが、逆に悪意版を残し続けた | 効かない |
| B. 乗っ取り後の新しいバージョン | 自動では入らない（MVS は利用者が `go get` するまで古いバージョンを使う）。更新したら入る | `^` の範囲内なら、lock を作り直したときに入る |
| C / D. 似た名前・同名 | 効かない。ただし import path が `github.com/<owner>/...` の完全なパスなので、短い名前の衝突は起きにくい。社内モジュールは `GOPRIVATE` で直接取得すれば公開側と混ざらない（推論） | 効かない。scope（`@org/`）と registry の設定で防ぐ |
| E. インストール時スクリプト | **そもそも無い**。`go get` / `go build` は依存のコードを実行しない。ただし cgo の穴（CVE-2023-29404）はあった | 既定で実行される（npm v11 時点）。`--ignore-scripts` などで止める |
| F. 経路での改ざん | 効く（署名付きの checksum database）。ただしクライアントの実装の穴は 2026-08 にも見つかった | 効く（HTTPS と lock の `integrity`）。lock に無い初回取得は registry を信じる |
| G. lock / go.sum の改ざん | **効かない**。go.sum に行があれば checksum database に問い合わせないので、偽の go.sum がマージされると正解として扱われる | 効かない。`lockfile-lint` で取得元を検査する |

## 3. Go の仕組みの詳細

- **go.sum + checksum database**: ダウンロード物のハッシュを go.sum と比べ、無ければ checksum database で検証する。https://go.dev/ref/mod#authenticating
  - 「同じバージョンなら中身は 1 つ」を保証するだけで、中身が安全かは保証しない。Go 公式も、依存はすべて信頼関係だと書いている。https://go.dev/blog/supply-chain （2022-03）
- **proxy の不変キャッシュ**: 同じバージョンには常に同じ内容を返す。改ざんには強いが、悪意版も残る（boltdb-go）。Google は proxy から削除し、脆弱性 DB に追加したと説明（二次: https://www.infoworld.com/article/3817921/malicious-package-found-in-the-go-ecosystem.html 。GO-ID は未確認）
- **ビルド時に依存のコードを実行しない**: インストールスクリプトが無い。https://go.dev/blog/supply-chain
  - ただし、ビルドしたプログラムやテストを実行すれば、依存の `init` は動く（同記事: ビルドの中にセキュリティ境界は無い）
  - `go generate` は build / test では自動で走らない。明示的に実行したときだけ任意のコマンドが動く
  - cgo: CVE-2023-29404（GO-2023-1841）。悪意ある `#cgo LDFLAGS` でビルド時に任意コード実行ができた。go1.19.10 / go1.20.5 で修正。https://pkg.go.dev/vuln/GO-2023-1841
- **toolchain**: `golang.org/toolchain` モジュールとして checksum database で検証する。`GOSUMDB=off` だと取得できず、`GOPRIVATE` / `GONOSUMDB` は効かない。https://go.dev/doc/toolchain#download
  - 手元で確認: 取得後の toolchain の `src/fmt` を書き換えても、ビルドは検出しなかった（記事の疑問 4）。検証は取得時の 1 回だけ
- **GOPRIVATE / GONOSUMDB / GOSUMDB=off**: private と判定したモジュールは proxy にも checksum database にも問い合わせない（検証しないで受け入れる）。https://go.dev/ref/mod#private-module-privacy
- **`-mod=readonly`**: Go 1.16 以降の既定。`go build` は go.mod / go.sum を勝手に書き換えない。https://go.dev/ref/mod#build-commands
- **govulncheck / vuln.go.dev**: 既知の脆弱性のうち、自分のコードから到達できるものだけを報告する。未知の悪意版は対象外。https://go.dev/doc/security/vuln/
- **retract**: 作者が自分のモジュールの版を取り下げる。`@latest` から外れるが、既存のビルドは壊さない。第三者の悪意版を止める手段ではない。https://go.dev/ref/mod#go-mod-file-retract
- **capslock**: パッケージが推移的に使う特権操作（ネットワーク・ファイル・exec など）を列挙する。依存を足したときの変化に気付ける。https://github.com/google/capslock

## 4. npm 側の対策

- **package-lock.json の `integrity` / `resolved`**: sha512 の SRI と取得元を記録する。https://docs.npmjs.com/cli/v11/configuring-npm/package-lock-json
- **npm audit**: 既知の脆弱性のみ。`npm audit signatures` で registry の署名と provenance を検証する。https://docs.npmjs.com/cli/v11/commands/npm-audit
- **provenance（Sigstore）**: どのソースからどこでビルドしたかを検証できる。公式は「悪意が無いことの保証ではない」と明記。https://docs.npmjs.com/generating-provenance-statements
- **trusted publishing（OIDC）**: 2025-07-31 に GA（npm CLI 11.5.1 以降）。CI から長期トークン無しで公開できる。https://github.blog/changelog/2025-07-31-npm-trusted-publishing-with-oidc-is-generally-available/
- **トークンと 2FA の強化**: 2025-09-22 に計画を発表し、2025-12-09 に classic token を全失効、`npm login` を 2 時間のセッション式に変えた。https://github.blog/changelog/2025-12-09-npm-classic-tokens-revoked-session-based-auth-and-cli-token-management-now-available/
- **スクリプトの抑止**: `--ignore-scripts`。git 依存経由の抜け道に `--allow-git=none`（v11.10.0 以降）。npm v12 では `allowScripts` が既定で off になる予定、と 2026-07 の changelog にある（本文は未精読）。https://github.blog/changelog/2026-07-07-npm-install-time-security-and-gat-bypass2fa-deprecation
- **クールダウン**: npm の `min-release-age`（日数）、pnpm の `minimumReleaseAge`（分）。pnpm は v10 から依存の postinstall を既定で実行しない。導入バージョンは未確認。https://docs.npmjs.com/cli/v11/using-npm/config , https://pnpm.io/supply-chain-security

## 5. 言語をまたぐ対策

- **クールダウン運用**: 公開直後のバージョンをすぐ取り込まない。chalk/debug や Shai-Hulud は数時間〜数日で検知・削除されたので、数日待つだけで避けられた可能性が高い（推論）
  - Dependabot: `cooldown.default-days`（gomod / npm 対応。セキュリティ更新には適用されない）。https://docs.github.com/en/code-security/dependabot/working-with-dependabot/dependabot-options-reference
  - Renovate: `minimumReleaseAge`。https://docs.renovatebot.com/configuration-options/#minimumreleaseage
- **SLSA**: ビルドの来歴（provenance）のレベル。L3 はビルドを分離し、署名鍵を利用者のビルド手順から隔離する。https://slsa.dev/spec/v1.0/levels
- **Sigstore**: OIDC に基づく短命証明書で署名し、透明性ログ（Rekor）に記録する。https://docs.sigstore.dev/about/overview/
- **OpenSSF Scorecard / deps.dev**: 依存先のリポジトリの運用（ブランチ保護、危険なワークフローなど）を採点する。依存グラフと OSV 勧告を引ける。https://github.com/ossf/scorecard , https://docs.deps.dev/
- **振る舞い検知（Socket など）**: boltdb-go は Socket が見つけた。既知の脆弱性 DB では見つからない悪意版に効く

## 6. チェックリスト

### Go

1. `go.mod` と `go.sum` を必ずコミットし、**go.sum の差分もレビューする**。特に、バージョンが変わっていないのにハッシュだけ変わった行は要注意（G 対策）
2. CI は `go build ./...` / `go test ./...` をそのまま使う（既定の `-mod=readonly`）。CI で `go mod tidy` して差分が出たら失敗させる
3. 社内モジュールは `GOPRIVATE` に登録する。`GOSUMDB=off` は使わない
4. Go 本体を最新のパッチに追従する（CVE-2026-56865 / 56864、cgo の CVE）。影響が気になる場合は `rm -r go.sum go.work.sum vendor/ && go mod tidy` で作り直して差分を見る
5. CI で `govulncheck ./...`。依存の更新は Dependabot の `cooldown` で数日待つ。新しい依存は `capslock` で使う権限を確認する

### npm

1. CI は `npm ci`。lockfile をコミットし、`lockfile-lint` で取得元のホストと HTTPS を検査する
2. `ignore-scripts=true` を既定にし、必要なパッケージだけ許可する（pnpm v10 以降は既定で止まる）
3. クールダウン（`min-release-age` / pnpm の `minimumReleaseAge` / Renovate の `minimumReleaseAge`）を設定する
4. CI で `npm audit signatures`
5. 公開する側は trusted publishing へ移行し、「Require two-factor authentication and disallow tokens」を選ぶ

## 7. 記事にするなら使えそうな論点

- Go の go.sum / checksum database が守るのは「同じバージョンの中身が変わらないこと」だけ。boltdb-go は、その不変性が悪意版を守る側に回った例
- Go にはインストールスクリプトが無いので、npm の被害の多く（E）は構造的に起きない。ただし実行時の `init` は防げない
- 最後の防御線は go.sum のレビュー。前回の記事の「go.sum の差分を読み飛ばさない」と、今回の「最初の 1 回だけ Google を信じる」がここでつながる

## 8. 未確認・要追跡

- Birsan の原典（Medium は 403）、Socket の公式ドキュメント、boltdb-go の Go 脆弱性 DB の ID
- `npm ci` の integrity 不一致時の挙動、npm `min-release-age` と pnpm `minimumReleaseAge` の導入バージョン
- npm の typosquatting の実例、lockfile injection の実被害の例
- CVE-2026-56865 / 56864 が入った各系列のパッチバージョン（第三者 DB は Go 1.25.13 以上・x/mod v0.40.0 以上と記述。二次）
