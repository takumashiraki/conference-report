// https://go.dev/ref/mod#go-mod-file-module
// module ディレクティブは、そのモジュール自身のモジュールパスを定義します。
module github.com/takumashiraki/conference-report/GoBash/ex-Go

// https://go.dev/ref/mod#go-mod-file-go
// go ディレクティブは、そのモジュールがどの Go バージョンの仕様・挙動を前提として書かれているかを示します。
go 1.26.8

// https://go.dev/ref/mod#go-mod-file-require
// require ディレクティブは、依存するモジュールについて必要となる最低バージョンを指定します。
require github.com/go-chi/chi/v5 v5.3.2
