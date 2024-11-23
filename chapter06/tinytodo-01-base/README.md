# tinytodo-01-base : 固定のToDoを表示するTiny ToDo

**書籍での説明箇所** :「6.3 ToDoアプリケーションで学ぶ基礎」p158

## 概要

本書で少しずつ作っていくToDoアプリケーション・「Tiny ToDo」の最初の一歩です。

まずはソースコード内にハードコードされた文字列をToDoとして表示します。
HTMLは、 [`templates/todo.html`](./templates/todo.html) にあるものを、Goのテンプレート機能で展開して返します。

本サンプルは、以下のURLで実働公開しています。

https://tinytodo-01-base.webtech.littleforest.jp/todo

## 主要コード

- **サーバサイド**
  - [`main.go`](./main.go) サーバサイドのメイン処理 (リスト6.4, p159)
- **HTMLテンプレート**
  - [`templates/todo.html`](./templates/todo.html) ToDo画面のHTMLテンプレート (リスト6.5, p162)

