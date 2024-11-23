# tinytodo-04-session : セッション管理を導入したTiny ToDo

**書籍での説明箇所** :「6.5.4 セッションの実装 - セッションID発行の実装例」p198

## 概要

セッション管理機能を導入したTiny ToDoです。

初回アクセス時にセッションIDを発行してCookieへ保持させます。
セッションIDに紐付いてToDoリストを管理するため、セッションごとに異なるToDoリストを管理できます。

本サンプルは、以下のURLで実働公開しています。

https://tinytodo-04-session.webtech.littleforest.jp/todo

## 主要コード

- **サーバサイド**
  - [`main.go`](./main.go) サーバサイドのメイン処理
  - [`session.go`](./session.go) セッション管理処理 (リスト6.10, p198)
  - [`todo.go`](./todo.go) ToDo管理処理 (リスト6.11, p200)
- **HTMLテンプレート**
  - [`templates/todo.html`](./templates/todo.html) ToDo画面のHTMLテンプレート

