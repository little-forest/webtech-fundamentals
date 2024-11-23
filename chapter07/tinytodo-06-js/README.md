# tinytodo-06-js : JavaScriptによる編集処理を組み込んだTiny ToDo

## 概要

SPA化の手始めとして、Tiny ToDoの画面でJavaScriptによる編集を可能にしたものです。
ToDoをクリックするとToDoを編集可能にし、Save、Cancelボタンを表示します。
本サンプルの時点では、編集内容をサーバへ送信する処理を実装していません。

本サンプルは、以下のURLで実働公開しています。

https://tinytodo-06-js.webtech.littleforest.jp/todo

## 主要コード

- **クライアントサイド**
  - [`static/todo.js`](./static/todo.js) ブラウザ上で動作するJavaScriptコード (リスト7.7, p256)
- **HTMLテンプレート**
  - [`templates/todo.html`](./templates/todo.html) ToDo画面のHTMLテンプレート (リスト7.5, p256)
- **サーバサイド**
  - [`main.go`](./main.go) サーバサイドのメイン処理
  - [`session.go`](./session.go) セッション管理処理

