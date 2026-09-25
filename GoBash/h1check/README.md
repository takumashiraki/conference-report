# README

[question.md](../question.md) の記載を試すためのディレクトリ

> Go から同じ計算をするなら、`golang.org/x/mod/sumdb/dirhash` の 1 行で済みます。
> 
> ```go
> dirhash.HashZip(".../v5.3.2.zip", dirhash.Hash1)
> // h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=  → go.sum と一致
> ```

```bash
go run . /Users/opm008296/go/pkg/mod/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.zip /Users/opm008296/go/pkg/mod/cache/download/github.com/go-chi/chi/v5/@v/v5.3.2.mod
zip    h1:5YQkICvTCSZ25hoRsyJazN0scjzKGiu4VAUc7H1o1nY=
go.mod h1:R+tYY2hNuVUUjxoPtqUdgBqevM9s9njzkTLutVsOCto=
```
