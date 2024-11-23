# tinytodo-03-prg : Post-Redirect-Getを導入したTiny ToDo

**書籍での説明箇所** :「6.4.3 Post-Redirect-Getの実装」p170

## 概要

[`tinytodo-02-add`](../tinytodo-02-add) におけるToDo追加時の画面遷移を、Post-Redirect-Getパターンに書き換えたものです。

## 主要コード

- **サーバサイド**
  - [`main.go`](./main.go) サーバサイドのメイン処理 (リスト6.9, p171)
- **HTMLテンプレート**
  - [`templates/todo.html`](./templates/todo.html) ToDo画面のHTMLテンプレート

