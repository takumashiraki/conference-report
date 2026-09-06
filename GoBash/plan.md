# 登壇の案

## タイトル

Go の人は go.sum をなんとなくロックファイルだと思っていることが多く、 package-lock.json を引き合いに出して、この問いについて深ぼる。

- 「Go にロックファイルが要らないのはなぜか」
- 「npm との違いを手がかりに Go の依存管理を理解したい」
- 「Go Modulesの裏側 ― go.mod / go.sumと依存解決の仕組み」

## 構成案

### 構成イメージ1

1. 「依存関係どうなってるの？？？」
   - go.mod
   - go.sum
   - $GOMODCACHE
2. npmとの比較
   - package.json
   - package-lock.json
   - node_modules
3. go.modを理解する
4. go.sumを理解する
5. $GOMODCACHEを理解する
6. 実際にGoが依存関係を解決する様子を見る
7. npmとGo Modulesの違い
8. 「なぜこういう設計なのか？」
9. Go Modulesの設計思想
10. Goの依存管理がちょっと分かった！
    - MVS の設計思想なので、それを記載

### 構成イメージ2

1. npm は lock がないと再現しない → なぜ？（SemVer レンジ + 解決アルゴリズムが非決定的）
2. Go は go.mod だけで再現する（MVS = 最小バージョン選択）
3. じゃあ go.sum は何者？ → ロックではなく検証（ハッシュ、sum.golang.org、Merkle tree）
4. dep 時代には Gopkg.lock があった → なぜ捨てられたか
5. オチ:「go.mod は宣言、go.sum は証明」

dep の話が「歴史コーナー」ではなく、論の必然として入るのがこの構成の強みです。